# NimbusU - Kafka Events Documentation

## Overview

NimbusU uses Apache Kafka as the backbone for asynchronous event-driven communication between microservices. This document provides comprehensive documentation for all Kafka topics, event schemas, and integration patterns.

## Architecture

```
                                    ┌─────────────────────────────────────────┐
                                    │            Apache Kafka                 │
                                    │         (Event Streaming)               │
                                    └─────────────────────────────────────────┘
                                                      │
        ┌─────────────────┬─────────────────┬─────────┴───────┬─────────────────┐
        │                 │                 │                 │                 │
        ▼                 ▼                 ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│ User Service  │ │Content Service│ │ Course Service│ │Timetable Svc  │ │Attendance Svc │
│               │ │               │ │               │ │               │ │               │
│ Produces:     │ │ Produces:     │ │ Produces:     │ │ Produces:     │ │ Produces:     │
│ - user.events │ │ - content.    │ │ - course.     │ │ - timetable.  │ │ - attendance. │
│ - auth.events │ │   events      │ │   events      │ │   events      │ │   events      │
└───────────────┘ └───────────────┘ └───────────────┘ └───────────────┘ └───────────────┘
        │                 │                 │                 │                 │
        └─────────────────┴─────────────────┴─────────┬───────┴─────────────────┘
                                                      │
                                                      ▼
                                    ┌─────────────────────────────────────────┐
                                    │       Notification Service              │
                                    │   (Consumes all event topics)           │
                                    │                                         │
                                    │   Subscribes to:                        │
                                    │   - notification.commands               │
                                    │   - user.events                         │
                                    │   - content.events                      │
                                    │   - attendance.events                   │
                                    │   - announcement.events                 │
                                    └─────────────────────────────────────────┘
```

---

## Kafka Cluster Configuration

### Broker Configuration

| Setting | Value | Description |
|---------|-------|-------------|
| Bootstrap Servers | `kafka:29092` (internal), `localhost:9092` (external) | Kafka broker addresses |
| Replication Factor | 3 (production), 1 (development) | Number of replicas per partition |
| Min In-Sync Replicas | 2 (production), 1 (development) | Minimum replicas for acknowledgment |
| Default Partitions | 6 | Default partition count for new topics |
| Message Retention | 7 days | Default message retention period |
| Max Message Size | 10 MB | Maximum message size |

### Zookeeper Configuration

| Setting | Value |
|---------|-------|
| Connection | `zookeeper:2181` |
| Session Timeout | 6000ms |
| Connection Timeout | 6000ms |

---

## Topic Definitions

### Topic Naming Convention

```
<domain>.<event-type>
```

Examples:
- `user.events` - User lifecycle events
- `notification.commands` - Notification requests (command pattern)
- `content.events` - Content lifecycle events

### Topic Registry

| Topic Name | Partitions | Retention | Replication | Description |
|------------|------------|-----------|-------------|-------------|
| `user.events` | 6 | 7 days | 3 | User lifecycle events (created, updated, deleted) |
| `auth.events` | 6 | 7 days | 3 | Authentication events (login, logout, password reset) |
| `content.events` | 12 | 7 days | 3 | Content lifecycle events |
| `notification.commands` | 12 | 1 day | 3 | Notification requests |
| `notification.status` | 6 | 3 days | 3 | Notification delivery status updates |
| `course.events` | 6 | 7 days | 3 | Course management events |
| `enrollment.events` | 6 | 7 days | 3 | Student enrollment events |
| `timetable.events` | 3 | 7 days | 3 | Timetable updates |
| `attendance.events` | 6 | 7 days | 3 | Attendance events |
| `announcement.events` | 3 | 7 days | 3 | Announcement events |
| `communication.events` | 6 | 3 days | 3 | Messaging events |
| `analytics.events` | 12 | 1 day | 3 | Analytics data collection |
| `dlq.events` | 3 | 30 days | 3 | Dead letter queue for failed events |

### Topic Creation Commands

```bash
# Create all topics
kafka-topics --create --bootstrap-server localhost:9092 --topic user.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic auth.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic content.events --partitions 12 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic notification.commands --partitions 12 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic notification.status --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic course.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic enrollment.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic timetable.events --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic attendance.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic announcement.events --partitions 3 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic communication.events --partitions 6 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic analytics.events --partitions 12 --replication-factor 3
kafka-topics --create --bootstrap-server localhost:9092 --topic dlq.events --partitions 3 --replication-factor 3
```

---

## Event Schema Definitions

### Base Event Structure

All events follow a consistent base structure:

```go
// BaseEvent is the common structure for all events
type BaseEvent struct {
    EventID       string            `json:"event_id"`        // UUID v4
    EventType     string            `json:"event_type"`      // Event type identifier
    EventVersion  string            `json:"event_version"`   // Schema version (e.g., "1.0")
    Timestamp     time.Time         `json:"timestamp"`       // Event creation time (ISO 8601)
    ServiceName   string            `json:"service_name"`    // Originating service
    CorrelationID string            `json:"correlation_id"`  // Request tracing ID
    Metadata      map[string]string `json:"metadata"`        // Additional metadata
}
```

