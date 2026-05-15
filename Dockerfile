FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wvp ./cmd/wvp/

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and config
COPY --from=builder /app/wvp /app/wvp
COPY --from=builder /app/configs /app/configs

# Expose ports
EXPOSE 18080 8116/tcp 8116/udp

# Set timezone
ENV TZ=Asia/Shanghai

# Run
ENTRYPOINT ["./wvp"]
CMD ["./configs/config.yaml"]
