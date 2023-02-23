#!/bin/bash

# app config
export APP_ENV=dev
export SERVER_PORT=8080

# tracing config
export OTEL_AGENT=http://localhost:14268/api/traces

go build -o main ./cmd/main.go && ./main
