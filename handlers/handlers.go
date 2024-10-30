package handlers

import (
	"csv-reader/models"
	"csv-reader/services"
	"encoding/json"
	"net/http"
)

type Handler struct {
	csvService *services.CSVService
}

func NewHandler(csvService *services.CSVService) *Handler {
	return &Handler{
		csvService: csvService,
	}
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func (h *Handler) HandleHeaders(w http.ResponseWriter, r *http.Request) {
	table := h.csvService.GetCurrentTable()
	json.NewEncoder(w).Encode(table.Headers)
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.csvService.ProcessCSVFile(file); err != nil {
		http.Error(w, "Error processing CSV file", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (h *Handler) HandleData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.csvService.GetCurrentTable())
}

func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.csvService.QueryRecords(req.Columns, req.Filter)
	json.NewEncoder(w).Encode(result)
}
