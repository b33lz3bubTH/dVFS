# Distributed Virtual File System Gateway

A gateway service for managing distributed file storage across multiple nodes.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Set up the database:
```bash
npx prisma migrate dev
```

3. Configure storage nodes in `config/nodes.json`:
```json
{
  "nodes": [
    {
      "id": "node1",
      "baseUrl": "http://localhost:8081",
      "healthEndpoint": "/api/v1/instance"
    }
  ]
}
```

## Running the Service

Development mode:
```bash
make dev
```

Production mode:
```bash
make build
make start
```

## API Usage

All requests require a `user-email` header for authentication.

### Upload File
```bash
curl -X POST http://localhost:3000/api/v1/files \
  -H "user-email: user@example.com" \
  -F "file=@/path/to/file.txt" \
  -F "virtualPath=/folder1/file.txt"

curl -X POST http://localhost:3000/api/v1/files   -H "user-email: user@example.com"   -F "file=@./secret.txt"   -F "virtualPath=/sourav-backup/secret.txt"
```

### Download File
```bash
curl http://localhost:3000/api/v1/files/{fileId} \
  -H "user-email: user@example.com"
```

### Get File Info
```bash
curl http://localhost:3000/api/v1/files/{fileId}/info \
  -H "user-email: user@example.com"
```

### Check File Exists
```bash
curl -I http://localhost:3000/api/v1/files/{fileId} \
  -H "user-email: user@example.com"
```

### Delete File
```bash
curl -X DELETE http://localhost:3000/api/v1/files/{fileId} \
  -H "user-email: user@example.com"
```

### Create Folder
```bash
curl -X POST http://localhost:3000/api/v1/folders \
  -H "user-email: user@example.com" \
  -H "Content-Type: application/json" \
  -d '{"path": "/folder1/subfolder"}'
```

### Get Folder Tree
```bash
curl http://localhost:3000/api/v1/tree \
  -H "user-email: user@example.com"
```

## Features

- User-based file system isolation
- Virtual folder structure
- Distributed storage across multiple nodes
- Health checking of storage nodes
- Round-robin load balancing for uploads
- SQLite database for metadata storage
- File and folder operations