package files

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dvfs/storage-node/pkg/models"
	"github.com/dvfs/storage-node/pkg/storage"
	"github.com/dvfs/storage-node/pkg/utils"
)

// Handler handles file-related HTTP requests
type Handler struct {
	storage *storage.FileStorage
}

// NewHandler creates a new files handler
func NewHandler(storage *storage.FileStorage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// UploadFile handles POST /api/v1/files
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form if present
	var content []byte
	var originalName, contentType string

	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Handle multipart form data
		err := r.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			h.sendError(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			h.sendError(w, "No file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err = io.ReadAll(file)
		if err != nil {
			h.sendError(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		originalName = header.Filename
		contentType = header.Header.Get("Content-Type")
	} else {
		// Handle raw file content
		var err error
		content, err = io.ReadAll(r.Body)
		if err != nil {
			h.sendError(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Get filename from header or use default
		originalName = r.Header.Get("X-Filename")
		if originalName == "" {
			originalName = "uploaded_file"
		}

		// Get content type from header or detect from filename
		contentType = r.Header.Get("Content-Type")
		if contentType == "" {
			contentType = utils.GetContentTypeFromExtension(originalName)
		}
	}

	// Store the file
	metadata, err := h.storage.Store(content, originalName, contentType)
	if err != nil {
		log.Printf("Failed to store file: %v", err)
		h.sendError(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	// Create response
	response := &models.FileUploadResponse{
		ID:           metadata.ID,
		OriginalName: metadata.OriginalName,
		ContentType:  metadata.ContentType,
		Size:         metadata.Size,
		Extension:    metadata.Extension,
		URL:          fmt.Sprintf("/api/v1/files/%s", metadata.ID),
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// GetFile handles GET /api/v1/files/{id}
func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := h.extractFileID(r.URL.Path)
	if fileID == "" {
		h.sendError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Get file content and metadata
	content, metadata, err := h.storage.Retrieve(fileID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendError(w, "File not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to retrieve file %s: %v", fileID, err)
			h.sendError(w, "Failed to retrieve file", http.StatusInternalServerError)
		}
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", metadata.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(metadata.Size, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", metadata.OriginalName))
	w.Header().Set("X-File-ID", metadata.ID)
	w.Header().Set("X-Original-Name", metadata.OriginalName)

	// Write file content
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// GetFileInfo handles GET /api/v1/files/{id}/info
func (h *Handler) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := h.extractFileIDFromInfoPath(r.URL.Path)
	if fileID == "" {
		h.sendError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Get metadata only
	metadata, err := h.storage.GetMetadata(fileID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendError(w, "File not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to get metadata for file %s: %v", fileID, err)
			h.sendError(w, "Failed to get file info", http.StatusInternalServerError)
		}
		return
	}

	// Create response
	response := &models.FileInfoResponse{
		ID:           metadata.ID,
		OriginalName: metadata.OriginalName,
		ContentType:  metadata.ContentType,
		Size:         metadata.Size,
		Extension:    metadata.Extension,
		CreatedAt:    metadata.CreatedAt,
		UpdatedAt:    metadata.UpdatedAt,
		URL:          fmt.Sprintf("/api/v1/files/%s", metadata.ID),
	}

	h.sendJSON(w, response, http.StatusOK)
}

// DeleteFile handles DELETE /api/v1/files/{id}
func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := h.extractFileID(r.URL.Path)
	if fileID == "" {
		h.sendError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Delete the file
	err := h.storage.Delete(fileID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendError(w, "File not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to delete file %s: %v", fileID, err)
			h.sendError(w, "Failed to delete file", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckFileExists handles HEAD /api/v1/files/{id}
func (h *Handler) CheckFileExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodHead {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := h.extractFileID(r.URL.Path)
	if fileID == "" {
		h.sendError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Check if file exists
	if h.storage.Exists(fileID) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// extractFileID extracts the file ID from the URL path
func (h *Handler) extractFileID(path string) string {
	// Expected format: /api/v1/files/{id}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "v1" && parts[2] == "files" {
		return parts[3]
	}
	return ""
}

// extractFileIDFromInfoPath extracts the file ID from info path
func (h *Handler) extractFileIDFromInfoPath(path string) string {
	// Expected format: /api/v1/files/{id}/info
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 5 && parts[0] == "api" && parts[1] == "v1" && parts[2] == "files" && parts[4] == "info" {
		return parts[3]
	}
	return ""
}

// sendJSON sends a JSON response
func (h *Handler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	errorResp := &models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}
