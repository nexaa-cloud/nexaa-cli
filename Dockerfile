# Build stage
FROM golang:1.23.0-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info
ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT_SHA
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/nexaa-cloud/nexaa-cli/cmd.Version=${VERSION} -X github.com/nexaa-cloud/nexaa-cli/cmd.BuildDate=${BUILD_DATE} -X github.com/nexaa-cloud/nexaa-cli/cmd.CommitSHA=${COMMIT_SHA}" \
    -o nexaa-cli .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S nexaa && \
    adduser -u 1001 -S nexaa -G nexaa

# Set working directory
WORKDIR /home/nexaa

# Copy binary from builder stage
COPY --from=builder /app/nexaa-cli /usr/local/bin/nexaa-cli

# Make sure binary is executable
RUN chmod +x /usr/local/bin/nexaa-cli

# Switch to non-root user
USER nexaa

# Set entrypoint
ENTRYPOINT ["nexaa-cli"]

# Default command shows help
CMD ["--help"]