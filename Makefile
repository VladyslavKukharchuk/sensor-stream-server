# Variables
SERVER_DIR=server
NODE_DIR=node
FQBN=esp32:esp32:esp32c6

.PHONY: server-lint server-run node-lint node-compile

# Server
server-lint:
	cd $(SERVER_DIR) && golangci-lint run

server-lint-fix:
	cd $(SERVER_DIR) && golangci-lint run --fix

server-run:
	cd $(SERVER_DIR) && go run main.go

# Node
node-lint:
	arduino-lint --path $(NODE_DIR) --recursive --project-type sketch

node-compile:
	arduino-cli compile --fqbn $(FQBN) $(NODE_DIR)/node.ino
