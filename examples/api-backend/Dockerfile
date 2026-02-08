# GFTD Auto-Generated Dockerfile
# Type: web | Stack: go

# ===== Build Stage =====
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/example_api ./cmd/server

# ===== Production Stage =====
FROM gcr.io/distroless/static-debian12 AS runner

# Copy the binary
COPY --from=builder /app/example_api /example_api

EXPOSE 8080

ENTRYPOINT ["/example_api"]
