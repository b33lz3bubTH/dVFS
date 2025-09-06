.PHONY: build run clean test docker-build docker-run docker-compose-up docker-compose-down

# Build the storage node binary
build:
	go build -o bin/storage-node ./cmd/node

# Run the storage node
run:
	go run ./cmd/node

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/storage-node-linux ./cmd/node

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/storage-node-darwin ./cmd/node

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/storage-node-windows.exe ./cmd/node

# Build all platforms
build-all: build-linux build-darwin build-windows

# Docker commands
docker-build:
	docker build -t storage-node .

docker-run:
	docker run -p 8080:8080 -e INSTANCE_ID=node-docker-1 storage-node

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f

# Development with custom instance ID
run-instance:
	INSTANCE_ID=node-dev-$(shell date +%s) go run ./cmd/node

# Test distributed setup
test-distributed:
	@echo "Testing distributed storage nodes..."
	@echo "Node 1: http://localhost:8081/api/v1/instance"
	@echo "Node 2: http://localhost:8082/api/v1/instance"
	@echo "Node 3: http://localhost:8083/api/v1/instance"
	@echo ""
	@echo "Upload to node 1:"
	@echo 'echo "Hello from Node 1" | curl -X POST http://localhost:8081/api/v1/files -H "Content-Type: text/plain" -H "X-Filename: hello1.txt" --data-binary @-'
	@echo ""
	@echo "Upload to node 2:"
	@echo 'echo "Hello from Node 2" | curl -X POST http://localhost:8082/api/v1/files -H "Content-Type: text/plain" -H "X-Filename: hello2.txt" --data-binary @-'
