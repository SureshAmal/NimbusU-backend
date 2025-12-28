# NimbusU - University Content Management System
# Complete Database Schema Documentation

## Overview

This document provides comprehensive documentation for the NimbusU database schema, designed for a microservices architecture using:
- **PostgreSQL**: Core relational data storage
- **MongoDB**: Document/file content storage (GridFS)
- **Kafka**: Event-driven communication between services

## Database Architecture

```
+------------------+     +------------------+     +------------------+
|   PostgreSQL     |     |     MongoDB      |     |      Kafka       |
|   (Relational)   |     |   (Documents)    |     |    (Events)      |
+------------------+     +------------------+     +------------------+
        |                        |                        |
        v                        v                        v
+-----------------------------------------------------------------------+
|                        MICROSERVICES                                   |
+-----------------------------------------------------------------------+
| User Service | Content Service | Notification Service | Timetable Svc |
+-----------------------------------------------------------------------+
```

---

## 1. USER MANAGEMENT SERVICE

### Purpose
Handles authentication, authorization, user profiles, and activity logging.

### Tables

#### 1.1 `users`
Primary user identity table for authentication.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| user_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| register_no | BIGINT | UNIQUE, NOT NULL | University registration number |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email address (login identifier) |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| role_id | UUID | FK -> roles.role_id, NOT NULL | User's role |
| status | VARCHAR(20) | DEFAULT 'active' | active, inactive, suspended |
| last_login | TIMESTAMPTZ | NULL | Last successful login timestamp |
| created_at | TIMESTAMPTZ | DEFAULT now() | Account creation time |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |
| created_by | UUID | FK -> users.user_id | Admin who created this user |

**Indexes**: email, role_id, status, register_no

---

#### 1.2 `roles`
Defines user roles in the system.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| role_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| role_name | VARCHAR(50) | UNIQUE, NOT NULL | Role name (admin, faculty, student) |
| description | TEXT | NULL | Role description |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

**Default Roles**:
- `admin` - System administrators
- `faculty` - Teaching staff
- `student` - Enrolled students
- `staff` - Non-teaching staff

---

#### 1.3 `permissions`
Defines granular permissions for RBAC.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| permission_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| permission_name | VARCHAR(100) | UNIQUE, NOT NULL | Human-readable name |
| resource | VARCHAR(100) | NOT NULL | Resource type (user, content, timetable) |
| action | VARCHAR(50) | NOT NULL | Action (create, read, update, delete) |
| description | TEXT | NULL | Permission description |

**Common Permissions**:
- `user:create`, `user:read`, `user:update`, `user:delete`
- `content:create`, `content:read`, `content:update`, `content:delete`, `content:publish`
- `timetable:create`, `timetable:read`, `timetable:update`, `timetable:publish`

---

#### 1.4 `role_permissions`
Maps roles to permissions (many-to-many).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| role_permission_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| role_id | UUID | FK -> roles.role_id, NOT NULL | Role reference |
| permission_id | UUID | FK -> permissions.permission_id, NOT NULL | Permission reference |
| created_at | TIMESTAMPTZ | DEFAULT now() | Mapping creation time |

**Unique Index**: (role_id, permission_id)

---

#### 1.5 `user_profiles`
Extended user profile information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| profile_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id (1:1), NOT NULL | User reference |
| register_no | BIGINT | UNIQUE, NOT NULL | Registration number (denormalized) |
| first_name | VARCHAR(100) | NOT NULL | First name |
| middle_name | VARCHAR(255) | NULL | Middle name |
| last_name | VARCHAR(100) | NOT NULL | Last name |
| phone | VARCHAR(20) | NULL | Phone number |
| gender | VARCHAR(10) | NULL | Gender (male, female, other) |
| profile_picture_url | TEXT | NULL | Profile picture URL/MongoDB ref |
| bio | TEXT | NULL | User biography |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 1.6 `user_activity_logs`
Audit trail for user actions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| log_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id, NOT NULL | User who performed action |
| action | VARCHAR(100) | NOT NULL | Action type (login, logout, create_content) |
| resource_type | VARCHAR(50) | NULL | Affected resource type |
| resource_id | UUID | NULL | Affected resource ID |
| ip_address | VARCHAR(45) | NULL | IPv4/IPv6 address |
| user_agent | TEXT | NULL | Browser/client user agent |
| details | JSONB | NULL | Additional action details |
| created_at | TIMESTAMPTZ | DEFAULT now() | Action timestamp |

