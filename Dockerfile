# Build stage
FROM golang:1.23.2-alpine AS builder

# Set working directory
WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /

# Copy the binary from builder
COPY --from=builder /main .

# Create directory for persistent storage
RUN mkdir -p /pb_data

# Expose the default PocketBase port
EXPOSE 8090

# Set volume for persistent data
VOLUME ["/pb_data"]

# Run the application
CMD ["./main", "serve", "--http=0.0.0.0:8090"]
