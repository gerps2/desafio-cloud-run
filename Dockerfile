# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk --no-cache add git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Install Wire for dependency injection
RUN go install github.com/google/wire/cmd/wire@latest

# Copy source code
COPY . .

# Generate Wire code
RUN wire ./cmd/api

# Build the application with optimizations for Cloud Run
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/api

# Final stage - minimal distroless image for Cloud Run
FROM gcr.io/distroless/static-debian12:nonroot

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder stage
COPY --from=builder /app/main /app/main

# Set working directory
WORKDIR /app

# Use non-root user (already set in distroless image)
USER nonroot:nonroot

# Expose port (Cloud Run uses PORT env var)
EXPOSE 8080

# Run the application
ENTRYPOINT ["./main"]