### JSON Schema (Base Event)

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "service_name"],
  "properties": {
    "event_id": {
      "type": "string",
      "format": "uuid"
    },
    "event_type": {
      "type": "string",
      "pattern": "^[A-Z_]+$"
    },
    "event_version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+$"
    },
    "timestamp": {
      "type": "string",
      "format": "date-time"
    },
    "service_name": {
      "type": "string"
    },
    "correlation_id": {
      "type": "string"
    },
    "metadata": {
      "type": "object",
      "additionalProperties": {
        "type": "string"
      }
    }
  }
}
```

---

## 1. User Events (`user.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `USER_CREATED` | New user account created | User registration, admin creation |
| `USER_UPDATED` | User information updated | Profile update, role change |
| `USER_DELETED` | User account deleted/deactivated | Admin action |
| `USER_ACTIVATED` | User account activated | Email verification, admin action |
| `USER_SUSPENDED` | User account suspended | Admin action, policy violation |
| `PROFILE_UPDATED` | User profile updated | Profile edit |
| `ROLE_ASSIGNED` | Role assigned to user | Admin action |
| `ROLE_REVOKED` | Role revoked from user | Admin action |

### Event Schemas

#### USER_CREATED

```go
type UserCreatedEvent struct {
    BaseEvent
    Payload UserCreatedPayload `json:"payload"`
}

type UserCreatedPayload struct {
    UserID           string `json:"user_id"`
    RegisterNo       int64  `json:"register_no"`
    Email            string `json:"email"`
    RoleID           string `json:"role_id"`
    RoleName         string `json:"role_name"`
    Status           string `json:"status"`
    CreatedBy        string `json:"created_by"`
    RequiresApproval bool   `json:"requires_approval"`
}
```

**Example JSON:**

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "USER_CREATED",
  "event_version": "1.0",
  "timestamp": "2024-12-27T10:30:00Z",
  "service_name": "user-service",
  "correlation_id": "req-12345-abcde",
  "metadata": {
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0..."
  },
  "payload": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "register_no": 2024001234,
    "email": "student@university.edu",
    "role_id": "role-student-uuid",
    "role_name": "student",
    "status": "active",
    "created_by": "admin-user-uuid",
    "requires_approval": false
  }
}
```

#### USER_UPDATED

```go
type UserUpdatedEvent struct {
    BaseEvent
    Payload UserUpdatedPayload `json:"payload"`
}

type UserUpdatedPayload struct {
    UserID      string                 `json:"user_id"`
    Changes     map[string]ChangeValue `json:"changes"`
    UpdatedBy   string                 `json:"updated_by"`
    UpdatedAt   time.Time              `json:"updated_at"`
}

type ChangeValue struct {
    OldValue interface{} `json:"old_value"`
    NewValue interface{} `json:"new_value"`
}
```

**Example JSON:**

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440001",
  "event_type": "USER_UPDATED",
  "event_version": "1.0",
  "timestamp": "2024-12-27T11:00:00Z",
  "service_name": "user-service",
  "correlation_id": "req-12345-abcdf",
  "payload": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "changes": {
      "status": {
        "old_value": "active",
        "new_value": "suspended"
      },
      "role_id": {
        "old_value": "role-student-uuid",
        "new_value": "role-faculty-uuid"
      }
    },
    "updated_by": "admin-user-uuid",
    "updated_at": "2024-12-27T11:00:00Z"
  }
}
```

#### USER_DELETED

```go
type UserDeletedEvent struct {
    BaseEvent
    Payload UserDeletedPayload `json:"payload"`
}

type UserDeletedPayload struct {
    UserID      string    `json:"user_id"`
    Email       string    `json:"email"`
    DeletedBy   string    `json:"deleted_by"`
    DeletedAt   time.Time `json:"deleted_at"`
    Reason      string    `json:"reason"`
    IsSoftDelete bool     `json:"is_soft_delete"`
}
```

---

## 2. Authentication Events (`auth.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `LOGIN_SUCCESS` | Successful login | User login |
| `LOGIN_FAILED` | Failed login attempt | Invalid credentials |
| `LOGOUT` | User logout | Manual logout, session expiry |
| `PASSWORD_RESET_REQUESTED` | Password reset initiated | Forgot password request |
| `PASSWORD_RESET_COMPLETED` | Password successfully reset | Password reset completion |
| `PASSWORD_CHANGED` | Password changed | User changes password |
| `SESSION_CREATED` | New session created | Login |
| `SESSION_EXPIRED` | Session expired | Timeout |
| `TOKEN_REFRESHED` | Access token refreshed | Token refresh |
| `ACCOUNT_LOCKED` | Account locked due to failed attempts | Multiple failed logins |

### Event Schemas

#### LOGIN_SUCCESS

```go
type LoginSuccessEvent struct {
    BaseEvent
    Payload LoginSuccessPayload `json:"payload"`
}

type LoginSuccessPayload struct {
    UserID      string    `json:"user_id"`
    Email       string    `json:"email"`
    RoleName    string    `json:"role_name"`
    SessionID   string    `json:"session_id"`
    IPAddress   string    `json:"ip_address"`
    UserAgent   string    `json:"user_agent"`
    DeviceType  string    `json:"device_type"`
    LoginAt     time.Time `json:"login_at"`
    ExpiresAt   time.Time `json:"expires_at"`
}
```

