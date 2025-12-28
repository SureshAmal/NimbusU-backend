# NimbusU - Microservice Architecture Document

## Overview

NimbusU is a University Content Management System built using a microservices architecture with:
- **Go (Golang)** - Backend services
- **Apache Kafka** - Event-driven communication
- **PostgreSQL** - Primary relational database
- **MongoDB** - Document/file storage
- **Redis** - Caching and session management

---

## Architecture Diagram

```
                                    ┌─────────────────┐
                                    │   Frontend      │
                                    │ (SvelteKit/Web) │
                                    └────────┬────────┘
                                             │
                                             ▼
                              ┌──────────────────────────┐
                              │      API Gateway         │
                              │    (Go + Chi Router)     │
                              │  - Auth/JWT Validation   │
                              │  - Rate Limiting         │
                              │  - Request Routing       │
                              └──────────┬───────────────┘
                                         │
           ┌─────────────────────────────┼─────────────────────────────┐
           │                             │                             │
           ▼                             ▼                             ▼
┌──────────────────┐        ┌──────────────────┐        ┌──────────────────┐
│   User Service   │        │ Content Service  │        │Notification Svc  │
│                  │        │                  │        │                  │
│ - Auth/Login     │        │ - Upload/Download│        │ - Email/SMS      │
│ - User CRUD      │        │ - Metadata       │        │ - Push Notifs    │
│ - Profile Mgmt   │        │ - Permissions    │        │ - Templates      │
│ - RBAC           │        │ - Search         │        │ - Queue Process  │
└────────┬─────────┘        └────────┬─────────┘        └────────┬─────────┘
         │                           │                           │
         │                           │                           │
         ▼                           ▼                           ▼
┌──────────────────┐        ┌──────────────────┐        ┌──────────────────┐
│   PostgreSQL     │        │   PostgreSQL     │        │   PostgreSQL     │
│  (User DB)       │        │  (Content DB)    │        │ (Notification DB)│
└──────────────────┘        └────────┬─────────┘        └──────────────────┘
                                     │
                                     ▼
                            ┌──────────────────┐
                            │     MongoDB      │
                            │  (File Storage)  │
                            │    GridFS        │
                            └──────────────────┘

           ┌─────────────────────────────────────────────────────────┐
           │                                                         │
           ▼                                                         ▼
┌──────────────────┐        ┌──────────────────┐        ┌──────────────────┐
│Timetable Service │        │Attendance Service│        │ Course Service   │
│                  │        │                  │        │                  │
│ - Schedule Mgmt  │        │ - Mark Attendance│        │ - Course CRUD    │
│ - Room Allocation│        │ - Reports        │        │ - Enrollments    │
│ - Conflict Check │        │ - Analytics      │        │ - Faculty Assign │
└────────┬─────────┘        └────────┬─────────┘        └────────┬─────────┘
         │                           │                           │
         ▼                           ▼                           ▼
┌──────────────────┐        ┌──────────────────┐        ┌──────────────────┐
│   PostgreSQL     │        │   PostgreSQL     │        │   PostgreSQL     │
│ (Timetable DB)   │        │ (Attendance DB)  │        │   (Course DB)    │
└──────────────────┘        └──────────────────┘        └──────────────────┘


                    ┌───────────────────────────────────┐
                    │           Apache Kafka            │
                    │  ┌─────────────────────────────┐  │
                    │  │  Topics:                    │  │
                    │  │  - user.events              │  │
                    │  │  - content.events           │  │
                    │  │  - notification.commands    │  │
                    │  │  - attendance.events        │  │
                    │  │  - timetable.events         │  │
                    │  │  - analytics.events         │  │
                    │  └─────────────────────────────┘  │
                    └───────────────────────────────────┘
                                     │
                                     ▼
                    ┌───────────────────────────────────┐
                    │             Redis                 │
                    │  - Session Store                  │
                    │  - Cache Layer                    │
                    │  - Rate Limiting                  │
                    │  - Real-time Presence             │
                    └───────────────────────────────────┘
```

---

## Microservices Overview

### 1. API Gateway Service

**Purpose:** Single entry point for all client requests.

**Responsibilities:**
- JWT validation and authentication
- Request routing to appropriate services
- Rate limiting and throttling
- Request/response logging
- CORS handling
- SSL termination

