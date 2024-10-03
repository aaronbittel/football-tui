PROJECT_NAME = algo
OUTPUT_DIR = ./bin
SOURCE_DIR = ./cmd/algo

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(OUTPUT_DIR)/$(PROJECT_NAME) $(SOURCE_DIR)

.PHONY: build-win
build-win:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(PROJECT_NAME).exe $(SOURCE_DIR)

.PHONY: run
run:
	go run $(SOURCE_DIR)/main.go

.PHONY: run-win
run-win: build-win
	$(OUTPUT_DIR)/$(PROJECT_NAME).exe

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)

.PHONY: help
help:
	@echo "Makefile for building and running $(PROJECT_NAME)"
	@echo "Available commands:"
	@echo "  make all          Build the executable for the current OS"
	@echo "  make build        Build the executable"
	@echo "  make build-win    Build the executable for Windows"
	@echo "  make run          Run the application"
	@echo "  make run-win      Run the Windows executable"
	@echo "  make clean        Remove build artifacts"
	@echo "  make help         Display this help message"