**Example JSON:**

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440002",
  "event_type": "LOGIN_SUCCESS",
  "event_version": "1.0",
  "timestamp": "2024-12-27T09:00:00Z",
  "service_name": "user-service",
  "correlation_id": "req-login-12345",
  "payload": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "student@university.edu",
    "role_name": "student",
    "session_id": "sess-abc123",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
    "device_type": "desktop",
    "login_at": "2024-12-27T09:00:00Z",
    "expires_at": "2024-12-27T10:00:00Z"
  }
}
```

#### LOGIN_FAILED

```go
type LoginFailedEvent struct {
    BaseEvent
    Payload LoginFailedPayload `json:"payload"`
}

type LoginFailedPayload struct {
    Email           string    `json:"email"`
    FailureReason   string    `json:"failure_reason"`
    AttemptCount    int       `json:"attempt_count"`
    IPAddress       string    `json:"ip_address"`
    UserAgent       string    `json:"user_agent"`
    IsAccountLocked bool      `json:"is_account_locked"`
    AttemptedAt     time.Time `json:"attempted_at"`
}
```

---

## 3. Content Events (`content.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `CONTENT_CREATED` | New content uploaded | File upload |
| `CONTENT_UPDATED` | Content metadata updated | Edit content |
| `CONTENT_DELETED` | Content deleted | Delete action |
| `CONTENT_PUBLISHED` | Content made public | Publish action |
| `CONTENT_UNPUBLISHED` | Content made private | Unpublish action |
| `CONTENT_VERSIONED` | New version created | Version upload |
| `CONTENT_ACCESSED` | Content viewed/downloaded | User access |
| `CONTENT_SHARED` | Content shared with users | Share action |
| `CONTENT_PERMISSION_GRANTED` | Permission granted | Admin/faculty action |
| `CONTENT_PERMISSION_REVOKED` | Permission revoked | Admin/faculty action |

### Event Schemas

#### CONTENT_CREATED

```go
type ContentCreatedEvent struct {
    BaseEvent
    Payload ContentCreatedPayload `json:"payload"`
}

type ContentCreatedPayload struct {
    ContentID       string   `json:"content_id"`
    MongoDocumentID string   `json:"mongo_document_id"`
    Title           string   `json:"title"`
    Description     string   `json:"description"`
    ContentType     string   `json:"content_type"`
    MimeType        string   `json:"mime_type"`
    FileSize        int64    `json:"file_size"`
    CreatedBy       string   `json:"created_by"`
    CourseIDs       []string `json:"course_ids"`
    SubjectIDs      []string `json:"subject_ids"`
    DepartmentID    string   `json:"department_id"`
    IsPublished     bool     `json:"is_published"`
    Tags            []string `json:"tags"`
}
```

**Example JSON:**

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440003",
  "event_type": "CONTENT_CREATED",
  "event_version": "1.0",
  "timestamp": "2024-12-27T14:30:00Z",
  "service_name": "content-service",
  "correlation_id": "req-upload-67890",
  "payload": {
    "content_id": "content-uuid-12345",
    "mongo_document_id": "507f1f77bcf86cd799439011",
    "title": "Introduction to Data Structures - Week 1",
    "description": "Lecture notes covering arrays and linked lists",
    "content_type": "document",
    "mime_type": "application/pdf",
    "file_size": 2048576,
    "created_by": "faculty-uuid-12345",
    "course_ids": ["course-uuid-001"],
    "subject_ids": ["subject-cs101"],
    "department_id": "dept-cse-uuid",
    "is_published": true,
    "tags": ["data-structures", "arrays", "linked-lists", "week-1"]
  }
}
```

#### CONTENT_ACCESSED

```go
type ContentAccessedEvent struct {
    BaseEvent
    Payload ContentAccessedPayload `json:"payload"`
}

type ContentAccessedPayload struct {
    ContentID   string    `json:"content_id"`
    UserID      string    `json:"user_id"`
    AccessType  string    `json:"access_type"` // view, download, share
    IPAddress   string    `json:"ip_address"`
    UserAgent   string    `json:"user_agent"`
    DeviceType  string    `json:"device_type"`
    Duration    int       `json:"duration_seconds"` // For view events
    AccessedAt  time.Time `json:"accessed_at"`
}
```

---

## 4. Notification Commands (`notification.commands`)

### Command Types

| Command Type | Description | Trigger |
|--------------|-------------|---------|
| `SEND_EMAIL` | Send email notification | Various events |
| `SEND_SMS` | Send SMS notification | Critical alerts |
| `SEND_PUSH` | Send push notification | Mobile app |
| `SEND_IN_APP` | Send in-app notification | All notifications |
| `SEND_BULK` | Send to multiple recipients | Announcements |
| `SCHEDULE_NOTIFICATION` | Schedule future notification | Scheduled events |
| `CANCEL_NOTIFICATION` | Cancel scheduled notification | Cancellation |

### Command Schemas

#### SEND_NOTIFICATION (Generic)

```go
type SendNotificationCommand struct {
    BaseEvent
    Payload NotificationPayload `json:"payload"`
}