**Indexes**: user_id, action, created_at, resource_type

---

#### 1.7 `password_reset_tokens`
Temporary tokens for password reset flow.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| token_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id, NOT NULL | User requesting reset |
| token | VARCHAR(255) | UNIQUE, NOT NULL | Secure random token |
| expires_at | TIMESTAMPTZ | NOT NULL | Token expiration time |
| used | BOOLEAN | DEFAULT false | Whether token was used |
| created_at | TIMESTAMPTZ | DEFAULT now() | Token creation time |

---

## 2. ACADEMIC STRUCTURE SERVICE

### Purpose
Manages departments, programs, faculty profiles, and student profiles.

### Tables

#### 2.1 `departments`
Academic departments in the university.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| department_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| department_name | VARCHAR(100) | UNIQUE, NOT NULL | Full department name |
| department_code | VARCHAR(20) | UNIQUE, NOT NULL | Short code (e.g., CSE, ECE) |
| head_of_department | UUID | FK -> faculties.faculty_id | Department head |
| description | TEXT | NULL | Department description |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 2.2 `programs`
Degree programs offered.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| program_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| program_name | VARCHAR(100) | NOT NULL | Program name (B.Tech Computer Science) |
| program_code | VARCHAR(20) | UNIQUE, NOT NULL | Short code (BTCS, MTech) |
| department_id | UUID | FK -> departments.department_id | Parent department |
| degree_type | VARCHAR(50) | NULL | Bachelor, Master, PhD |
| duration_years | INTEGER | NOT NULL | Program duration in years |
| total_credits | INTEGER | NULL | Total credits required |
| description | TEXT | NULL | Program description |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 2.3 `faculties`
Faculty member profiles (extends users).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| faculty_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id (1:1), NOT NULL | User account reference |
| employee_id | VARCHAR(50) | UNIQUE, NOT NULL | Employee ID number |
| department_id | UUID | FK -> departments.department_id | Department assignment |
| designation | VARCHAR(100) | NULL | Job title (Professor, Assistant Professor) |
| qualification | VARCHAR(255) | NULL | Academic qualifications |
| specialization | TEXT | NULL | Research/teaching specialization |
| joining_date | DATE | NULL | Date of joining |
| office_room | VARCHAR(50) | NULL | Office room number |
| office_hours | TEXT | NULL | Available office hours |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 2.4 `students`
Student profiles (extends users).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| student_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id (1:1), NOT NULL | User account reference |
| registration_number | VARCHAR(50) | UNIQUE, NOT NULL | University registration number |
| roll_number | VARCHAR(50) | NULL | Class roll number |
| department_id | UUID | FK -> departments.department_id | Department |
| program_id | UUID | FK -> programs.program_id | Enrolled program |
| semester | INTEGER | NOT NULL | Current semester (1-8) |
| batch_year | INTEGER | NOT NULL | Admission batch year |
| admission_date | DATE | NULL | Date of admission |
| current_cgpa | NUMERIC(3,2) | NULL | Current CGPA (0.00-10.00) |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

**Indexes**: user_id, registration_number, department_id, program_id, semester, batch_year

---

## 3. COURSE MANAGEMENT SERVICE

### Purpose
Manages courses, subjects, enrollments, and academic calendar.

### Tables

