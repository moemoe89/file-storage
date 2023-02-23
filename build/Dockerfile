# Stage 1: Build the Go binary
FROM golang:1.19-alpine AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o file-storage ./cmd/main.go

# Stage 2: Copy the binary to a minimal Alpine image
FROM alpine:3.17 As server
WORKDIR /app
COPY --from=builder /app/file-storage /app
RUN chmod +x ./file-storage
CMD ["./file-storage"]
