# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
# Point the compiler to the exact location of your main package
RUN go build -o server ./cmd/server/main.go

# Production stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
