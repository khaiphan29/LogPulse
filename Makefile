# Define path to main source files
API_FILE = cmd/api/main.go
API_PORT = 8080
AIR_CMD = air
TEMP_DIR = tmp

KAFKA_BROKER1_PORT = 9094
KAFKA_BROKER1_LOG_DIR = logs/broker1-logs
KAFKA_BROKER1_METADATA_DIR = logs/broker1-metadata

# Ensure that clean, build, and run are treated as commands to execute, not as files or directories
.PHONY: dev run utest tes setup-kafka-brokers setup-kafka-topics help

setup-kafka-brokers:
	@echo "Setting up Kafka..."
	@echo "Formatting broker 1..."
	kafka-storage format --config ./configs/kafka/broker1.properties --cluster-id $(shell kafka-storage random-uuid)
	@echo "Starting Kafka Broker 1 on port $(KAFKA_BROKER1_PORT)..."
	kafka-server-start ./configs/kafka/broker1.properties

setup-kafka-topics:
	@echo "Creating Kafka topics..."
	go run ./deployments/create_kafka_topics.go

start-kafka-broker1:
	@echo "Starting Kafka Broker 1..."
	kafka-server-start ./configs/kafka/broker1.properties

show-kafka-topics:
	kafka-topics --bootstrap-server localhost:$(KAFKA_BROKER1_PORT) --list

dev:
	export APP_ENV=development
	@echo "Starting Go Server with Air..."
	@$(AIR_CMD)

run:
	export APP_ENV=development
	@echo "Starting Go Server..."
	go build -o ./tmp/app $(API_FILE)
	./tmp/app $(API_PORT)

clean:
	@echo "Cleaning up temporary files..."
	rm -rf $(TEMP_DIR)

utest:
	export APP_ENV=test
	@echo "Running unit tests..."
	go test -v ./tests/unit

help:
	@echo "Available commands:"
	@echo "  make dev   - Start the Go server with Air for live reloading"
	@echo "  make run   - Start the Go server"
	@echo "  make clean - Clean up temporary files"
