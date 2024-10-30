package main

import (
	"csv-reader/handlers"
	"csv-reader/services"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Initialize logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Start system tray
	go systray.Run(onReady, onExit)

	// Initialize services
	csvService := services.NewCSVService()

	// Create a sub-filesystem for the static directory
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to create static sub-filesystem:", err)
	}

	// Initialize handlers with the embedded filesystem
	handler, err := handlers.NewHandler(csvService, staticFS)
	if err != nil {
		log.Fatal("Failed to initialize handlers:", err)
	}

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      setupRoutes(handler, staticFS),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Open browser automatically
	go func() {
		// Wait a bit for the server to start
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:8080")
	}()

	fmt.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func setupRoutes(handler *handlers.Handler, staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// Serve static files from embedded filesystem
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// API routes
	mux.HandleFunc("/", handler.HandleHome)
	mux.HandleFunc("/api/headers", handler.HandleHeaders)
	mux.HandleFunc("/api/query", handler.HandleQuery)
	mux.HandleFunc("/api/upload", handler.HandleUpload)
	mux.HandleFunc("/api/data", handler.HandleData)

	// Add middleware
	return addMiddleware(mux)
}

func addMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add common headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Log request
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Serve request
		handler.ServeHTTP(w, r)
	})
}

func onReady() {
	systray.SetTitle("CSV Reader")
	systray.SetTooltip("CSV Reader Application")

	mOpen := systray.AddMenuItem("Open UI", "Open the web interface")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser("http://localhost:8080")
			case <-mQuit.ClickedCh:
				gracefulShutdown()
			}
		}
	}()
}

func onExit() {
	// Cleanup when app exits
	log.Println("Application shutting down...")
}

func gracefulShutdown() {
	log.Println("Initiating graceful shutdown...")
	// Add any cleanup tasks here
	systray.Quit()
	os.Exit(0)
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