#### 3.1 `subjects`
Academic subjects/modules.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| subject_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| subject_name | VARCHAR(255) | NOT NULL | Full subject name |
| subject_code | VARCHAR(20) | UNIQUE, NOT NULL | Subject code (CS101) |
| department_id | UUID | FK -> departments.department_id | Offering department |
| credits | INTEGER | NOT NULL | Credit hours |
| subject_type | VARCHAR(50) | NULL | theory, practical, project |
| description | TEXT | NULL | Subject description |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 3.2 `courses`
Course instances for a semester.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| course_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| course_code | VARCHAR(20) | UNIQUE, NOT NULL | Unique course code |
| course_name | VARCHAR(255) | NOT NULL | Course name |
| subject_id | UUID | FK -> subjects.subject_id, NOT NULL | Base subject |
| department_id | UUID | FK -> departments.department_id | Department |
| program_id | UUID | FK -> programs.program_id | Program |
| semester | INTEGER | NOT NULL | Target semester |
| academic_year | INTEGER | NOT NULL | Academic year (2024) |
| max_students | INTEGER | NULL | Maximum enrollment |
| description | TEXT | NULL | Course description |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 3.3 `faculty_courses`
Faculty-to-course assignments.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| faculty_course_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| faculty_id | UUID | FK -> faculties.faculty_id, NOT NULL | Faculty member |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Course |
| role | VARCHAR(50) | DEFAULT 'instructor' | instructor, co-instructor, teaching_assistant |
| assigned_at | TIMESTAMPTZ | DEFAULT now() | Assignment timestamp |

**Unique Index**: (faculty_id, course_id)

---

#### 3.4 `course_enrollments`
Student course enrollments.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| enrollment_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| student_id | UUID | FK -> students.student_id, NOT NULL | Student |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Course |
| enrollment_status | VARCHAR(20) | DEFAULT 'enrolled' | enrolled, dropped, completed, failed |
| enrollment_date | TIMESTAMPTZ | DEFAULT now() | Enrollment timestamp |
| dropped_date | TIMESTAMPTZ | NULL | Drop timestamp (if applicable) |
| completion_date | TIMESTAMPTZ | NULL | Completion timestamp |

**Unique Index**: (student_id, course_id)

---

#### 3.5 `semesters`
Academic semester definitions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| semester_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| semester_name | VARCHAR(50) | NOT NULL | Fall 2024, Spring 2025 |
| academic_year | INTEGER | NOT NULL | Academic year |
| start_date | DATE | NOT NULL | Semester start date |
| end_date | DATE | NOT NULL | Semester end date |
| is_current | BOOLEAN | DEFAULT false | Current semester flag |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

## 4. CONTENT MANAGEMENT SERVICE

### Purpose
Manages content metadata, permissions, and organization. Actual file content stored in MongoDB.

### Tables

#### 4.1 `content_metadata`
Metadata for uploaded content (files stored in MongoDB).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| content_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| mongo_document_id | VARCHAR(24) | UNIQUE, NOT NULL | MongoDB ObjectId reference |
| title | VARCHAR(255) | NOT NULL | Content title |
| description | TEXT | NULL | Content description |
| content_type | VARCHAR(50) | NOT NULL | document, video, image, link, presentation |
| file_size | BIGINT | NULL | File size in bytes |
| mime_type | VARCHAR(100) | NULL | MIME type (application/pdf) |
| version | INTEGER | DEFAULT 1 | Current version number |
| is_published | BOOLEAN | DEFAULT false | Publication status |
| published_at | TIMESTAMPTZ | NULL | Publication timestamp |
| scheduled_publish_at | TIMESTAMPTZ | NULL | Scheduled publication time |
| expires_at | TIMESTAMPTZ | NULL | Content expiration time |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator user |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 4.2 `content_course_mapping`
Maps content to courses/subjects.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| mapping_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| content_id | UUID | FK -> content_metadata.content_id, NOT NULL | Content reference |
| course_id | UUID | FK -> courses.course_id | Course reference |
| subject_id | UUID | FK -> subjects.subject_id | Subject reference |
| semester | INTEGER | NULL | Target semester |
| topic | VARCHAR(255) | NULL | Topic/chapter name |
| keywords | TEXT | NULL | Comma-separated keywords |
| language | VARCHAR(50) | DEFAULT 'English' | Content language |

---

