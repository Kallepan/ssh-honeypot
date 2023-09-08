#!/bin/bash

export $(grep -v '^#' .dev.env | xargs)

go run cmd/ssh-honeypot/main.go