type NotificationPayload struct {
    NotificationID   string                 `json:"notification_id"`
    RecipientUserID  string                 `json:"recipient_user_id"`
    RecipientEmail   string                 `json:"recipient_email"`
    RecipientPhone   string                 `json:"recipient_phone"`
    NotificationType string                 `json:"notification_type"`
    Title            string                 `json:"title"`
    Message          string                 `json:"message"`
    Channels         []string               `json:"channels"` // email, sms, push, in_app
    Priority         string                 `json:"priority"` // low, normal, high, urgent
    TemplateID       string                 `json:"template_id"`
    TemplateData     map[string]interface{} `json:"template_data"`
    ActionURL        string                 `json:"action_url"`
    ScheduledAt      *time.Time             `json:"scheduled_at"`
    ExpiresAt        *time.Time             `json:"expires_at"`
}
```

**Example JSON:**

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440004",
  "event_type": "SEND_NOTIFICATION",
  "event_version": "1.0",
  "timestamp": "2024-12-27T15:00:00Z",
  "service_name": "content-service",
  "correlation_id": "req-content-share-001",
  "payload": {
    "notification_id": "notif-uuid-12345",
    "recipient_user_id": "student-uuid-001",
    "recipient_email": "student@university.edu",
    "notification_type": "content",
    "title": "New Course Material Available",
    "message": "New lecture notes for CS101 have been uploaded.",
    "channels": ["email", "push", "in_app"],
    "priority": "normal",
    "template_id": "template-new-content",
    "template_data": {
      "course_name": "Data Structures",
      "content_title": "Week 1 - Introduction",
      "faculty_name": "Dr. Smith"
    },
    "action_url": "/content/view/content-uuid-12345",
    "scheduled_at": null,
    "expires_at": null
  }
}
```

#### SEND_BULK_NOTIFICATION

```go
type SendBulkNotificationCommand struct {
    BaseEvent
    Payload BulkNotificationPayload `json:"payload"`
}

type BulkNotificationPayload struct {
    BatchID          string                 `json:"batch_id"`
    RecipientType    string                 `json:"recipient_type"` // all, department, course, group
    RecipientFilter  RecipientFilter        `json:"recipient_filter"`
    NotificationType string                 `json:"notification_type"`
    Title            string                 `json:"title"`
    Message          string                 `json:"message"`
    Channels         []string               `json:"channels"`
    Priority         string                 `json:"priority"`
    TemplateID       string                 `json:"template_id"`
    TemplateData     map[string]interface{} `json:"template_data"`
}

type RecipientFilter struct {
    DepartmentIDs []string `json:"department_ids"`
    CourseIDs     []string `json:"course_ids"`
    ProgramIDs    []string `json:"program_ids"`
    GroupIDs      []string `json:"group_ids"`
    RoleNames     []string `json:"role_names"`
    Semesters     []int    `json:"semesters"`
}
```

---

## 5. Notification Status (`notification.status`)

### Event Types

| Event Type | Description |
|------------|-------------|
| `NOTIFICATION_QUEUED` | Notification added to delivery queue |
| `NOTIFICATION_SENT` | Notification sent successfully |
| `NOTIFICATION_DELIVERED` | Notification confirmed delivered |
| `NOTIFICATION_FAILED` | Notification delivery failed |
| `NOTIFICATION_READ` | Notification marked as read |

### Event Schema

```go
type NotificationStatusEvent struct {
    BaseEvent
    Payload NotificationStatusPayload `json:"payload"`
}

type NotificationStatusPayload struct {
    NotificationID string    `json:"notification_id"`
    QueueID        string    `json:"queue_id"`
    UserID         string    `json:"user_id"`
    Channel        string    `json:"channel"`
    Status         string    `json:"status"`
    AttemptCount   int       `json:"attempt_count"`
    ErrorMessage   string    `json:"error_message"`
    SentAt         time.Time `json:"sent_at"`
    DeliveredAt    time.Time `json:"delivered_at"`
    ReadAt         time.Time `json:"read_at"`
}
```

---

## 6. Course Events (`course.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `COURSE_CREATED` | New course created | Admin action |
| `COURSE_UPDATED` | Course information updated | Edit course |
| `COURSE_DELETED` | Course deleted | Admin action |
| `COURSE_ACTIVATED` | Course made active | Status change |
| `COURSE_DEACTIVATED` | Course made inactive | Status change |
| `FACULTY_ASSIGNED` | Faculty assigned to course | Admin action |
| `FACULTY_UNASSIGNED` | Faculty removed from course | Admin action |

### Event Schema

#### COURSE_CREATED

```go
type CourseCreatedEvent struct {
    BaseEvent
    Payload CourseCreatedPayload `json:"payload"`
}

type CourseCreatedPayload struct {
    CourseID     string `json:"course_id"`
    CourseCode   string `json:"course_code"`
    CourseName   string `json:"course_name"`
    SubjectID    string `json:"subject_id"`
    DepartmentID string `json:"department_id"`
    ProgramID    string `json:"program_id"`
    Semester     int    `json:"semester"`
    AcademicYear int    `json:"academic_year"`
    MaxStudents  int    `json:"max_students"`
    CreatedBy    string `json:"created_by"`
}
```

---

