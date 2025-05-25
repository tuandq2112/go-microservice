.PHONY: proto build run-gateway run-user clean init deps

# Initialize project
init: deps proto
	@echo "Project initialized successfully!"

# Install Go dependencies
deps:
	@echo "Installing dependencies..."
	@go install github.com/bufbuild/buf/cmd/buf@v1.26.1
	@cd shared && go mod tidy
	@cd shared/proto/schema && $(MAKE) generate
	@cd api-gateway && go mod tidy
	@cd user-service && go mod tidy


# Build all services
build: proto
	@echo "Building services..."
	@mkdir -p bin
	@go build -o bin/api-gateway api-gateway/main.go
	@go build -o bin/user-service user-service/main.go

# Run API Gateway
run-gateway:
	@echo "Starting API Gateway..."
	@go run api-gateway/main.go

# Run User Service
run-user:
	@echo "Starting User Service..."
	@go run user-service/main.go

# Clean generated files and binaries
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@cd shared/proto && $(MAKE) clean

# Run all services (in separate terminals)
run-all:
	@echo "Starting all services..."
	@make run-user & make run-gateway

# Help command
help:
	@echo "🛠  Available commands:"
	@echo "   make init       - Initialize project (install dependencies)"
	@echo "   make proto      - Generate protobuf files"
	@echo "   make build      - Build all services"
	@echo "   make run-gateway - Run API Gateway"
	@echo "   make run-user   - Run User Service"
	@echo "   make run-all    - Run all services"
	@echo "   make clean      - Clean generated files"
	@echo "   make deps       - Install dependencies"
	@echo "   make help       - Show this help message"