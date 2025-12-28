# NimbusU Backend - Microservices Architecture

University Content Management System built with **Go microservices**, **PostgreSQL**, **Redis**, **Kafka**, and **Gin framework**.

## 🏗️ Architecture Overview

```
NimbusU-backend/
├── shared/                          # Shared libraries for all services
│   ├── config/                      # Configuration management
│   ├── database/                    # PostgreSQL & Redis connection pools
│   ├── kafka/                       # Kafka producer/consumer wrappers
│   ├── logger/                      # Zap structured logging
│   ├── middleware/                  # HTTP middleware (auth, CORS, logging, rate limit)
│   ├── models/                      # Kafka event schemas
│   └── utils/                       # JWT, password hashing, API responses
│
├── services/                        # Microservices
│   └── user-service/                # User management & authentication service
│       ├── cmd/                     # Application entry points
│       ├── internal/                # Private application code
│       │   ├── domain/              # Business entities & interfaces
│       │   ├── dto/                 # Data Transfer Objects
│       │   ├── handler/             # HTTP handlers & routes
│       │   ├── repository/          # Data access layer (PostgreSQL)
│       │   └── service/             # Business logic layer
│       └── migrations/              # Database migration files
│
└── kafka/                           # Kafka configuration
    └── docker-compose.yaml          # Local Kafka setup
```

## 📦 Microservices

### ✅ User Service (Port 8081)
**Status:** Fully Implemented

**Responsibilities:**
- User authentication (login, logout, JWT tokens)
- User management (CRUD operations)
- Role-based access control (RBAC)
- Password management (reset, change)
- Session management
- User activity logging
- Bulk user import

**API Endpoints:**
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/password/change` - Change password
- `POST /auth/password/reset-request` - Request password reset
- `POST /auth/password/reset` - Reset password with token
- `GET /auth/sessions` - Get active sessions
- `DELETE /auth/sessions` - Revoke all sessions
- `GET /users/me` - Get current user
- `PUT /users/me` - Update current user
- `POST /admin/users` - Create user (admin)
- `GET /admin/users` - List users (admin)
- `GET /admin/users/:id` - Get user by ID (admin)
- `PUT /admin/users/:id` - Update user (admin)
- `DELETE /admin/users/:id` - Delete user (admin)
- `POST /admin/users/:id/activate` - Activate user (admin)
- `POST /admin/users/:id/suspend` - Suspend user (admin)
- `POST /admin/users/bulk-import` - Bulk import users (admin)

**Events Published:**
- `user.events` - USER_CREATED, USER_UPDATED, USER_DELETED, USER_ACTIVATED, USER_SUSPENDED
- `auth.events` - LOGIN_SUCCESS, LOGIN_FAILED, LOGOUT, PASSWORD_CHANGED

**Database Tables:**
- `users` - Core user identity
- `user_profiles` - Extended profile information
- `roles` - User roles (admin, faculty, student, staff)
- `permissions` - Granular permissions
- `role_permissions` - Role-permission mapping
- `user_activity_logs` - Audit trail
- `password_reset_tokens` - Password reset flow
- `active_sessions` - Session tracking

### 🔜 Planned Services

- **Content Service** (Port 8082) - Document management, file storage
- **Course Service** (Port 8084) - Course management, enrollments
- **Timetable Service** (Port 8085) - Schedule management
- **Attendance Service** (Port 8086) - Attendance tracking
- **Notification Service** (Port 8083) - Email, SMS, push notifications
- **Announcement Service** (Port 8087) - Announcements and notices
- **Communication Service** (Port 8088) - Messaging and forums
- **Analytics Service** (Port 8089) - Reporting and analytics

## 🛠️ Technology Stack

### Core Technologies
- **Language:** Go 1.21+
- **Web Framework:** Gin (HTTP router)
- **Database:** PostgreSQL 15+ (pgx/v5 driver)
- **Cache:** Redis 7+
- **Message Broker:** Apache Kafka (Sarama client)

### Libraries
- **Logging:** Zap (structured logging)
- **Auth:** JWT (golang-jwt/jwt/v5)
- **Security:** bcrypt (password hashing)
- **Validation:** go-playground/validator
- **UUID:** google/uuid

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Apache Kafka 3.0+
- Docker & Docker Compose (optional, for local development)

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
ENV=development
PORT=8081

# PostgreSQL
DATABASE_URL=postgres://nimbusu:password@localhost:5432/user_service_db
DB_MAX_CONNECTIONS=25
DB_MIN_CONNECTIONS=5

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=user-service-group

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=3600        # 1 hour
JWT_REFRESH_TOKEN_EXPIRY=604800     # 7 days

# Logging
LOG_LEVEL=debug
```

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd NimbusU-backend
   ```

2. **Install dependencies:**
   ```bash
   # Install shared library dependencies
   cd shared && go mod download && cd ..
   
   # Install user-service dependencies
   cd services/user-service && go mod download && cd ../..
   ```

3. **Start infrastructure (Docker Compose):**
   ```bash
   # Start PostgreSQL, Redis, and Kafka
   cd kafka
   docker-compose up -d
   cd ..
   ```

4. **Run database migrations:**
   ```bash
   cd services/user-service
   # Run migrations using your preferred migration tool
   # e.g., migrate -path migrations -database "postgres://..." up
   ```

5. **Run the service:**
   ```bash
   cd services/user-service
   go run cmd/main.go
   ```

### Development

```bash
# Run tests
go test ./...

