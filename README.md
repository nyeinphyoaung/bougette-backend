# Bougette Backend

A budget tracking backend application built with Go, Echo framework, MySQL, and Redis.

## Features

- User authentication with JWT
- Budget management
- Category management
- Wallet tracking
- Notifications
- File uploads
- WebSocket support
- Redis caching

## Prerequisites

- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

## Installation

1. Clone the repository
```bash
git clone <repository-url>
cd bougette-backend
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_IP=localhost
SERVER_PORT=9984

# Database Configuration
DB_NAME=bougette
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=your_username
DB_PASSWORD=your_password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Email Configuration
MAIL_SENDER=your_email@example.com
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=your_email@example.com
MAIL_PASSWORD=your_email_password

# AWS Configuration (for file uploads)
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BUCKET_NAME=your_bucket_name

# App Configuration
VIA_APP_NAME=Bougette
```

4. Start MySQL and Redis
```bash
# MySQL
mysql -u root -p < create_database.sql  # If you have a schema file

# Redis
redis-server  # or use Docker: docker run -p 6379:6379 redis
```

5. Run the application
```bash
go run cmd/server.go
```

## Redis Integration

The application now includes Redis support for caching and session management.

### Environment Variables

Add these to your `.env` file:

```env
REDIS_HOST=localhost    # Redis host
REDIS_PORT=6379         # Redis port
REDIS_PASSWORD=         # Redis password (leave empty if not set)
REDIS_DB=0              # Redis database number
```

### Using Redis in Your Code

#### 1. Using the Redis Utility

```go
import "bougette-backend/utilities"

func SomeFunction() {
    redis := utilities.NewRedisClient()
    
    // Set a value
    redis.SetString("key", "value", time.Minute*5)
    
    // Get a value
    value, _ := redis.GetString("key")
    
    // Set complex objects
    redis.Set("user:1", userData, time.Hour)
    
    // Get complex objects
    var userData User
    redis.GetObject("user:1", &userData)
}
```

#### 2. Using the Direct Redis Client

```go
import "bougette-backend/common"

func SomeFunction() {
    redisClient, _ := common.GetRedis()
    
    // Use the Redis client directly for advanced operations
    ctx := context.Background()
    redisClient.Set(ctx, "key", "value", time.Minute)
}
```

### Redis Use Cases in This Application

- **Caching**: Cache frequently accessed data like user profiles
- **Session Management**: Store JWT tokens and session data
- **Rate Limiting**: Implement rate limiting for API endpoints
- **Real-time Features**: Store WebSocket connection metadata
- **Temporary Storage**: Store temporary data like password reset tokens

See `examples/redis-usage.go` for comprehensive examples.

## Project Structure

```
bougette-backend/
├── cmd/
│   └── server.go          # Application entry point
├── configs/
│   ├── config.go           # Configuration and database connections
│   └── websocket.go        # WebSocket configuration
├── controllers/            # Request handlers
├── services/               # Business logic
├── repositories/           # Data access layer
├── models/                 # Database models
├── middlewares/            # HTTP middlewares
├── routes/                  # Route definitions
├── dtos/                    # Data transfer objects
├── helper/                  # Utility functions
├── utilities/              # Utility modules (including Redis)
├── common/                 # Shared utilities
├── validation/             # Input validation
├── seeder/                 # Database seeders
└── examples/               # Example code (including Redis usage)
```

## API Endpoints

The application provides various endpoints for:
- User authentication and registration
- Password reset
- Budget management
- Category management
- Wallet operations
- Notification handling
- File uploads
- WebSocket connections

## Development

### Running with Hot Reload

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run air
air
```

### Running Tests

```bash
go test ./...
```

## License

See LICENSE file for details.

