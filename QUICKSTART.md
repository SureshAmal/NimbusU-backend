# NimbusU User Service - Quick Start Guide

This guide will help you get the User Service up and running in under 5 minutes.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL client (psql)
- Make (optional, but recommended)

## Quick Start (Using Makefile)

### 1. Clone and Navigate
```bash
cd /home/suresh/dev/NimbusU/NimbusU-backend
```

### 2. Create Environment File
```bash
make env
# Edit .env with your configuration
```

### 3. Complete Setup (All-in-One)
```bash
make setup
# This will:
# - Start Docker infrastructure (Kafka, PostgreSQL, Redis)
# - Run database migrations
# - Seed the database with default users
```

Infrastructure Ports:
- PostgreSQL: 5433 (to avoid conflict with local DB)
- Redis: 6380 (to avoid conflict with local Redis)
- Kafka: 9094-9096 (external ports)

### 4. Run the Service
```bash
make run
```

The service should now be running on `http://localhost:8081`

### 5. Test It
```bash
curl http://localhost:8081/health
```

## Manual Setup (Without Make)

### 1. Create .env File
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start Infrastructure
```bash
cd kafka
docker-compose up -d
cd ..
```

Wait for services to be ready (about 30 seconds).

### 3. Create Database
```bash
# Connect to PostgreSQL (Port 5433)
psql -U postgres -h localhost -p 5433

# Create database
CREATE DATABASE user_service_db;
\q
```

### 4. Install Migration Tool
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 5. Run Migrations
```bash
cd services/user-service
migrate -path migrations \
  -database "postgres://nimbusu:password@localhost:5433/user_service_db?sslmode=disable" \
  up
```

### 6. Seed Database (Optional)
```bash
psql "postgres://nimbusu:password@localhost:5433/user_service_db?sslmode=disable" \
  -f migrations/seed.sql
```

### 7. Build and Run
```bash
go build -o bin/user-service cmd/main.go
./bin/user-service
```

Or run directly:
```bash
go run cmd/main.go
```

## Testing the API

### Health Check
```bash
curl http://localhost:8081/health
```

Expected response:
```json
{
  "success": true,
  "message": "User service is healthy",
  "data": {
    "service": "user-service",
    "status": "healthy"
  }
}
```

### Login (If you seeded the database)
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@nimbusu.edu",
    "password": "Admin@123"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "...",
      "email": "admin@nimbusu.edu",
      "role": "admin",
      "status": "active"
    },
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "expires_in": 3600
  }
}
```

### Get Current User (Authenticated)
```bash
# Use the access_token from login response
curl http://localhost:8081/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### List Users (Admin Only)
```bash
curl http://localhost:8081/admin/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Default Credentials (After Seeding)

### Admin
- Email: `admin@nimbusu.edu`
- Password: `Admin@123`

### Faculty
- Email: `john.doe@nimbusu.edu`
- Password: `Faculty@123`

- Email: `jane.smith@nimbusu.edu`
- Password: `Faculty@123`

### Students
- Email: `alice.johnson@student.nimbusu.edu`
- Password: `Student@123`

- Email: `bob.williams@student.nimbusu.edu`
- Password: `Student@123`

- Email: `carol.davis@student.nimbusu.edu`
- Password: `Student@123`

### Staff
- Email: `david.brown@nimbusu.edu`
- Password: `Staff@123`

## Makefile Commands

```bash
make help              # Show all available commands
make install           # Install dependencies
make build             # Build the service
make run               # Run the service
make test              # Run tests
make docker-up         # Start infrastructure
make docker-down       # Stop infrastructure
make migrate-up        # Run migrations
make migrate-down      # Rollback migration
make seed              # Seed database
make clean             # Clean build artifacts
make setup             # Complete setup (docker + migrate + seed)
```

## Troubleshooting

### Port Already in Use
If port 8081 is already in use, change it in your `.env` file:
```env
PORT=8082
```

### Database Connection Failed
Make sure PostgreSQL is running and credentials in `.env` are correct:
```bash
docker-compose -f kafka/docker-compose.yaml ps
```

### Kafka Connection Failed
Verify Kafka is running:
```bash
docker-compose -f kafka/docker-compose.yaml logs kafka
```

### Migration Errors
Reset and rerun migrations:
```bash
make migrate-reset
```

### Module Errors
Tidy up dependencies:
```bash
make tidy
```

## Project Structure

```
services/user-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── domain/                 # Business entities and interfaces
│   ├── dto/                    # Request/Response DTOs
│   ├── handler/http/           # HTTP handlers
│   ├── repository/postgres/    # Database layer
│   └── service/                # Business logic
├── migrations/                 # Database migrations
│   └── seed.sql                # Seed data
├── tools/                      # Utility scripts
│   └── generate_hashes.go      # Password hash generator
└── bin/                        # Built binaries
```

## Next Steps

1. **Add More Users**: Use the `/admin/users` endpoint to create users
2. **Implement Other Services**: Course Service, Enrollment Service, etc.
3. **Add Tests**: Write unit and integration tests
4. **Deploy**: Create Dockerfile and deploy to production
5. **API Documentation**: Generate Swagger/OpenAPI documentation
6. **Monitoring**: Add Prometheus metrics and alerting

## Development

### Live Reload
Install `air` for live reload during development:
```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Code Formatting
```bash
make fmt
```

### Linting
```bash
# Install golangci-lint first
make lint
```

### Run Tests with Coverage
```bash
make test-coverage
```

## API Documentation

Full API documentation is available in the main README.md file.

Key endpoints:
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout
- `GET /users/me` - Get current user
- `GET /admin/users` - List all users (admin)
- `POST /admin/users` - Create user (admin)

## Support

For issues or questions:
1. Check the main README.md
2. Review the code documentation
3. Check Docker logs: `make docker-logs`
4. Review application logs

## License

Copyright © 2024 NimbusU