#### 4.3 `content_folders`
Hierarchical folder organization.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| folder_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| folder_name | VARCHAR(255) | NOT NULL | Folder name |
| parent_folder_id | UUID | FK -> content_folders.folder_id | Parent folder (NULL for root) |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator user |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 4.4 `content_permissions`
Access control for content.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| permission_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| content_id | UUID | FK -> content_metadata.content_id, NOT NULL | Content reference |
| user_id | UUID | FK -> users.user_id | Specific user access |
| role_id | UUID | FK -> roles.role_id | Role-based access |
| department_id | UUID | FK -> departments.department_id | Department access |
| course_id | UUID | FK -> courses.course_id | Course-based access |
| group_id | UUID | FK -> groups.group_id | Group-based access |
| permission_type | VARCHAR(20) | NOT NULL | view, download, edit |
| granted_by | UUID | FK -> users.user_id, NOT NULL | Granting user |
| granted_at | TIMESTAMPTZ | DEFAULT now() | Grant timestamp |

---

## 5. CONTENT TRACKING SERVICE

### Purpose
Tracks content access, engagement, and delivery status.

### Tables

#### 5.1 `content_access_logs`
Access log for content views/downloads.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| access_log_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| content_id | UUID | FK -> content_metadata.content_id, NOT NULL | Content accessed |
| user_id | UUID | FK -> users.user_id, NOT NULL | Accessing user |
| access_type | VARCHAR(20) | NOT NULL | view, download, share |
| ip_address | VARCHAR(45) | NULL | Client IP address |
| user_agent | TEXT | NULL | Client user agent |
| device_type | VARCHAR(50) | NULL | desktop, mobile, tablet |
| accessed_at | TIMESTAMPTZ | DEFAULT now() | Access timestamp |

---

#### 5.2 `content_engagement_metrics`
User engagement with content.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| metric_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| content_id | UUID | FK -> content_metadata.content_id, NOT NULL | Content reference |
| user_id | UUID | FK -> users.user_id, NOT NULL | User reference |
| time_spent_seconds | INTEGER | DEFAULT 0 | Total time spent |
| completion_percentage | INTEGER | DEFAULT 0 | Content completion (0-100) |
| last_accessed_at | TIMESTAMPTZ | NULL | Last access timestamp |
| first_accessed_at | TIMESTAMPTZ | DEFAULT now() | First access timestamp |

---

#### 5.3 `content_bookmarks`
User bookmarks for quick access.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| bookmark_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| content_id | UUID | FK -> content_metadata.content_id, NOT NULL | Content reference |
| user_id | UUID | FK -> users.user_id, NOT NULL | User reference |
| notes | TEXT | NULL | User notes |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

**Unique Index**: (content_id, user_id)

---

## 6. TIMETABLE MANAGEMENT SERVICE

### Purpose
Manages class schedules, room allocations, and schedule change requests.

### Tables

#### 6.1 `rooms`
Classrooms and facilities.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| room_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| room_number | VARCHAR(50) | UNIQUE, NOT NULL | Room number (A101) |
| building | VARCHAR(100) | NULL | Building name |
| floor | INTEGER | NULL | Floor number |
| capacity | INTEGER | NOT NULL | Seating capacity |
| room_type | VARCHAR(50) | NULL | classroom, lab, auditorium, seminar_hall |
| facilities | TEXT | NULL | Available facilities (projector, AC) |
| is_available | BOOLEAN | DEFAULT true | Availability status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 6.2 `time_slots`
Standard time slots for scheduling.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| slot_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| slot_name | VARCHAR(50) | NOT NULL | Slot name (Period 1) |
| start_time | TIME | NOT NULL | Start time |
| end_time | TIME | NOT NULL | End time |
| duration_minutes | INTEGER | NOT NULL | Duration in minutes |
| slot_order | INTEGER | NULL | Display order |

---

