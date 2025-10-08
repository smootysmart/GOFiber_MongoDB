FROM golang:1.25-alpine AS builder
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o server ./main.go
FROM alpine:latest
WORKDIR /
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]