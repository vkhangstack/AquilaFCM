# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set environment variables
ENV GOGC=50

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aquilafcm .

# Stage 2: Create a lightweight container with the built binary
FROM alpine:latest

# Install curl and htop
RUN apk --no-cache add curl htop git nano

# Set the working directory
WORKDIR /app
RUN mkdir -p config

# Copy the binary from the builder stage
COPY --from=builder /app/aquilafcm .
COPY --from=builder /app/config/ config/

# Ensure the binary is executable
RUN chmod +x /app/aquilafcm

# Expose the port the service listens on
EXPOSE 8080

# Debugging: List files to verify if main exists
RUN ls -l /app/

# Command to run the application
CMD ["./aquilafcm"]
