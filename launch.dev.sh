#!/bin/bash

# Set environment variables from .env
export $(grep -v '^#' ./.env | xargs)

# Clean up go.mod
go mod tidy

# Install dependencies from go.mod
go mod download

# Run the application
go run main.go