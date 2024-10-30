package handlers

import (
	"csv-reader/models"
	"csv-reader/services"
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
)

type Handler struct {
	csvService *services.CSVService
	staticFS   fs.FS
}

func NewHandler(csvService *services.CSVService, staticFS fs.FS) (*Handler, error) {
	return &Handler{
		csvService: csvService,
		staticFS:   staticFS,
	}, nil
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Read the embedded index.html file
	content, err := fs.ReadFile(h.staticFS, "index.html")
	if err != nil {
		http.Error(w, "Error reading index.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func (h *Handler) HandleHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	table := h.csvService.GetCurrentTable()
	if table == nil {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(table.Headers); err != nil {
		http.Error(w, "Error encoding headers", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit file size to 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension
	if ext := filepath.Ext(header.Filename); ext != ".csv" {
		http.Error(w, "Only CSV files are allowed", http.StatusBadRequest)
		return
	}

	if err := h.csvService.ProcessCSVFile(file); err != nil {
		http.Error(w, "Error processing CSV file: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "File uploaded successfully",
	})
}

func (h *Handler) HandleData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	table := h.csvService.GetCurrentTable()
	if table == nil {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(table); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.Columns) == 0 {
		http.Error(w, "No columns specified", http.StatusBadRequest)
		return
	}

	result := h.csvService.QueryRecords(req.Columns, req.Filter)
	if result == nil {
		result = []models.Record{} // Return empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding query results", http.StatusInternalServerError)
		return
	}
}

// HandleStatic serves static files from the embedded filesystem
func (h *Handler) HandleStatic(w http.ResponseWriter, r *http.Request) {
	// Remove /static/ prefix from path
	path := r.URL.Path[8:]

	content, err := fs.ReadFile(h.staticFS, path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set content type based on file extension
	contentType := "text/plain"
	switch filepath.Ext(path) {
	case ".html":
		contentType = "text/html"
	case ".js":
		contentType = "application/javascript"
	case ".css":
		contentType = "text/css"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}
