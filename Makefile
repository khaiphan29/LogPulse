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
.ONESHELL:

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
	make start-containers
	sleep 5 # wait for services to be fully up
	@echo "Setting up Elasticsearch indexes..."
	@set -a; . .env; set +a; go run ./cmd/create_es_indexes/main.go $(ES_PORT)
	@echo "Creating Kafka topics..."
	@set -a; . .env; set +a; go run ./cmd/create_kafka_topics/main.go $(KAFKA_PORT)

dev:
	make start-containers
	sleep 5 # wait for services to be fully up
	@echo "Starting Go Server with Air..."
	@set -a; . .env; set +a; APP_ENV=dev $(AIR_CMD) # temporarily export all env vars from .env

run:
	make start-containers
	sleep 5 # wait for services to be fully up
	@echo "Starting Go Server..."
	@set -a; . .env; set +a; APP_ENV=dev go run $(API_FILE)

clean:
	@echo "Cleaning up temporary files..."
	rm -rf $(TEMP_DIR)

test:
	export APP_ENV=test; \
		make start-containers; \
		sleep 5; \
		echo "Running tests..."; \
		set -a; . .env; set +a; \
		gotestsum --format testname -- -v ./... \

p ?= ./...
testfile:
		export APP_ENV=test; \
		make start-containers; \
		sleep 5; \
		echo "Running tests..."; \
		set -a; . .env; set +a; \
		gotestsum $(p) \

help:
	@echo "Available commands:"
	@echo "  make provision - Set up Elasticsearch indexes and create Kafka topics"
	@echo "  make dev   - Start the Go server with Air for live reloading"
	@echo "  make run   - Start the Go server"
	@echo "  make clean - Clean up temporary files"
	@echo "  make test - Run unit tests"