#### 6.3 `timetables`
Timetable headers.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| timetable_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| semester_id | UUID | FK -> semesters.semester_id, NOT NULL | Semester reference |
| department_id | UUID | FK -> departments.department_id | Department |
| program_id | UUID | FK -> programs.program_id | Program |
| semester_number | INTEGER | NULL | Semester number (1-8) |
| is_published | BOOLEAN | DEFAULT false | Publication status |
| published_at | TIMESTAMPTZ | NULL | Publication timestamp |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 6.4 `timetable_entries`
Individual schedule entries.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| entry_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| timetable_id | UUID | FK -> timetables.timetable_id, NOT NULL | Timetable reference |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Course |
| faculty_id | UUID | FK -> faculties.faculty_id, NOT NULL | Assigned faculty |
| room_id | UUID | FK -> rooms.room_id, NOT NULL | Assigned room |
| time_slot_id | UUID | FK -> time_slots.slot_id, NOT NULL | Time slot |
| day_of_week | INTEGER | NOT NULL | 1=Monday through 7=Sunday |
| class_type | VARCHAR(50) | DEFAULT 'lecture' | lecture, lab, tutorial, practical |
| effective_from | DATE | NULL | Start date |
| effective_to | DATE | NULL | End date |
| is_active | BOOLEAN | DEFAULT true | Active status |

---

## 7. ATTENDANCE SERVICE

### Purpose
Tracks student attendance for classes.

### Tables

#### 7.1 `attendance_records`
Individual attendance records.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| attendance_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| entry_id | UUID | FK -> timetable_entries.entry_id, NOT NULL | Timetable entry |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Course |
| student_id | UUID | FK -> students.student_id, NOT NULL | Student |
| attendance_date | DATE | NOT NULL | Class date |
| status | VARCHAR(20) | NOT NULL | present, absent, late, excused |
| marked_by | UUID | FK -> faculties.faculty_id, NOT NULL | Faculty who marked |
| marked_at | TIMESTAMPTZ | DEFAULT now() | Marking timestamp |
| remarks | TEXT | NULL | Additional remarks |

---

#### 7.2 `attendance_summary`
Aggregated attendance per student per course.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| summary_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| student_id | UUID | FK -> students.student_id, NOT NULL | Student |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Course |
| total_classes | INTEGER | DEFAULT 0 | Total classes held |
| classes_attended | INTEGER | DEFAULT 0 | Classes attended |
| attendance_percentage | NUMERIC(5,2) | NULL | Percentage (0.00-100.00) |
| last_updated | TIMESTAMPTZ | DEFAULT now() | Last calculation time |

**Unique Index**: (student_id, course_id)

---

## 8. NOTIFICATION SERVICE

### Purpose
Manages notifications across multiple channels (email, SMS, push, in-app).

### Tables

#### 8.1 `notification_templates`
Reusable notification templates.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| template_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| template_name | VARCHAR(100) | UNIQUE, NOT NULL | Template name |
| template_type | VARCHAR(50) | NOT NULL | email, sms, push, in_app |
| subject | VARCHAR(255) | NULL | Email subject line |
| body | TEXT | NOT NULL | Message body with placeholders |
| variables | JSONB | NULL | Variable definitions ({{user_name}}) |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 8.2 `notifications`
Notification instances.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| notification_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id, NOT NULL | Target user |
| notification_type | VARCHAR(50) | NOT NULL | assignment, grade, announcement, reminder |
| title | VARCHAR(255) | NOT NULL | Notification title |
| message | TEXT | NOT NULL | Message content |
| priority | VARCHAR(20) | DEFAULT 'normal' | low, normal, high, urgent |
| status | VARCHAR(20) | DEFAULT 'unread' | unread, read, archived |
| action_url | TEXT | NULL | Deep link URL |
| metadata | JSONB | NULL | Additional data |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| read_at | TIMESTAMPTZ | NULL | Read timestamp |

---

#### 8.3 `notification_queue`
Delivery queue for notifications.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| queue_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| notification_id | UUID | FK -> notifications.notification_id, NOT NULL | Notification reference |
| channel | VARCHAR(50) | NOT NULL | email, sms, push, in_app |
| recipient_address | TEXT | NOT NULL | Delivery address |
| status | VARCHAR(20) | DEFAULT 'pending' | pending, sent, failed, retry |
| attempt_count | INTEGER | DEFAULT 0 | Delivery attempts |
| last_attempt_at | TIMESTAMPTZ | NULL | Last attempt timestamp |
| sent_at | TIMESTAMPTZ | NULL | Successful delivery time |
| error_message | TEXT | NULL | Error details |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

