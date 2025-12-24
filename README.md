# User Service

This is the User Service for the E-commerce system, built with Go, Gin, GORM, and Redis.

## Architecture

- **Framework**: Gin
- **Database**: PostgreSQL (GORM)
- **Cache**: Redis
- **Config**: Fetched from Config Service via gRPC

## Setup

1. **Prerequisites**:
    - PostgreSQL
    - Redis
    - Config Service running on port 50051 (default)

2. **Configuration**:
    - Update `.env` with `CONFIG_SERVICE_URL`, `JWT_SECRET`, and `PORT`.
    - Ensure Config Service has `user_service` configuration.

3. **Run**:
    ```bash
    go mod tidy
    go run cmd/api/main.go
    ```

## APIs

- `POST /v1/user/register`: Register a new user.
- `POST /v1/user/login`: Login and get JWT.
- `GET /v1/user/profile`: Get user profile (Protected).
