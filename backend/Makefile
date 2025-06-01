.PHONY: build run test clean docker-build docker-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=ruixin-api
MAIN_PATH=cmd/api/main.go

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Test the project
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Download dependencies
deps:
	$(GOMOD) download

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -s -w .

# Build Docker image
docker-build:
	docker build -t ruixin-platform:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --name ruixin-api ruixin-platform:latest

# Database migrations
migrate-up:
	migrate -path ./migrations -database "mysql://root:password@tcp(localhost:3306)/ruixin_platform" up

migrate-down:
	migrate -path ./migrations -database "mysql://root:password@tcp(localhost:3306)/ruixin_platform" down

# Development with hot reload
dev:
	air -c .air.toml

# Initialize database
init-db:
	cd scripts && $(GOCMD) run init_db.go