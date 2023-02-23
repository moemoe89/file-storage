#!/bin/bash

# app config
export APP_ENV=dev
export SERVER_PORT=8080

# minio config
export MINIO_HOST=localhost:9000
export MINIO_ACCESS_KEY_ID=minioadmin
export MINIO_SECRET_ACCESS_KEY=minioadmin

# tracing config
export OTEL_AGENT=http://localhost:14268/api/traces

go build -o main ./cmd/main.go && ./main