**Technology Stack:**
- Go 1.21+
- Chi Router
- Redis (rate limiting)

**Port:** 8080

---

### 2. User Service

**Purpose:** Manages user identity, authentication, and authorization.

**Responsibilities:**
- User registration and management
- Authentication (login/logout)
- JWT token generation and refresh
- Password management
- Role-based access control (RBAC)
- User profile management
- Activity logging

**Database:** PostgreSQL (user_service_db)

**Tables Owned:**
- users
- roles
- permissions
- role_permissions
- user_profiles
- user_activity_logs
- password_reset_tokens
- active_sessions

**Port:** 8081

**Kafka Topics:**
- Publishes to: `user.events`, `auth.events`
- Subscribes to: None

---

### 3. Content Service

**Purpose:** Manages educational content and media files.

**Responsibilities:**
- File upload/download
- Content metadata management
- Content versioning
- Folder organization
- Content permissions
- Search and filtering
- Integration with MongoDB for file storage

**Databases:**
- PostgreSQL (content_service_db) - Metadata
- MongoDB (content_files) - Binary files

**Tables Owned:**
- content_metadata
- content_course_mapping
- content_categories
- content_category_mapping
- content_tags
- content_tag_mapping
- content_folders
- content_folder_items
- content_permissions
- content_versions
- content_assignments
- content_access_logs
- content_delivery_status
- content_engagement_metrics
- content_bookmarks

**Port:** 8082

**Kafka Topics:**
- Publishes to: `content.events`
- Subscribes to: `user.events` (for permission updates)

---

### 4. Notification Service

**Purpose:** Manages all notification channels.

**Responsibilities:**
- Email notifications (SMTP/SendGrid)
- SMS notifications (Twilio)
- Push notifications (FCM/APNS)
- In-app notifications
- Notification templates
- Queue processing with retries
- Notification preferences

**Database:** PostgreSQL (notification_service_db)

**Tables Owned:**
- notification_templates
- notifications
- notification_queue
- notification_preferences
- notification_statistics

**Port:** 8083

**Kafka Topics:**
- Publishes to: None
- Subscribes to: `notification.commands`, `user.events`, `content.events`, `attendance.events`

---

### 5. Course Service

**Purpose:** Manages academic courses and enrollments.

**Responsibilities:**
- Department management
- Program management
- Subject and course management
- Faculty-course assignments
- Student enrollments
- Academic calendar
- Semester management

**Database:** PostgreSQL (course_service_db)

**Tables Owned:**
- departments
- programs
- subjects
- courses
- course_prerequisites
- course_corequisites
- faculty_courses
- course_enrollments
- semesters
- academic_calendar
- faculties (profile)
- students (profile)

**Port:** 8084

**Kafka Topics:**
- Publishes to: `course.events`
- Subscribes to: `user.events`

---

### 6. Timetable Service

**Purpose:** Manages class schedules and room allocations.

**Responsibilities:**
- Timetable creation and management
- Room allocation
- Time slot management
- Conflict detection and resolution
- Schedule change requests
- Timetable publishing

**Database:** PostgreSQL (timetable_service_db)

**Tables Owned:**
- rooms
- time_slots
- timetables
- timetable_entries
- timetable_conflicts
- schedule_change_requests

**Port:** 8085

**Kafka Topics:**
- Publishes to: `timetable.events`
- Subscribes to: `course.events`

---

### 7. Attendance Service

**Purpose:** Tracks student attendance.

**Responsibilities:**
- Mark attendance
- Attendance corrections
- Attendance reports
- Attendance summary calculation
- Low attendance alerts

**Database:** PostgreSQL (attendance_service_db)

**Tables Owned:**
- attendance_records
- attendance_corrections
- attendance_summary

**Port:** 8086

**Kafka Topics:**
- Publishes to: `attendance.events`
- Subscribes to: `timetable.events`, `course.events`

---

### 8. Announcement Service

**Purpose:** Manages university announcements.

**Responsibilities:**
- Create/publish announcements
- Target audience selection
- Scheduled publishing
- Announcement read tracking
- Priority announcements

**Database:** PostgreSQL (announcement_service_db)

**Tables Owned:**
- announcements
- announcement_targets
- announcement_reads

