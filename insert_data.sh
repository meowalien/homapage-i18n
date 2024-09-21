#!/bin/bash

docker run --rm \
  -v "$(pwd)":/app \
  -w /app \
  golang:1.22.4 \
  sh -c "go run ./insert_data.go"