# Run with hot reload (using air)
air

# Build binary
go build -o bin/user-service cmd/main.go

# Run linter
golangci-lint run
```

## 📚 API Documentation

API documentation will be available at:
- Swagger UI: `http://localhost:8081/swagger/index.html` (planned)
- Postman Collection: `docs/postman/` (planned)

## 🔒 Security Features

- **Password Hashing:** Bcrypt with default cost
- **JWT Authentication:** Access & refresh tokens
- **Session Management:** Redis-backed with expiry
- **Rate Limiting:** 100 requests/minute per IP (configurable)
- **RBAC:** Role-based access control
- **Audit Logging:** All user actions logged
- **Input Validation:** Request validation with go-playground/validator
- **SQL Injection Prevention:** Parameterized queries with pgx

## 📊 Event-Driven Architecture

Services communicate asynchronously via Kafka:

### Topics
- `user.events` - User lifecycle events
- `auth.events` - Authentication events
- `content.events` - Content management events (planned)
- `notification.events` - Notification requests (planned)

### Event Schema
```json
{
  "event_id": "uuid",
  "event_type": "USER_CREATED",
  "timestamp": "2024-12-28T12:00:00Z",
  "service_name": "user-service",
  "user_id": "uuid",
  "email": "user@example.com",
  "metadata": {}
}
```

## 🗄️ Database Schema

Complete database schema documentation: `/database/DATABASE_SCHEMA_COMPLETE.md`

### User Service Tables
- **users** - Primary user identity table
- **user_profiles** - Extended user information
- **roles** - System roles
- **permissions** - Granular permissions
- **role_permissions** - Many-to-many mapping
- **user_activity_logs** - Audit trail
- **password_reset_tokens** - Password reset workflow
- **active_sessions** - Session tracking

## 🐳 Docker Deployment

```bash
# Build image
docker build -t nimbusu/user-service:latest services/user-service

# Run container
docker run -p 8081:8081 --env-file .env nimbusu/user-service:latest
```

## 📈 Monitoring & Observability

- **Logging:** Structured JSON logs with Zap
- **Metrics:** Prometheus metrics (planned)
- **Tracing:** OpenTelemetry (planned)
- **Health Checks:** `/health` endpoint

## 🧪 Testing

```bash
# Run unit tests
go test ./services/user-service/internal/...

# Run integration tests
go test -tags=integration ./services/user-service/...

# Generate coverage report
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 Code Standards

- **Architecture:** Clean Architecture (Domain → Service → Handler)
- **Naming:** Go conventions (camelCase for unexported, PascalCase for exported)
- **Error Handling:** Explicit error returns, wrapped errors
- **Logging:** Structured logging with context
- **Comments:** Exported functions and types must have doc comments

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Write tests
4. Run linter and tests
5. Submit a pull request

## 📄 License

[Add your license here]

## 🔗 Related Documentation

- [Database Schema](/database/DATABASE_SCHEMA_COMPLETE.md)
- [API Documentation](/database/USER_API_DOCUMENTATION.md)
- [Microservice Architecture](/database/MICROSERVICE_ARCHITECTURE.md)
- [Kafka Events](/database/KAFKA_EVENTS_DOCUMENTATION.md)
- [Requirements](/requirements.md)

## 👥 Team

[Add team information]

---

**Built with ❤️ using Go**