#### 8.4 `notification_preferences`
User notification preferences.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| preference_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id (1:1), NOT NULL | User reference |
| notification_type | VARCHAR(50) | NOT NULL | Notification category |
| email_enabled | BOOLEAN | DEFAULT true | Email delivery |
| sms_enabled | BOOLEAN | DEFAULT false | SMS delivery |
| push_enabled | BOOLEAN | DEFAULT true | Push notification |
| in_app_enabled | BOOLEAN | DEFAULT true | In-app notification |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

**Unique Index**: (user_id, notification_type)

---

## 9. ANNOUNCEMENT SERVICE

### Purpose
Manages university-wide and targeted announcements.

### Tables

#### 9.1 `announcements`
Announcement records.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| announcement_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| title | VARCHAR(255) | NOT NULL | Announcement title |
| content | TEXT | NOT NULL | Announcement content |
| announcement_type | VARCHAR(50) | NOT NULL | general, academic, event, urgent |
| priority | VARCHAR(20) | DEFAULT 'normal' | normal, high, urgent |
| target_audience | VARCHAR(50) | NOT NULL | all, students, faculty, department, course |
| is_published | BOOLEAN | DEFAULT false | Publication status |
| published_at | TIMESTAMP | NULL | Publication timestamp |
| scheduled_publish_at | TIMESTAMP | NULL | Scheduled publication |
| expires_at | TIMESTAMP | NULL | Expiration timestamp |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 9.2 `announcement_targets`
Specific targeting for announcements.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| target_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| announcement_id | UUID | FK -> announcements.announcement_id, NOT NULL | Announcement |
| target_type | VARCHAR(50) | NOT NULL | department, program, course, group, individual |
| department_id | UUID | FK -> departments.department_id | Target department |
| program_id | UUID | FK -> programs.program_id | Target program |
| course_id | UUID | FK -> courses.course_id | Target course |
| group_id | UUID | FK -> groups.group_id | Target group |
| user_id | UUID | FK -> users.user_id | Target user |

---

## 10. COMMUNICATION SERVICE

### Purpose
Manages messaging, discussions, and office hours.

### Tables

#### 10.1 `conversations`
Message threads/conversations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| conversation_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| conversation_type | VARCHAR(50) | NOT NULL | direct, group, course |
| title | VARCHAR(255) | NULL | Conversation title |
| course_id | UUID | FK -> courses.course_id | Associated course |
| group_id | UUID | FK -> groups.group_id | Associated group |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| last_message_at | TIMESTAMPTZ | NULL | Last message timestamp |

---

#### 10.2 `messages`
Individual messages.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| message_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| conversation_id | UUID | FK -> conversations.conversation_id, NOT NULL | Conversation |
| sender_id | UUID | FK -> users.user_id, NOT NULL | Sender |
| message_text | TEXT | NOT NULL | Message content |
| message_type | VARCHAR(50) | DEFAULT 'text' | text, file, image, link |
| attachment_url | TEXT | NULL | Attachment URL |
| reply_to_message_id | UUID | FK -> messages.message_id | Reply reference |
| is_edited | BOOLEAN | DEFAULT false | Edit status |
| edited_at | TIMESTAMPTZ | NULL | Edit timestamp |
| sent_at | TIMESTAMPTZ | DEFAULT now() | Send timestamp |

---

#### 10.3 `discussion_forums`
Course discussion forums.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| forum_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| course_id | UUID | FK -> courses.course_id, NOT NULL | Associated course |
| title | VARCHAR(255) | NOT NULL | Forum title |
| description | TEXT | NULL | Forum description |
| created_by | UUID | FK -> faculties.faculty_id, NOT NULL | Creator faculty |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |

---

## 11. GROUP MANAGEMENT

### Purpose
Manages study groups, project groups, and sections.

### Tables

#### 11.1 `groups`
User groups for content sharing and communication.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| group_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| group_name | VARCHAR(255) | NOT NULL | Group name |
| group_type | VARCHAR(50) | NOT NULL | study_group, project_group, section, batch |
| description | TEXT | NULL | Group description |
| course_id | UUID | FK -> courses.course_id | Associated course |
| department_id | UUID | FK -> departments.department_id | Associated department |
| program_id | UUID | FK -> programs.program_id | Associated program |
| semester | INTEGER | NULL | Associated semester |
| created_by | UUID | FK -> users.user_id, NOT NULL | Creator |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMPTZ | DEFAULT now() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT now() | Last update time |

