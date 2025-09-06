package utils

import (
	"mime"
	"path/filepath"
	"strings"
)

// GetContentTypeFromExtension returns MIME type based on file extension
func GetContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Common extensions mapping
	extToMime := map[string]string{
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".pdf":  "application/pdf",
		".zip":  "application/zip",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}
	
	if mimeType, exists := extToMime[ext]; exists {
		return mimeType
	}
	
	// Fallback to Go's built-in MIME type detection
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	
	// Default to binary if unknown
	return "application/octet-stream"
}

// GetExtensionFromContentType returns file extension based on MIME type
func GetExtensionFromContentType(contentType string) string {
	// Remove charset and other parameters
	contentType = strings.Split(contentType, ";")[0]
	
	mimeToExt := map[string]string{
		"text/plain":                      ".txt",
		"text/html":                       ".html",
		"text/css":                        ".css",
		"application/javascript":          ".js",
		"application/json":                ".json",
		"application/xml":                 ".xml",
		"application/pdf":                 ".pdf",
		"application/zip":                 ".zip",
		"application/x-tar":               ".tar",
		"application/gzip":                ".gz",
		"image/png":                       ".png",
		"image/jpeg":                      ".jpg",
		"image/gif":                       ".gif",
		"image/svg+xml":                   ".svg",
		"video/mp4":                       ".mp4",
		"audio/mpeg":                      ".mp3",
		"audio/wav":                       ".wav",
		"video/x-msvideo":                 ".avi",
		"video/quicktime":                 ".mov",
		"application/msword":              ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel":        ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"application/vnd.ms-powerpoint":   ".ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	}
	
	if ext, exists := mimeToExt[contentType]; exists {
		return ext
	}
	
	// Default extension for unknown types
	return ".bin"
}

// SanitizeFileName removes dangerous characters from filename
func SanitizeFileName(filename string) string {
	// Remove path separators and other dangerous characters
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")
	filename = strings.ReplaceAll(filename, "~", "_")
	
	// Limit length
	if len(filename) > 255 {
		filename = filename[:255]
	}
	
	return filename
}
