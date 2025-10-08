# ---- Stage 1: Build the binary ----
FROM golang:1.25-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the Go Fiber binary
RUN go build -o server ./main.go

# ---- Stage 2: Create minimal runtime image ----
FROM alpine:latest

WORKDIR /

# Copy binary from builder stage
COPY --from=builder /app/server .

# Expose Fiber port
EXPOSE 8080

# Run the binary
CMD ["./server"]
