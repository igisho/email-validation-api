# Stage 1: Build the Go binary
FROM golang:1.18-alpine AS builder

WORKDIR /eva

# Install dependencies
RUN apk add --no-cache git tzdata

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy the binary and config directory from the builder stage
COPY --from=builder /eva/main .
COPY --from=builder /eva/config ./config

# Create logs directory
RUN mkdir logs

# Expose the port the app runs on
EXPOSE 80

# Define healthcheck
HEALTHCHECK CMD wget --spider --quiet --tries=1 --timeout=5 http://localhost:80/ || exit 1

# Command to run the binary
CMD ["./main"]
