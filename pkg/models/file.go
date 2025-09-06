package models

import (
	"time"
)

// FileMetadata represents metadata for a stored file
type FileMetadata struct {
	ID          string    `json:"id"`
	OriginalName string   `json:"original_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Extension   string    `json:"extension"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileUploadRequest represents the request structure for file upload
type FileUploadRequest struct {
	Content     []byte `json:"-"`
	ContentType string `json:"content_type,omitempty"`
	FileName    string `json:"file_name,omitempty"`
}

// FileUploadResponse represents the response structure for file upload
type FileUploadResponse struct {
	ID          string `json:"id"`
	OriginalName string `json:"original_name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Extension   string `json:"extension"`
	URL         string `json:"url"`
}

// FileInfoResponse represents file information response
type FileInfoResponse struct {
	ID          string    `json:"id"`
	OriginalName string   `json:"original_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Extension   string    `json:"extension"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	URL         string    `json:"url"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
