# Makefile for Yato

# Check for .env file
ifneq ($(wildcard .env),)
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Ensure required environment variables are set
ifndef MAL_CLIENT_ID
$(error MAL_CLIENT_ID is not set. Please check your .env file)
endif
ifndef MAL_CLIENT_SECRET
$(error MAL_CLIENT_SECRET is not set. Please check your .env file)
endif

# Variables
BUILD_DIR=build
BINARY_NAME=yato
MAIN_FILE=main.go
ENCODED_CLIENT_ID=$(shell printf '%s' "$(MAL_CLIENT_ID)" | base64)
ENCODED_CLIENT_SECRET=$(shell printf '%s' "$(MAL_CLIENT_SECRET)" | base64)

.PHONY: build run clean

build:
	@echo "Building Yato..."
	@echo "Encoded Client ID: $(ENCODED_CLIENT_ID)"
	@echo "Encoded Client Secret: $(ENCODED_CLIENT_SECRET)"
	@go build -ldflags '-X "yato/config.encodedClientID=$(ENCODED_CLIENT_ID)" -X "yato/config.encodedClientSecret=$(ENCODED_CLIENT_SECRET)"' -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete. Binary '$(BINARY_NAME)' created."

run:
	@echo "Running Yato..."
	@go run .

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup complete."