## 7. Enrollment Events (`enrollment.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `STUDENT_ENROLLED` | Student enrolled in course | Enrollment |
| `STUDENT_DROPPED` | Student dropped course | Drop action |
| `ENROLLMENT_COMPLETED` | Course completion recorded | End of semester |
| `ENROLLMENT_FAILED` | Student failed course | Grading |
| `WAITLIST_ADDED` | Student added to waitlist | Full course |
| `WAITLIST_PROMOTED` | Student promoted from waitlist | Spot available |

### Event Schema

```go
type StudentEnrolledEvent struct {
    BaseEvent
    Payload EnrollmentPayload `json:"payload"`
}

type EnrollmentPayload struct {
    EnrollmentID   string    `json:"enrollment_id"`
    StudentID      string    `json:"student_id"`
    CourseID       string    `json:"course_id"`
    CourseName     string    `json:"course_name"`
    EnrollmentDate time.Time `json:"enrollment_date"`
    EnrolledBy     string    `json:"enrolled_by"` // self, admin
    CurrentCount   int       `json:"current_count"`
    MaxCapacity    int       `json:"max_capacity"`
}
```

---

## 8. Timetable Events (`timetable.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `TIMETABLE_CREATED` | New timetable created | Admin action |
| `TIMETABLE_PUBLISHED` | Timetable published | Admin action |
| `TIMETABLE_UPDATED` | Timetable modified | Edit action |
| `ENTRY_ADDED` | Schedule entry added | Admin action |
| `ENTRY_MODIFIED` | Schedule entry modified | Edit action |
| `ENTRY_DELETED` | Schedule entry deleted | Delete action |
| `ROOM_CONFLICT_DETECTED` | Room booking conflict | Validation |
| `FACULTY_CONFLICT_DETECTED` | Faculty schedule conflict | Validation |

### Event Schema

```go
type TimetablePublishedEvent struct {
    BaseEvent
    Payload TimetablePublishedPayload `json:"payload"`
}

type TimetablePublishedPayload struct {
    TimetableID    string    `json:"timetable_id"`
    SemesterID     string    `json:"semester_id"`
    SemesterName   string    `json:"semester_name"`
    DepartmentID   string    `json:"department_id"`
    DepartmentName string    `json:"department_name"`
    ProgramID      string    `json:"program_id"`
    SemesterNumber int       `json:"semester_number"`
    EntryCount     int       `json:"entry_count"`
    PublishedBy    string    `json:"published_by"`
    PublishedAt    time.Time `json:"published_at"`
    EffectiveFrom  time.Time `json:"effective_from"`
}
```

---

## 9. Attendance Events (`attendance.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `ATTENDANCE_MARKED` | Attendance recorded for class | Faculty action |
| `ATTENDANCE_UPDATED` | Attendance record modified | Correction |
| `ATTENDANCE_BULK_MARKED` | Bulk attendance recorded | Batch update |
| `LOW_ATTENDANCE_ALERT` | Student below attendance threshold | System check |
| `ATTENDANCE_REPORT_GENERATED` | Report generated | Request |

### Event Schema

#### ATTENDANCE_MARKED

```go
type AttendanceMarkedEvent struct {
    BaseEvent
    Payload AttendanceMarkedPayload `json:"payload"`
}

type AttendanceMarkedPayload struct {
    AttendanceID   string              `json:"attendance_id"`
    EntryID        string              `json:"entry_id"`
    CourseID       string              `json:"course_id"`
    CourseName     string              `json:"course_name"`
    ClassDate      string              `json:"class_date"` // YYYY-MM-DD
    MarkedBy       string              `json:"marked_by"`
    MarkedAt       time.Time           `json:"marked_at"`
    TotalStudents  int                 `json:"total_students"`
    PresentCount   int                 `json:"present_count"`
    AbsentCount    int                 `json:"absent_count"`
    StudentRecords []StudentAttendance `json:"student_records"`
}

type StudentAttendance struct {
    StudentID string `json:"student_id"`
    Status    string `json:"status"` // present, absent, late, excused
    Remarks   string `json:"remarks"`
}
```

#### LOW_ATTENDANCE_ALERT

```go
type LowAttendanceAlertEvent struct {
    BaseEvent
    Payload LowAttendanceAlertPayload `json:"payload"`
}

type LowAttendanceAlertPayload struct {
    StudentID            string  `json:"student_id"`
    StudentName          string  `json:"student_name"`
    CourseID             string  `json:"course_id"`
    CourseName           string  `json:"course_name"`
    AttendancePercentage float64 `json:"attendance_percentage"`
    ThresholdPercentage  float64 `json:"threshold_percentage"`
    TotalClasses         int     `json:"total_classes"`
    ClassesAttended      int     `json:"classes_attended"`
    AlertLevel           string  `json:"alert_level"` // warning, critical
}
```

---

## 10. Announcement Events (`announcement.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `ANNOUNCEMENT_CREATED` | New announcement created | Admin/faculty action |
| `ANNOUNCEMENT_PUBLISHED` | Announcement published | Publish action |
| `ANNOUNCEMENT_UPDATED` | Announcement modified | Edit action |
| `ANNOUNCEMENT_DELETED` | Announcement deleted | Delete action |
| `ANNOUNCEMENT_EXPIRED` | Announcement expired | Auto-expiry |
| `ANNOUNCEMENT_SCHEDULED` | Announcement scheduled | Schedule action |

