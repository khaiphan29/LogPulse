# Define path to main source files
API_FILE = cmd/api/main.go
API_PORT = 8080

# Ensure that clean, build, and run are treated as commands to execute, not as files or directories
.PHONY: clean build run help

run:
	@echo "Starting Go Server..."
	go run $(API_FILE) $(API_PORT)

help:
	@echo "Available commands:"
	@echo "  make run   - Start the Go server"
