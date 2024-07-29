# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o portpecker .

# Final stage
FROM alpine:latest  

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/portpecker .

# Run the binary
CMD ["./portpecker", "/app/config.json"]