### Event Schema

```go
type AnnouncementPublishedEvent struct {
    BaseEvent
    Payload AnnouncementPublishedPayload `json:"payload"`
}

type AnnouncementPublishedPayload struct {
    AnnouncementID   string    `json:"announcement_id"`
    Title            string    `json:"title"`
    Content          string    `json:"content"`
    AnnouncementType string    `json:"announcement_type"`
    Priority         string    `json:"priority"`
    TargetAudience   string    `json:"target_audience"`
    TargetDetails    Target    `json:"target_details"`
    CreatedBy        string    `json:"created_by"`
    PublishedAt      time.Time `json:"published_at"`
    ExpiresAt        time.Time `json:"expires_at"`
}

type Target struct {
    DepartmentIDs []string `json:"department_ids"`
    ProgramIDs    []string `json:"program_ids"`
    CourseIDs     []string `json:"course_ids"`
    GroupIDs      []string `json:"group_ids"`
    UserIDs       []string `json:"user_ids"`
}
```

---

## 11. Communication Events (`communication.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `MESSAGE_SENT` | New message sent | User action |
| `MESSAGE_READ` | Message marked as read | Read action |
| `MESSAGE_DELETED` | Message deleted | Delete action |
| `CONVERSATION_CREATED` | New conversation started | User action |
| `CONVERSATION_ARCHIVED` | Conversation archived | Archive action |
| `PARTICIPANT_ADDED` | User added to conversation | Invite action |
| `PARTICIPANT_REMOVED` | User removed from conversation | Remove action |
| `FORUM_POST_CREATED` | New forum post created | User action |
| `FORUM_REPLY_ADDED` | Reply added to forum post | User action |

### Event Schema

```go
type MessageSentEvent struct {
    BaseEvent
    Payload MessageSentPayload `json:"payload"`
}

type MessageSentPayload struct {
    MessageID        string    `json:"message_id"`
    ConversationID   string    `json:"conversation_id"`
    ConversationType string    `json:"conversation_type"`
    SenderID         string    `json:"sender_id"`
    SenderName       string    `json:"sender_name"`
    MessageText      string    `json:"message_text"`
    MessageType      string    `json:"message_type"`
    AttachmentURL    string    `json:"attachment_url"`
    ReplyToMessageID string    `json:"reply_to_message_id"`
    RecipientIDs     []string  `json:"recipient_ids"`
    SentAt           time.Time `json:"sent_at"`
}
```

---

## 12. Analytics Events (`analytics.events`)

### Event Types

| Event Type | Description | Trigger |
|------------|-------------|---------|
| `PAGE_VIEW` | Page/screen viewed | User navigation |
| `USER_ACTION` | User performed action | Any user action |
| `SESSION_START` | User session started | Login |
| `SESSION_END` | User session ended | Logout/timeout |
| `FEATURE_USED` | Feature usage tracked | Feature interaction |
| `ERROR_OCCURRED` | Error tracked | System error |
| `PERFORMANCE_METRIC` | Performance data | System measurement |

### Event Schema

```go
type AnalyticsEvent struct {
    BaseEvent
    Payload AnalyticsPayload `json:"payload"`
}

type AnalyticsPayload struct {
    UserID       string                 `json:"user_id"`
    SessionID    string                 `json:"session_id"`
    EventName    string                 `json:"event_name"`
    Category     string                 `json:"category"`
    Action       string                 `json:"action"`
    Label        string                 `json:"label"`
    Value        float64                `json:"value"`
    PageURL      string                 `json:"page_url"`
    Referrer     string                 `json:"referrer"`
    DeviceType   string                 `json:"device_type"`
    Browser      string                 `json:"browser"`
    OS           string                 `json:"os"`
    ScreenSize   string                 `json:"screen_size"`
    Properties   map[string]interface{} `json:"properties"`
}
```

---

## Consumer Groups

### Configuration

| Consumer Group | Service | Topics Subscribed |
|----------------|---------|-------------------|
| `user-service-group` | User Service | - |
| `content-service-group` | Content Service | `user.events` |
| `notification-service-group` | Notification Service | `notification.commands`, `user.events`, `content.events`, `attendance.events`, `announcement.events` |
| `course-service-group` | Course Service | `user.events` |
| `timetable-service-group` | Timetable Service | `course.events` |
| `attendance-service-group` | Attendance Service | `timetable.events`, `course.events`, `enrollment.events` |
| `announcement-service-group` | Announcement Service | `user.events` |
| `communication-service-group` | Communication Service | `user.events`, `course.events` |
| `analytics-service-group` | Analytics Service | All event topics |

### Consumer Group Code Example