**Port:** 8087

**Kafka Topics:**
- Publishes to: `announcement.events`
- Subscribes to: `user.events`

---

### 9. Communication Service

**Purpose:** Manages messaging and discussions.

**Responsibilities:**
- Direct messaging
- Group conversations
- Discussion forums
- Office hours scheduling
- Message notifications

**Database:** PostgreSQL (communication_service_db)

**Tables Owned:**
- conversations
- conversation_participants
- messages
- message_reads
- discussion_forums
- forum_threads
- forum_posts
- office_hours
- office_hour_bookings

**Port:** 8088

**Kafka Topics:**
- Publishes to: `communication.events`
- Subscribes to: `user.events`, `course.events`

---

### 10. Analytics Service

**Purpose:** Aggregates and analyzes system data.

**Responsibilities:**
- Real-time active user monitoring
- Usage analytics
- Content engagement metrics
- Performance dashboards
- Report generation

**Database:** PostgreSQL (analytics_service_db) + TimescaleDB for time-series

**Tables Owned:**
- user_activity_tracking
- system_metrics

**Port:** 8089

**Kafka Topics:**
- Publishes to: `analytics.events`
- Subscribes to: All event topics (consumer group)

---

## Go Project Structure

Each microservice follows a consistent structure:

```
service-name/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── domain/
│   │   ├── entity.go           # Domain entities
│   │   ├── repository.go       # Repository interfaces
│   │   └── service.go          # Domain service interfaces
│   ├── handler/
│   │   ├── http/
│   │   │   ├── handler.go      # HTTP handlers
│   │   │   ├── middleware.go   # HTTP middleware
│   │   │   └── routes.go       # Route definitions
│   │   └── kafka/
│   │       ├── consumer.go     # Kafka consumers
│   │       └── producer.go     # Kafka producers
│   ├── repository/
│   │   ├── postgres/
│   │   │   └── repository.go   # PostgreSQL implementation
│   │   └── mongodb/
│   │       └── repository.go   # MongoDB implementation
│   ├── service/
│   │   └── service.go          # Business logic implementation
│   └── dto/
│       ├── request.go          # Request DTOs
│       └── response.go         # Response DTOs
├── pkg/
│   ├── validator/
│   │   └── validator.go        # Input validation
│   └── errors/
│       └── errors.go           # Custom error types
├── migrations/
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── .env.example
```

---

## Shared Libraries

Located in `shared/` directory:

```
shared/
├── config/
│   └── config.go               # Common configuration loading
├── kafka/
│   ├── producer.go             # Kafka producer wrapper
│   ├── consumer.go             # Kafka consumer wrapper
│   └── events.go               # Event type definitions
├── logger/
│   └── logger.go               # Structured logging (Zap)
├── middleware/
│   ├── auth.go                 # JWT authentication middleware
│   ├── cors.go                 # CORS middleware
│   ├── logging.go              # Request logging middleware
│   └── ratelimit.go            # Rate limiting middleware
├── database/
│   ├── postgres.go             # PostgreSQL connection pool
│   ├── mongodb.go              # MongoDB connection
│   └── redis.go                # Redis client
├── models/
│   └── events.go               # Shared event structures
└── utils/
    ├── hash.go                 # Password hashing (bcrypt)
    ├── jwt.go                  # JWT utilities
    └── response.go             # Standard API response
```

---

## Kafka Configuration

### Topics

| Topic Name | Partitions | Retention | Description |
|------------|------------|-----------|-------------|
| user.events | 6 | 7 days | User lifecycle events |
| auth.events | 6 | 7 days | Authentication events |
| content.events | 12 | 7 days | Content lifecycle events |
| notification.commands | 12 | 1 day | Notification requests |
| course.events | 6 | 7 days | Course management events |
| timetable.events | 3 | 7 days | Timetable updates |
| attendance.events | 6 | 7 days | Attendance events |
| announcement.events | 3 | 7 days | Announcement events |
| communication.events | 6 | 3 days | Messaging events |
| analytics.events | 12 | 1 day | Analytics data |

### Event Schema

