#!/bin/bash

docker run --rm \
  -v "$(pwd)":/app \
  -w /app \
  -e MONGODB_URI="mongodb://admin:password@mongodb-server.mongodb.svc.cluster.local:27017" \
  golang:1.22.4 \
  sh -c "go run ./insert_data.go"