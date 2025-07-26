# Define path to main source files
API_FILE = cmd/api/main.go
API_PORT = 8080
AIR_CMD = air
TEMP_DIR = tmp

KAFKA_PORT ?= 9092
KAFKA_LOCAL_LOG_DIR = logs/broker1-logs
KAFKA_LOCAL_METADATA_DIR = logs/broker1-metadata

ES_PORT ?= 9200

# Ensure that clean, build, and run are treated as commands to execute, not as files or directories
.PHONY: dev run utest test setup-kafka-brokers setup-kafka-topics help

# Local Kafka setup - Deprecated
# setup-kafka-brokers:
# 	@echo "Setting up Kafka..."
# 	@echo "Formatting broker 1..."
# 	kafka-storage format --config ./configs/kafka/broker1.properties --cluster-id $(shell kafka-storage random-uuid)
# 	@echo "Starting Kafka Broker 1 on port $(KAFKA_PORT)..."
# 	kafka-server-start ./configs/kafka/broker1.properties
#
# start-kafka-broker:
# 	@echo "Starting Kafka Broker 1..."
# 	kafka-server-start ./configs/kafka/broker1.properties

list-kafka-topics:
	kafka-topics --bootstrap-server localhost:$(KAFKA_PORT) --list

start-containers:
	@echo "Starting Docker containers (DETACHED mode)..."
	docker-compose up -d

# Example: make provision KAFKA_PORT=9092
provision:
	@echo "Setting up Elasticsearch indexes..."
	go run ./cmd/create_es_indexes/main.go $(ES_PORT)
	@echo "Creating Kafka topics..."
	go run ./cmd/create_kafka_topics/main.go $(KAFKA_PORT)

dev:
	make start-containers; \
	@export APP_ENV=dev; \
	echo "Starting Go Server with Air..."; \
	@$(AIR_CMD)

run:
	export APP_ENV=dev
	@echo "Starting Go Server..."
	go build -o ./tmp/app $(API_FILE)
	./tmp/app $(API_PORT)

clean:
	@echo "Cleaning up temporary files..."
	rm -rf $(TEMP_DIR)

utest:
	@export APP_ENV=test; \
	echo "Running unit tests..."; \
	go test -v ./tests/unit

help:
	@echo "Available commands:"
	@echo "  make provision - Set up Elasticsearch indexes and create Kafka topics"
	@echo "  make dev   - Start the Go server with Air for live reloading"
	@echo "  make run   - Start the Go server"
	@echo "  make clean - Clean up temporary files"
	@echo "  make utest - Run unit tests"