```go
// shared/models/events.go

package models

import "time"

type BaseEvent struct {
    EventID     string    `json:"event_id"`
    EventType   string    `json:"event_type"`
    Timestamp   time.Time `json:"timestamp"`
    ServiceName string    `json:"service_name"`
    Version     string    `json:"version"`
}

// User Events
type UserCreatedEvent struct {
    BaseEvent
    Payload UserCreatedPayload `json:"payload"`
}

type UserCreatedPayload struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    CreatedBy string `json:"created_by"`
}

type UserUpdatedEvent struct {
    BaseEvent
    Payload UserUpdatedPayload `json:"payload"`
}

type UserUpdatedPayload struct {
    UserID    string                 `json:"user_id"`
    Changes   map[string]interface{} `json:"changes"`
    UpdatedBy string                 `json:"updated_by"`
}

// Notification Commands
type SendNotificationCommand struct {
    BaseEvent
    Payload NotificationPayload `json:"payload"`
}

type NotificationPayload struct {
    UserID           string                 `json:"user_id"`
    NotificationType string                 `json:"notification_type"`
    Title            string                 `json:"title"`
    Message          string                 `json:"message"`
    Channels         []string               `json:"channels"`
    Priority         string                 `json:"priority"`
    Data             map[string]interface{} `json:"data"`
}

// Content Events
type ContentCreatedEvent struct {
    BaseEvent
    Payload ContentCreatedPayload `json:"payload"`
}

type ContentCreatedPayload struct {
    ContentID   string   `json:"content_id"`
    Title       string   `json:"title"`
    ContentType string   `json:"content_type"`
    CreatedBy   string   `json:"created_by"`
    CourseIDs   []string `json:"course_ids"`
}
```

### Consumer Group Configuration

```go
// shared/kafka/consumer.go

package kafka

import (
    "github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
    Brokers       []string
    GroupID       string
    Topic         string
    MinBytes      int
    MaxBytes      int
    MaxWait       time.Duration
    StartOffset   int64
    CommitInterval time.Duration
}

func NewConsumerGroup(cfg ConsumerConfig) *kafka.Reader {
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers:        cfg.Brokers,
        GroupID:        cfg.GroupID,
        Topic:          cfg.Topic,
        MinBytes:       cfg.MinBytes,      // 10KB
        MaxBytes:       cfg.MaxBytes,      // 10MB
        MaxWait:        cfg.MaxWait,       // 3s
        StartOffset:    cfg.StartOffset,   // kafka.FirstOffset
        CommitInterval: cfg.CommitInterval, // 1s
    })
}

// Consumer group IDs per service
const (
    UserServiceConsumerGroup         = "user-service-group"
    ContentServiceConsumerGroup      = "content-service-group"
    NotificationServiceConsumerGroup = "notification-service-group"
    CourseServiceConsumerGroup       = "course-service-group"
    TimetableServiceConsumerGroup    = "timetable-service-group"
    AttendanceServiceConsumerGroup   = "attendance-service-group"
    AnalyticsServiceConsumerGroup    = "analytics-service-group"
)
```

---

## PostgreSQL Configuration

### Database per Service

Each service has its own database:

```sql
-- Create databases
CREATE DATABASE user_service_db;
CREATE DATABASE content_service_db;
CREATE DATABASE notification_service_db;
CREATE DATABASE course_service_db;
CREATE DATABASE timetable_service_db;
CREATE DATABASE attendance_service_db;
CREATE DATABASE announcement_service_db;
CREATE DATABASE communication_service_db;
CREATE DATABASE analytics_service_db;
```

### Connection Pool Configuration

```go
// shared/database/postgres.go

package database

import (
    "context"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
    Host            string
    Port            int
    Database        string
    Username        string
    Password        string
    MaxConns        int32
    MinConns        int32
    MaxConnLifetime time.Duration
    MaxConnIdleTime time.Duration
}

func NewPostgresPool(cfg PostgresConfig) (*pgxpool.Pool, error) {
    connString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=require",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
    )
    
    poolConfig, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }
    
    poolConfig.MaxConns = cfg.MaxConns           // 25
    poolConfig.MinConns = cfg.MinConns           // 5
    poolConfig.MaxConnLifetime = cfg.MaxConnLifetime // 1 hour
    poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime // 30 minutes
    
    return pgxpool.NewWithConfig(context.Background(), poolConfig)
}
```

### Migration Strategy

Using golang-migrate:

