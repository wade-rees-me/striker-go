# Stage 1: Build the Go application
FROM golang:1.22 as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with CGO disabled for static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o striker .

# Stage 2: Create a small image for the application
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/striker .

# Expose the application port (if any)
# EXPOSE 8080

# Command to run the application
CMD ["./striker", "-queue"]

