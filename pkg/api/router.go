package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dvfs/storage-node/pkg/api/resources/files"
	"github.com/dvfs/storage-node/pkg/storage"
)

// Router handles all API routing
type Router struct {
	filesHandler *files.Handler
	instanceID   string
	startTime    time.Time
}

// NewRouter creates a new API router
func NewRouter(storage *storage.FileStorage, instanceID string) *Router {
	return &Router{
		filesHandler: files.NewHandler(storage),
		instanceID:   instanceID,
		startTime:    time.Now(),
	}
}

// Routes returns the configured HTTP handler
func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	// API v1 routes
	mux.HandleFunc("/api/v1/files", r.handleFiles)
	mux.HandleFunc("/api/v1/files/", r.handleFilesWithID)
	
	// Instance-specific routes
	mux.HandleFunc("/api/v1/instance", r.getInstanceInfo)
	
	// Health check
	mux.HandleFunc("/health", r.healthCheck)
	
	// Root endpoint
	mux.HandleFunc("/", r.rootHandler)

	return r.loggingMiddleware(mux)
}

// handleFiles routes requests to /api/v1/files
func (r *Router) handleFiles(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		r.filesHandler.UploadFile(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleFilesWithID routes requests to /api/v1/files/{id} and /api/v1/files/{id}/info
func (r *Router) handleFilesWithID(w http.ResponseWriter, req *http.Request) {
	// Check if it's an info request
	if r.isInfoRequest(req.URL.Path) {
		switch req.Method {
		case http.MethodGet:
			r.filesHandler.GetFileInfo(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Regular file operations
	switch req.Method {
	case http.MethodGet:
		r.filesHandler.GetFile(w, req)
	case http.MethodDelete:
		r.filesHandler.DeleteFile(w, req)
	case http.MethodHead:
		r.filesHandler.CheckFileExists(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getInstanceInfo handles GET /api/v1/instance
func (r *Router) getInstanceInfo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(r.startTime)
	
	instanceInfo := map[string]interface{}{
		"instance_id": r.instanceID,
		"service":     "storage-node",
		"version":     "1.0.0",
		"uptime":      uptime.String(),
		"started_at":  r.startTime.Format(time.RFC3339),
		"endpoints": map[string]string{
			"files":   "/api/v1/files",
			"health":  "/health",
			"instance": "/api/v1/instance",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(instanceInfo)
}

// isInfoRequest checks if the request is for file info
func (r *Router) isInfoRequest(path string) bool {
	return len(path) > 5 && path[len(path)-5:] == "/info"
}

// healthCheck handles GET /health
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(r.startTime)
	
	healthInfo := map[string]interface{}{
		"status":      "healthy",
		"service":     "storage-node",
		"instance_id": r.instanceID,
		"uptime":      uptime.String(),
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthInfo)
}

// rootHandler handles requests to the root path
func (r *Router) rootHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	rootInfo := map[string]interface{}{
		"service":     "Distributed File System Storage Node",
		"version":     "1.0.0",
		"instance_id": r.instanceID,
		"endpoints": map[string]string{
			"upload":    "POST /api/v1/files",
			"download":  "GET /api/v1/files/{id}",
			"info":      "GET /api/v1/files/{id}/info",
			"delete":    "DELETE /api/v1/files/{id}",
			"exists":    "HEAD /api/v1/files/{id}",
			"health":    "GET /health",
			"instance":  "GET /api/v1/instance",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rootInfo)
}

// loggingMiddleware logs HTTP requests
func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s %s %s", r.instanceID, req.Method, req.URL.Path, req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}
