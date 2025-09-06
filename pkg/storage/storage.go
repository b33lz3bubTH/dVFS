package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dvfs/storage-node/pkg/models"
	"github.com/dvfs/storage-node/pkg/utils"
	"github.com/google/uuid"
)

// FileStorage handles local file operations with metadata
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new FileStorage instance
func NewFileStorage(basePath string) (*FileStorage, error) {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create metadata directory
	metadataPath := filepath.Join(basePath, "metadata")
	if err := os.MkdirAll(metadataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metadata directory: %w", err)
	}

	return &FileStorage{
		basePath: basePath,
	}, nil
}

// Store saves a file with metadata and returns file information
func (fs *FileStorage) Store(content []byte, originalName, contentType string) (*models.FileMetadata, error) {
	// Generate a new UUID for the file
	fileID := uuid.New().String()
	
	// Determine file extension
	var extension string
	if contentType != "" {
		extension = utils.GetExtensionFromContentType(contentType)
	} else if originalName != "" {
		extension = filepath.Ext(originalName)
	} else {
		extension = ".bin"
	}
	
	// Create filename with extension
	filename := fileID + extension
	filePath := filepath.Join(fs.basePath, filename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to file
	_, err = file.Write(content)
	if err != nil {
		// Clean up the file if write fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	// Create metadata
	now := time.Now()
	metadata := &models.FileMetadata{
		ID:           fileID,
		OriginalName: utils.SanitizeFileName(originalName),
		ContentType:  contentType,
		Size:         int64(len(content)),
		Extension:    extension,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save metadata
	if err := fs.saveMetadata(metadata); err != nil {
		// Clean up the file if metadata save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return metadata, nil
}

// Retrieve returns file content and metadata for the given file ID
func (fs *FileStorage) Retrieve(fileID string) ([]byte, *models.FileMetadata, error) {
	// Load metadata first
	metadata, err := fs.loadMetadata(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("metadata not found: %s", fileID)
	}

	// Construct file path
	filename := fileID + metadata.Extension
	filePath := filepath.Join(fs.basePath, filename)
	
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, metadata, nil
}

// GetMetadata returns only metadata for the given file ID
func (fs *FileStorage) GetMetadata(fileID string) (*models.FileMetadata, error) {
	return fs.loadMetadata(fileID)
}

// Delete removes a file and its metadata by ID
func (fs *FileStorage) Delete(fileID string) error {
	// Load metadata to get the extension
	metadata, err := fs.loadMetadata(fileID)
	if err != nil {
		return fmt.Errorf("metadata not found: %s", fileID)
	}

	// Construct file path
	filename := fileID + metadata.Extension
	filePath := filepath.Join(fs.basePath, filename)
	
	// Remove file
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Remove metadata
	metadataPath := fs.getMetadataPath(fileID)
	err = os.Remove(metadataPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}

// Exists checks if a file exists
func (fs *FileStorage) Exists(fileID string) bool {
	metadata, err := fs.loadMetadata(fileID)
	if err != nil {
		return false
	}

	filename := fileID + metadata.Extension
	filePath := filepath.Join(fs.basePath, filename)
	_, err = os.Stat(filePath)
	return !os.IsNotExist(err)
}

// saveMetadata saves file metadata to disk
func (fs *FileStorage) saveMetadata(metadata *models.FileMetadata) error {
	metadataPath := fs.getMetadataPath(metadata.ID)
	
	file, err := os.Create(metadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metadata)
}

// loadMetadata loads file metadata from disk
func (fs *FileStorage) loadMetadata(fileID string) (*models.FileMetadata, error) {
	metadataPath := fs.getMetadataPath(fileID)
	
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata models.FileMetadata
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// getMetadataPath returns the path for metadata file
func (fs *FileStorage) getMetadataPath(fileID string) string {
	return filepath.Join(fs.basePath, "metadata", fileID+".json")
}
