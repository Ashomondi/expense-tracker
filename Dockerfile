# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY backend/ ./backend/

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./backend

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy frontend files
COPY frontend/ ./frontend/

# JSON data files will be created by the application at runtime

# Set frontend directory path
ENV FRONTEND_DIR=./frontend/

# Expose port (Render will set PORT environment variable)
EXPOSE 8080

# Run the server
CMD ["./server"]
