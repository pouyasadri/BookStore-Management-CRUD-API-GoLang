FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build a static binary; avoids requiring gcc inside the builder image.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/main

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .

EXPOSE 8080

CMD ["./main"]