```go
package kafka

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
    Brokers        []string
    GroupID        string
    Topics         []string
    MinBytes       int
    MaxBytes       int
    MaxWait        time.Duration
    CommitInterval time.Duration
}

type Consumer struct {
    readers []*kafka.Reader
    handler EventHandler
}

type EventHandler interface {
    HandleEvent(ctx context.Context, event BaseEvent, payload json.RawMessage) error
}

func NewConsumer(cfg ConsumerConfig, handler EventHandler) *Consumer {
    var readers []*kafka.Reader
    
    for _, topic := range cfg.Topics {
        reader := kafka.NewReader(kafka.ReaderConfig{
            Brokers:        cfg.Brokers,
            GroupID:        cfg.GroupID,
            Topic:          topic,
            MinBytes:       cfg.MinBytes,      // 10KB
            MaxBytes:       cfg.MaxBytes,      // 10MB
            MaxWait:        cfg.MaxWait,       // 3s
            CommitInterval: cfg.CommitInterval, // 1s
            StartOffset:    kafka.FirstOffset,
        })
        readers = append(readers, reader)
    }
    
    return &Consumer{
        readers: readers,
        handler: handler,
    }
}

func (c *Consumer) Start(ctx context.Context) {
    for _, reader := range c.readers {
        go c.consume(ctx, reader)
    }
}

func (c *Consumer) consume(ctx context.Context, reader *kafka.Reader) {
    for {
        select {
        case <-ctx.Done():
            reader.Close()
            return
        default:
            msg, err := reader.ReadMessage(ctx)
            if err != nil {
                log.Printf("Error reading message: %v", err)
                continue
            }
            
            var event struct {
                BaseEvent
                Payload json.RawMessage `json:"payload"`
            }
            
            if err := json.Unmarshal(msg.Value, &event); err != nil {
                log.Printf("Error unmarshaling event: %v", err)
                continue
            }
            
            if err := c.handler.HandleEvent(ctx, event.BaseEvent, event.Payload); err != nil {
                log.Printf("Error handling event: %v", err)
                // Send to DLQ if needed
            }
        }
    }
}
```

---

## Producer Configuration

### Producer Code Example

```go
package kafka

import (
    "context"
    "encoding/json"
    "time"

    "github.com/google/uuid"
    "github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
    Brokers      []string
    BatchSize    int
    BatchTimeout time.Duration
    RequiredAcks kafka.RequiredAcks
    Async        bool
}

type Producer struct {
    writer *kafka.Writer
    config ProducerConfig
}

func NewProducer(cfg ProducerConfig) *Producer {
    writer := &kafka.Writer{
        Addr:         kafka.TCP(cfg.Brokers...),
        BatchSize:    cfg.BatchSize,
        BatchTimeout: cfg.BatchTimeout,
        RequiredAcks: cfg.RequiredAcks,
        Async:        cfg.Async,
    }
    
    return &Producer{
        writer: writer,
        config: cfg,
    }
}

func (p *Producer) Publish(ctx context.Context, topic string, event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    return p.writer.WriteMessages(ctx, kafka.Message{
        Topic: topic,
        Key:   []byte(uuid.New().String()),
        Value: data,
        Time:  time.Now(),
    })
}

func (p *Producer) PublishWithKey(ctx context.Context, topic, key string, event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    return p.writer.WriteMessages(ctx, kafka.Message{
        Topic: topic,
        Key:   []byte(key),
        Value: data,
        Time:  time.Now(),
    })
}

func (p *Producer) Close() error {
    return p.writer.Close()
}
```

---

## Error Handling & Dead Letter Queue

### DLQ Event Structure

```go
type DeadLetterEvent struct {
    OriginalTopic     string          `json:"original_topic"`
    OriginalEvent     json.RawMessage `json:"original_event"`
    OriginalKey       string          `json:"original_key"`
    OriginalTimestamp time.Time       `json:"original_timestamp"`
    FailureReason     string          `json:"failure_reason"`
    FailureCount      int             `json:"failure_count"`
    LastAttemptAt     time.Time       `json:"last_attempt_at"`
    ConsumerGroup     string          `json:"consumer_group"`
    ServiceName       string          `json:"service_name"`
    StackTrace        string          `json:"stack_trace"`
}
```

### Retry Strategy

```go
type RetryConfig struct {
    MaxRetries     int
    InitialBackoff time.Duration
    MaxBackoff     time.Duration
    BackoffFactor  float64
}

var DefaultRetryConfig = RetryConfig{
    MaxRetries:     3,
    InitialBackoff: 100 * time.Millisecond,
    MaxBackoff:     5 * time.Second,
    BackoffFactor:  2.0,
}

func (c *Consumer) processWithRetry(ctx context.Context, msg kafka.Message, cfg RetryConfig) error {
    var lastErr error
    backoff := cfg.InitialBackoff
    
    for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
        err := c.processMessage(ctx, msg)
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        if attempt < cfg.MaxRetries {
            time.Sleep(backoff)
            backoff = time.Duration(float64(backoff) * cfg.BackoffFactor)
            if backoff > cfg.MaxBackoff {
                backoff = cfg.MaxBackoff
            }
        }
    }
    
    // Send to DLQ after all retries exhausted
    return c.sendToDLQ(ctx, msg, lastErr)
}
```

---

## Monitoring & Observability

### Metrics to Track

| Metric | Type | Description |
|--------|------|-------------|
| `kafka_messages_produced_total` | Counter | Total messages produced per topic |
| `kafka_messages_consumed_total` | Counter | Total messages consumed per topic |
| `kafka_consumer_lag` | Gauge | Consumer lag per partition |
| `kafka_message_processing_duration` | Histogram | Message processing time |
| `kafka_producer_errors_total` | Counter | Producer errors |
| `kafka_consumer_errors_total` | Counter | Consumer errors |
| `kafka_dlq_messages_total` | Counter | Messages sent to DLQ |

