# WVP-PRO-GO

GB28181 Video Surveillance Platform - Go Implementation

Based on [WVP-PRO v2.7.4](https://github.com/648540858/wvp-GB28181-pro), rewritten in Go.

## Features

- **GB28181-2016 Protocol**: Full SIP-based device management, live streaming, playback, PTZ control
- **JT/T 1078 Protocol**: Transport industry video terminal support
- **ZLMediaKit Integration**: Stream processing, recording, snapshot
- **REST API**: ~150+ endpoints compatible with original Java API
- **Multi-database**: MySQL and PostgreSQL support
- **Redis**: Session management, caching, cluster communication
- **JWT Authentication**: Secure API access

## Quick Start

### Prerequisites

- Go 1.22+
- MySQL 5.7+ or PostgreSQL 12+
- Redis 6+
- ZLMediaKit (for media streaming)

### Linux Build

```bash
# Install Go dependencies
make deps

# Build
make build

# Run (with config file)
./wvp ./configs/config.yaml
```

### Windows Build

```cmd
:: Install Go 1.22+
:: Build
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o wvp.exe .\cmd\wvp\

:: Run
wvp.exe .\configs\config.yaml
```

### Docker

```bash
# Build image
make docker-build

# Run
docker run -d \
  -p 18080:18080 \
  -p 8116:8116/tcp \
  -p 8116:8116/udp \
  -v $(pwd)/configs:/app/configs \
  --name wvp-go \
  wvp-pro-go:latest
```

## Configuration

Edit `configs/config.yaml`:

```yaml
server:
  port: 18080

database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  dbname: wvp
  username: root
  password: your_password

redis:
  host: 127.0.0.1
  port: 6379

sip:
  port: 8116
  domain: "4101050000"
  id: "41010500002000000001"
  password: "12345678"

media:
  ip: "192.168.1.10"
  http-port: 9092
  secret: "your-zlm-secret"
```

## API Documentation

After starting the server, access Swagger UI at:
```
http://localhost:18080/swagger/index.html
```

Generate docs:
```bash
make swag
```

## Project Structure

```
cmd/wvp/                    # Application entry point
internal/
  config/                   # Configuration (Viper)
  model/                    # GORM models
  database/                 # Database connection
  redis/                    # Redis client
  sip/                      # GB28181 SIP protocol
    xml/                    # GB28181 XML message types
  zlm/                      # ZLMediaKit integration
  handler/                  # HTTP handlers (Gin)
  service/                  # Business logic
  middleware/               # HTTP middleware
  event/                    # Event bus
  task/                     # Scheduled tasks
  utils/                    # Utilities
pkg/gb28181/                # Reusable GB28181 types
configs/                    # Configuration files
sql/                        # Database schemas
```

## API Compatibility

All API endpoints maintain the same contract as the Java version:
- Same paths, methods, request/response formats
- Same error codes (0=success, 100=failure, 400=param error, etc.)
- Same JSON field names (camelCase)
- Same pagination format (PageInfo)

## Tech Stack

| Component | Technology |
|-----------|------------|
| HTTP Framework | Gin |
| ORM | GORM |
| SIP Protocol | gosip |
| Redis | go-redis/v9 |
| Configuration | Viper |
| Logging | Zap |
| JWT | golang-jwt |

## License

MIT License - Same as original WVP-PRO

## Notes

This is a skeleton/foundation implementation. The REST API handlers provide the complete endpoint structure with stub responses. The SIP protocol layer provides the infrastructure. Business logic in the service layer connects handlers with the SIP/ZLM components.

The project structure and API contracts are designed to be 100% compatible with the Java WVP-PRO v2.7.4 frontend.
