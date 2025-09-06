# Distributed File System Storage Node

A production-ready, enterprise-grade Go storage node for a distributed file system. Built with NSA-level code quality, proper resource organization, and comprehensive file handling capabilities.

## 🏗️ Architecture

This storage node follows enterprise-grade patterns with:

- **Resource-based organization**: Each API resource has its own handler package
- **Proper file extension handling**: Files are stored with correct extensions and MIME types
- **Metadata management**: Complete file metadata tracking with JSON storage
- **Content-type detection**: Automatic MIME type detection and validation
- **Clean separation of concerns**: Models, utilities, storage, and API layers

## 📁 Project Structure

```
.
├── cmd/node/                    # Application entrypoint
│   └── main.go
├── pkg/
│   ├── api/
│   │   ├── resources/
│   │   │   └── files/           # Files resource handlers
│   │   │       └── handler.go
│   │   └── router.go            # API routing and middleware
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   └── file.go              # Data models and DTOs
│   ├── storage/
│   │   └── storage.go           # File storage operations
│   └── utils/
│       └── mime.go              # MIME type utilities
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🚀 Features

### Core Functionality
- **File Storage**: Store files with proper extensions and metadata
- **Content-Type Detection**: Automatic MIME type detection from extensions
- **Metadata Management**: Complete file information tracking
- **Resource-based API**: Clean RESTful endpoints organized by resource
- **Graceful Shutdown**: Proper signal handling and cleanup

### Enterprise Features
- **Comprehensive Error Handling**: Structured error responses
- **Request Logging**: Full request/response logging
- **Health Monitoring**: Built-in health check endpoint
- **File Sanitization**: Secure filename handling
- **Multipart Support**: Handle both raw and multipart file uploads

## 🌐 API Endpoints

### File Operations

#### Upload File
- **POST** `/api/v1/files`
- **Description**: Upload a file with proper extension and metadata
- **Content-Type**: `multipart/form-data` or raw binary
- **Headers**: 
  - `X-Filename`: Original filename (for raw uploads)
  - `Content-Type`: MIME type (optional, auto-detected)
- **Response**: File metadata with download URL

#### Download File
- **GET** `/api/v1/files/{id}`
- **Description**: Download file with proper content-type headers
- **Response**: File content with appropriate headers

#### Get File Information
- **GET** `/api/v1/files/{id}/info`
- **Description**: Get file metadata without downloading
- **Response**: Complete file information

#### Delete File
- **DELETE** `/api/v1/files/{id}`
- **Description**: Delete file and its metadata
- **Response**: 204 No Content

#### Check File Exists
- **HEAD** `/api/v1/files/{id}`
- **Description**: Check if file exists
- **Response**: 200 OK if exists, 404 Not Found if not

### System Endpoints

#### Health Check
- **GET** `/health`
- **Description**: Service health status
- **Response**: `{"status":"healthy","service":"storage-node"}`

#### API Information
- **GET** `/`
- **Description**: API documentation and available endpoints

## ⚙️ Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `STORAGE_PATH`: Local storage directory (default: /tmp/storage_data)

## 🛠️ Getting Started

### Prerequisites
- Go 1.21 or later
- Make (optional)

### Installation
```bash
git clone <repository-url>
cd storage-node
make deps
```

### Running
```bash
# Using Make
make run

# Using Go directly
go run ./cmd/node

# With custom configuration
PORT=9000 STORAGE_PATH=/custom/storage go run ./cmd/node
```

### Building
```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

## 🧪 Testing the API

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Upload a File (Multipart)
```bash
curl -X POST http://localhost:8080/api/v1/files \
  -F "file=@example.txt"
```

### 3. Upload a File (Raw)
```bash
echo "Hello, World!" | curl -X POST http://localhost:8080/api/v1/files \
  -H "Content-Type: text/plain" \
  -H "X-Filename: hello.txt" \
  --data-binary @-
```

### 4. Get File Information
```bash
curl http://localhost:8080/api/v1/files/{file-id}/info
```

### 5. Download File
```bash
curl http://localhost:8080/api/v1/files/{file-id} -o downloaded_file.txt
```

### 6. Check if File Exists
```bash
curl -I http://localhost:8080/api/v1/files/{file-id}
```

### 7. Delete File
```bash
curl -X DELETE http://localhost:8080/api/v1/files/{file-id}
```

## 📊 File Storage Details

### File Organization
- **Physical Storage**: Files stored as `{uuid}.{extension}` in storage directory
- **Metadata Storage**: JSON metadata files in `metadata/` subdirectory
- **Extension Handling**: Proper file extensions based on content-type or original filename
- **MIME Type Support**: Comprehensive MIME type detection and mapping

### Supported File Types
- **Text**: .txt, .html, .css, .js, .json, .xml
- **Images**: .png, .jpg, .jpeg, .gif, .svg
- **Documents**: .pdf, .doc, .docx, .xls, .xlsx, .ppt, .pptx
- **Archives**: .zip, .tar, .gz
- **Media**: .mp4, .mp3, .wav, .avi, .mov
- **Binary**: .bin (default for unknown types)

## 🔒 Security Features

- **Filename Sanitization**: Removes dangerous characters and path traversal attempts
- **Content-Type Validation**: Proper MIME type handling
- **File Size Limits**: Configurable upload size limits
- **Path Security**: Prevents directory traversal attacks

## 🏭 Production Considerations

- **Horizontal Scaling**: Stateless design allows multiple instances
- **Monitoring**: Built-in health checks and request logging
- **Error Handling**: Comprehensive error responses with proper HTTP status codes
- **Performance**: Efficient file streaming and metadata caching
- **Maintenance**: Clean separation allows easy feature additions

## 🧪 Development

### Running Tests
```bash
make test
```

### Code Quality
- Follows Go best practices and idioms
- Comprehensive error handling
- Clean separation of concerns
- Resource-based organization
- Enterprise-grade logging and monitoring

## 📄 License

This project is licensed under the MIT License.