---

#### 11.2 `group_members`
Group membership.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| member_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| group_id | UUID | FK -> groups.group_id, NOT NULL | Group reference |
| user_id | UUID | FK -> users.user_id, NOT NULL | Member user |
| role | VARCHAR(50) | DEFAULT 'member' | admin, moderator, member |
| joined_at | TIMESTAMPTZ | DEFAULT now() | Join timestamp |
| left_at | TIMESTAMPTZ | NULL | Leave timestamp |
| is_active | BOOLEAN | DEFAULT true | Active status |

**Unique Index**: (group_id, user_id)

---

## 12. SESSION MONITORING SERVICE

### Purpose
Tracks active users and system metrics.

### Tables

#### 12.1 `active_sessions`
Current user sessions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| session_id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| user_id | UUID | FK -> users.user_id, NOT NULL | Session user |
| session_token | VARCHAR(255) | UNIQUE, NOT NULL | JWT/session token |
| ip_address | VARCHAR(45) | NULL | Client IP |
| user_agent | TEXT | NULL | Browser/client info |
| device_type | VARCHAR(50) | NULL | desktop, mobile, tablet |
| started_at | TIMESTAMPTZ | DEFAULT now() | Session start |
| last_activity_at | TIMESTAMPTZ | DEFAULT now() | Last activity |
| ended_at | TIMESTAMPTZ | NULL | Session end |
| is_active | BOOLEAN | DEFAULT true | Active status |

---

## MongoDB Collections (Document Storage)

### `content_files`
Stores actual file content using GridFS.

```json
{
  "_id": ObjectId("..."),
  "filename": "lecture_notes.pdf",
  "contentType": "application/pdf",
  "length": 1048576,
  "uploadDate": ISODate("2024-01-15T10:30:00Z"),
  "metadata": {
    "postgres_content_id": "uuid-string",
    "original_name": "Lecture Notes - Week 1.pdf",
    "uploaded_by": "user-uuid",
    "checksum": "sha256-hash"
  }
}
```

### `content_versions`
Stores versioned file content.

```json
{
  "_id": ObjectId("..."),
  "content_id": "postgres-uuid",
  "version": 2,
  "file_data": Binary("..."),
  "created_at": ISODate("2024-01-20T14:00:00Z"),
  "created_by": "user-uuid",
  "change_description": "Updated formatting"
}
```

---

## Data Flow Examples

### 1. User Registration Flow
```
1. Create user in `users` table (credentials)
2. Create `user_profiles` entry (personal info)
3. Create `students` or `faculties` entry based on role
4. Create default `notification_preferences`
5. Publish USER_CREATED event to Kafka
```

### 2. Content Upload Flow
```
1. Upload file to MongoDB GridFS -> get ObjectId
2. Create `content_metadata` entry with mongo_document_id
3. Create `content_course_mapping` for course association
4. Create `content_permissions` for access control
5. Publish CONTENT_CREATED event to Kafka
```

### 3. Attendance Marking Flow
```
1. Faculty selects timetable_entry for class
2. Create `attendance_records` for each student
3. Update `attendance_summary` aggregate
4. Publish ATTENDANCE_MARKED event to Kafka
```

---

## Indexing Strategy

### Primary Indexes
- All tables use UUID primary keys with `gen_random_uuid()`
- Supports distributed systems and avoids ID conflicts

### Foreign Key Indexes
- All foreign keys are indexed for join performance

### Query-Optimized Indexes
- `users.email` - Login queries
- `users.register_no` - Student lookup
- `content_metadata.created_at` - Recent content
- `notifications.user_id + status` - Unread notifications
- `attendance_records.course_id + attendance_date` - Daily attendance

### Composite Indexes
- `(student_id, course_id)` - Enrollment lookup
- `(content_id, user_id)` - Bookmark/engagement queries
- `(role_id, permission_id)` - RBAC queries

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024-12-27 | Initial schema design |
