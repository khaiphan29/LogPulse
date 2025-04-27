# Define path to main source files
API_FILE = cmd/api/main.go
API_PORT = 8080
AIR_CMD = air
TEMP_DIR = tmp

# Ensure that clean, build, and run are treated as commands to execute, not as files or directories
.PHONY: dev run utest test help

dev:
	@echo "Starting Go Server with Air..."
	@$(AIR_CMD)

run:
	@echo "Starting Go Server..."
	go run $(API_FILE) $(API_PORT)

clean:
	@echo "Cleaning up temporary files..."
	rm -rf $(TEMP_DIR)

utest:
	@echo "Running unit tests..."
	go test -v ./tests/unit

help:
	@echo "Available commands:"
	@echo "  make dev   - Start the Go server with Air for live reloading"
	@echo "  make run   - Start the Go server"
	@echo "  make clean - Clean up temporary files"