### Prometheus Metrics Example

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    MessagesProduced = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "kafka_messages_produced_total",
            Help: "Total number of messages produced",
        },
        []string{"topic", "event_type"},
    )

    MessagesConsumed = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "kafka_messages_consumed_total",
            Help: "Total number of messages consumed",
        },
        []string{"topic", "consumer_group", "event_type"},
    )

    ConsumerLag = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "kafka_consumer_lag",
            Help: "Consumer lag per partition",
        },
        []string{"topic", "consumer_group", "partition"},
    )

    ProcessingDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "kafka_message_processing_duration_seconds",
            Help:    "Message processing duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"topic", "event_type"},
    )

    DLQMessages = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "kafka_dlq_messages_total",
            Help: "Total messages sent to DLQ",
        },
        []string{"original_topic", "failure_reason"},
    )
)
```

---

## Best Practices

### 1. Event Design

- **Immutability**: Events should be immutable once published
- **Self-contained**: Include all necessary data; avoid requiring lookups
- **Versioning**: Always include schema version for backward compatibility
- **Idempotency**: Design consumers to handle duplicate events

### 2. Partitioning Strategy

```go
// Use meaningful partition keys for ordering guarantees
// User events: partition by user_id
// Content events: partition by content_id
// Course events: partition by course_id

func getPartitionKey(event interface{}) string {
    switch e := event.(type) {
    case *UserCreatedEvent:
        return e.Payload.UserID
    case *ContentCreatedEvent:
        return e.Payload.ContentID
    case *AttendanceMarkedEvent:
        return e.Payload.CourseID
    default:
        return uuid.New().String()
    }
}
```

### 3. Consumer Best Practices

- Use consumer groups for horizontal scaling
- Implement graceful shutdown
- Handle rebalancing properly
- Commit offsets after successful processing
- Implement circuit breakers for downstream services

### 4. Producer Best Practices

- Use batching for high-throughput scenarios
- Implement proper error handling with retries
- Use appropriate acknowledgment levels
- Consider async vs sync based on requirements

### 5. Schema Evolution

```go
// Version 1.0
type UserCreatedPayloadV1 struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
}

// Version 1.1 - Added new field (backward compatible)
type UserCreatedPayloadV11 struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name,omitempty"` // New optional field
}
```

---

## Docker Compose (Kafka Setup)

```yaml
version: '3.8'

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    hostname: zookeeper
    container_name: nimbusu-zookeeper
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - zookeeper_data:/var/lib/zookeeper/data
      - zookeeper_logs:/var/lib/zookeeper/log
    networks:
      - nimbusu-network

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    hostname: kafka
    container_name: nimbusu-kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
      - "29092:29092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "false"
      KAFKA_LOG_RETENTION_HOURS: 168
      KAFKA_LOG_RETENTION_BYTES: 1073741824
      KAFKA_MESSAGE_MAX_BYTES: 10485760
    volumes:
      - kafka_data:/var/lib/kafka/data
    networks:
      - nimbusu-network
    healthcheck:
      test: kafka-broker-api-versions --bootstrap-server localhost:9092
      interval: 10s
      timeout: 10s
      retries: 5

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    container_name: nimbusu-kafka-ui
    ports:
      - "8090:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: nimbusu-local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
      KAFKA_CLUSTERS_0_ZOOKEEPER: zookeeper:2181
    depends_on:
      - kafka
    networks:
      - nimbusu-network

  kafka-init:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      kafka:
        condition: service_healthy
    entrypoint: ["/bin/sh", "-c"]
    command: |
      "
      echo 'Creating Kafka topics...'
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic user.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic auth.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic content.events --partitions 12 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic notification.commands --partitions 12 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic notification.status --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic course.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic enrollment.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic timetable.events --partitions 3 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic attendance.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic announcement.events --partitions 3 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic communication.events --partitions 6 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic analytics.events --partitions 12 --replication-factor 1
      kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic dlq.events --partitions 3 --replication-factor 1
      echo 'Topics created successfully!'
      kafka-topics --list --bootstrap-server kafka:29092
      "
    networks:
      - nimbusu-network

volumes:
  zookeeper_data:
  zookeeper_logs:
  kafka_data:

networks:
  nimbusu-network:
    driver: bridge
```

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-12-27 | Initial Kafka events documentation |

---

## Quick Reference

### Event Type Lookup

| Service | Produces To | Subscribes To |
|---------|-------------|---------------|
| User Service | `user.events`, `auth.events` | - |
| Content Service | `content.events` | `user.events` |
| Course Service | `course.events`, `enrollment.events` | `user.events` |
| Timetable Service | `timetable.events` | `course.events` |
| Attendance Service | `attendance.events` | `timetable.events`, `course.events`, `enrollment.events` |
| Announcement Service | `announcement.events` | `user.events` |
| Communication Service | `communication.events` | `user.events`, `course.events` |
| Notification Service | `notification.status` | `notification.commands`, `user.events`, `content.events`, `attendance.events`, `announcement.events` |
| Analytics Service | `analytics.events` | All topics |
