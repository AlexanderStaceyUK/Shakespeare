# Build stage
FROM golang:alpine AS builder

# Install git for go modules (if needed)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy only the Go source file
COPY shakespeare.go .

# Download dependencies and compile the Go program
RUN go mod init shakespeare && go mod tidy
RUN go build -o shakespeare

# Final image (runtime only)
FROM alpine:latest

# Install certs for HTTPS access if needed
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder
COPY --from=builder /app/shakespeare /shakespeare

# Run the app
ENTRYPOINT ["/shakespeare"]