```go
// migrations/runner.go

package migrations

import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string, migrationsPath string) error {
    m, err := migrate.New(
        "file://"+migrationsPath,
        dbURL,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

---

## Docker Compose Configuration

```yaml
# docker-compose.yml

version: '3.8'

services:
  # ===================
  # Infrastructure
  # ===================
  
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: nimbusu
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_MULTIPLE_DATABASES: user_service_db,content_service_db,notification_service_db,course_service_db,timetable_service_db,attendance_service_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-multiple-dbs.sh:/docker-entrypoint-initdb.d/init-multiple-dbs.sh
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U nimbusu"]
      interval: 10s
      timeout: 5s
      retries: 5

  mongodb:
    image: mongo:7.0
    environment:
      MONGO_INITDB_ROOT_USERNAME: nimbusu
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
    volumes:
      - mongodb_data:/data/db
    ports:
      - "27017:27017"
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - zookeeper_data:/var/lib/zookeeper/data

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
    volumes:
      - kafka_data:/var/lib/kafka/data
    healthcheck:
      test: kafka-broker-api-versions --bootstrap-server localhost:9092
      interval: 10s
      timeout: 5s
      retries: 5

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports:
      - "8090:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: nimbusu
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
    depends_on:
      - kafka

  # ===================
  # Services
  # ===================

  api-gateway:
    build:
      context: ./api-gateway
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - redis
      - user-service
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  user-service:
    build:
      context: ./services/user-service
      dockerfile: Dockerfile
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/user_service_db
      - REDIS_URL=redis://redis:6379
      - KAFKA_BROKERS=kafka:29092
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
      - kafka

  content-service:
    build:
      context: ./services/content-service
      dockerfile: Dockerfile
    ports:
      - "8082:8082"
    environment:
      - PORT=8082
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/content_service_db
      - MONGODB_URL=mongodb://nimbusu:${MONGO_PASSWORD}@mongodb:27017
      - KAFKA_BROKERS=kafka:29092
    depends_on:
      - postgres
      - mongodb
      - kafka

  notification-service:
    build:
      context: ./services/notification-service
      dockerfile: Dockerfile
    ports:
      - "8083:8083"
    environment:
      - PORT=8083
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/notification_service_db
      - KAFKA_BROKERS=kafka:29092
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASSWORD=${SMTP_PASSWORD}
      - TWILIO_SID=${TWILIO_SID}
      - TWILIO_TOKEN=${TWILIO_TOKEN}
    depends_on:
      - postgres
      - kafka

  course-service:
    build:
      context: ./services/course-service
      dockerfile: Dockerfile
    ports:
      - "8084:8084"
    environment:
      - PORT=8084
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/course_service_db
      - KAFKA_BROKERS=kafka:29092
    depends_on:
      - postgres
      - kafka

  timetable-service:
    build:
      context: ./services/timetable-service
      dockerfile: Dockerfile
    ports:
      - "8085:8085"
    environment:
      - PORT=8085
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/timetable_service_db
      - KAFKA_BROKERS=kafka:29092
    depends_on:
      - postgres
      - kafka

  attendance-service:
    build:
      context: ./services/attendance-service
      dockerfile: Dockerfile
    ports:
      - "8086:8086"
    environment:
      - PORT=8086
      - DATABASE_URL=postgres://nimbusu:${POSTGRES_PASSWORD}@postgres:5432/attendance_service_db
      - KAFKA_BROKERS=kafka:29092
    depends_on:
      - postgres
      - kafka

volumes:
  postgres_data:
  mongodb_data:
  redis_data:
  zookeeper_data:
  kafka_data:

networks:
  default:
    name: nimbusu-network
```

---

## Service Communication Patterns

### 1. Synchronous (REST/gRPC)

Used for:
- User authentication/authorization validation
- Real-time data queries
- CRUD operations

```go
// Example: Content service validating user permissions

