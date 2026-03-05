# Variables
SERVER_DIR=server
NODE_DIR=node

.PHONY: help server-lint server-lint-fix server-run node-lint node-compile node-upload node-monitor

help:
	@echo "Available commands:"
	@echo "  make server-lint    - Run golangci-lint for the server"
	@echo "  make server-run     - Run the Go server"
	@echo "  make node-compile   - Compile the ESP-IDF project"
	@echo "  make node-upload    - Upload firmware via idf.py flash"
	@echo "  make node-monitor   - Open serial monitor via idf.py monitor"

# Server
server-lint:
	cd $(SERVER_DIR) && golangci-lint run

server-lint-fix:
	cd $(SERVER_DIR) && golangci-lint run --fix

server-run:
	cd $(SERVER_DIR) && go run main.go

IDF_EXPORT = . $(HOME)/esp/esp-idf/export.sh

# Node
node-compile:
	($(IDF_EXPORT) && cd $(NODE_DIR) && idf.py build)

node-upload:
	($(IDF_EXPORT) && cd $(NODE_DIR) && idf.py flash)

node-monitor:
	($(IDF_EXPORT) && cd $(NODE_DIR) && idf.py monitor)