func (s *ContentService) GetContent(ctx context.Context, contentID, userID string) (*Content, error) {
    // Call User Service to validate permissions
    permissions, err := s.userClient.GetUserPermissions(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if !hasPermission(permissions, "content:read") {
        return nil, ErrForbidden
    }
    
    return s.repo.FindByID(ctx, contentID)
}
```

### 2. Asynchronous (Kafka Events)

Used for:
- Event notifications
- Data synchronization
- Audit logging
- Analytics

```go
// Example: Publishing user created event

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Publish event asynchronously
    event := UserCreatedEvent{
        BaseEvent: BaseEvent{
            EventID:     uuid.New().String(),
            EventType:   "USER_CREATED",
            Timestamp:   time.Now(),
            ServiceName: "user-service",
            Version:     "1.0",
        },
        Payload: UserCreatedPayload{
            UserID:    user.ID,
            Email:     user.Email,
            Role:      user.Role,
            CreatedBy: req.CreatedBy,
        },
    }
    
    go s.producer.Publish("user.events", event)
    
    return user, nil
}
```

### 3. Request-Reply Pattern (via Kafka)

For complex cross-service operations:

```go
// Example: Saga pattern for enrollment

type EnrollmentSaga struct {
    Steps []SagaStep
}

func (s *EnrollmentSaga) Execute(ctx context.Context) error {
    for _, step := range s.Steps {
        if err := step.Execute(ctx); err != nil {
            return s.Compensate(ctx)
        }
    }
    return nil
}
```

---

## Error Handling

### Standard Error Response

```go
// shared/pkg/errors/errors.go

package errors

type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    Details    any    `json:"details,omitempty"`
    HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Common errors
var (
    ErrNotFound = &AppError{
        Code:       "NOT_FOUND",
        Message:    "Resource not found",
        HTTPStatus: 404,
    }
    
    ErrUnauthorized = &AppError{
        Code:       "UNAUTHORIZED",
        Message:    "Authentication required",
        HTTPStatus: 401,
    }
    
    ErrForbidden = &AppError{
        Code:       "FORBIDDEN",
        Message:    "Access denied",
        HTTPStatus: 403,
    }
    
    ErrValidation = &AppError{
        Code:       "VALIDATION_ERROR",
        Message:    "Validation failed",
        HTTPStatus: 400,
    }
)
```

---

## Observability

### Logging (Zap)

```go
// shared/logger/logger.go

package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(env string) (*zap.Logger, error) {
    var config zap.Config
    
    if env == "production" {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        config = zap.NewDevelopmentConfig()
    }
    
    return config.Build()
}
```

### Metrics (Prometheus)

```go
// Each service exposes /metrics endpoint

import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)
```

### Tracing (OpenTelemetry)

```go
// shared/tracing/tracing.go

package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(serviceName, jaegerURL string) (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint(jaegerURL),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return tp, nil
}
```

---

## Health Checks

Each service exposes health endpoints:

```go
// internal/handler/http/health.go

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    health := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Checks: map[string]HealthCheck{
            "database": h.checkDatabase(),
            "kafka":    h.checkKafka(),
            "redis":    h.checkRedis(),
        },
    }
    
    if !health.IsHealthy() {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    json.NewEncoder(w).Encode(health)
}
```

---

## Deployment Architecture

### Kubernetes (Production)

```yaml
# kubernetes/user-service/deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: nimbusu
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: nimbusu/user-service:latest
        ports:
        - containerPort: 8081
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: user-service-secrets
              key: database-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## Security Considerations

1. **Network Security**
   - All inter-service communication over internal network
   - TLS for external traffic
   - mTLS for service-to-service (optional)

2. **Authentication**
   - JWT tokens with short expiry (1 hour)
   - Refresh tokens with longer expiry (7 days)
   - Token blacklisting via Redis

3. **Authorization**
   - RBAC at API Gateway level
   - Fine-grained permissions per service

4. **Data Security**
   - Encryption at rest (PostgreSQL, MongoDB)
   - Encryption in transit (TLS)
   - PII data masking in logs

5. **Secrets Management**
   - Environment variables in development
   - Kubernetes Secrets / HashiCorp Vault in production

---

## Performance Considerations

1. **Caching Strategy**
   - Redis for session data
   - Redis for frequently accessed data (user permissions, roles)
   - Cache invalidation via Kafka events

2. **Database Optimization**
   - Connection pooling
   - Prepared statements
   - Proper indexing
   - Read replicas for analytics

3. **Kafka Optimization**
   - Appropriate partition count
   - Batch message production
   - Consumer group parallelism

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-12-27 | Initial architecture document |
