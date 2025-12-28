# NimbusU Course Service - API Documentation

This document defines the REST API for the **Course Service**, managing academic structure, courses, enrollments, and faculty assignments.

**Base URL:** `/api/v1`  
**Port:** 8084

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [Departments](#2-departments)
3. [Programs](#3-programs)
4. [Subjects](#4-subjects)
5. [Semesters](#5-semesters)
6. [Courses](#6-courses)
7. [Faculty Assignments](#7-faculty-assignments)
8. [Enrollments](#8-enrollments)
9. [Faculty Profiles](#9-faculty-profiles)
10. [Student Profiles](#10-student-profiles)
11. [Academic Calendar](#11-academic-calendar)
12. [Kafka Events](#12-kafka-events)
13. [Error Responses](#13-error-responses)
14. [Permissions (RBAC)](#14-permissions-rbac)

---

## 1. Authentication

All endpoints except public ones require a valid JWT token:

```
Authorization: Bearer <token>
```

The token is validated by the API Gateway and user context is passed via headers:
- `X-User-ID`: Authenticated user's UUID
- `X-User-Role`: User's role (admin, faculty, student)

---

## 2. Departments

### 2.1. List Departments

Get all departments with optional filtering.

- **GET** `/departments`
- **Auth:** Public
- **Query Params:**
  - `is_active` (boolean): Filter by active status
  - `page` (int): Page number (default: 1)
  - `limit` (int): Items per page (default: 20, max: 100)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "department_id": "uuid",
      "department_name": "Computer Science and Engineering",
      "department_code": "CSE",
      "head_of_department": {
        "faculty_id": "uuid",
        "name": "Dr. John Smith"
      },
      "description": "Department of Computer Science",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

### 2.2. Get Department by ID

- **GET** `/departments/{department_id}`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "department_id": "uuid",
  "department_name": "Computer Science and Engineering",
  "department_code": "CSE",
  "head_of_department": {
    "faculty_id": "uuid",
    "name": "Dr. John Smith",
    "designation": "Professor"
  },
  "description": "Department of Computer Science",
  "programs_count": 5,
  "faculty_count": 25,
  "students_count": 500,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### 2.3. Create Department

- **POST** `/departments`
- **Auth:** Admin only

**Request:**

```json
{
  "department_name": "Computer Science and Engineering",
  "department_code": "CSE",
  "head_of_department": "faculty-uuid",
  "description": "Department of Computer Science"
}
```

**Response:** `201 Created`

```json
{
  "department_id": "uuid",
  "department_name": "Computer Science and Engineering",
  "department_code": "CSE",
  "message": "Department created successfully"
}
```

---

### 2.4. Update Department

- **PUT** `/departments/{department_id}`
- **Auth:** Admin only

**Request:**

```json
{
  "department_name": "Computer Science & Engineering",
  "head_of_department": "new-faculty-uuid",
  "description": "Updated description",
  "is_active": true
}
```

**Response:** `200 OK`

---

### 2.5. Delete Department

Soft delete a department.

- **DELETE** `/departments/{department_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

---

## 3. Programs

### 3.1. List Programs

- **GET** `/programs`
- **Auth:** Public
- **Query Params:**
  - `department_id` (uuid): Filter by department
  - `degree_type` (string): Filter by degree type (Bachelor, Master, PhD)
  - `is_active` (boolean): Filter by active status
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "program_id": "uuid",
      "program_name": "Bachelor of Technology in Computer Science",
      "program_code": "BTCS",
      "department": {
        "department_id": "uuid",
        "department_name": "Computer Science and Engineering"
      },
      "degree_type": "Bachelor",
      "duration_years": 4,
      "total_credits": 160,
      "is_active": true
    }
  ],
  "pagination": { ... }
}
```

---

### 3.2. Get Program by ID

- **GET** `/programs/{program_id}`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "program_id": "uuid",
  "program_name": "Bachelor of Technology in Computer Science",
  "program_code": "BTCS",
  "department": {
    "department_id": "uuid",
    "department_name": "Computer Science and Engineering"
  },
  "degree_type": "Bachelor",
  "duration_years": 4,
  "total_credits": 160,
  "description": "4-year undergraduate program in Computer Science",
  "semester_count": 8,
  "students_count": 240,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### 3.3. Create Program

- **POST** `/programs`
- **Auth:** Admin only

**Request:**

```json
{
  "program_name": "Bachelor of Technology in Computer Science",
  "program_code": "BTCS",
  "department_id": "uuid",
  "degree_type": "Bachelor",
  "duration_years": 4,
  "total_credits": 160,
  "description": "4-year undergraduate program"
}
```

**Response:** `201 Created`

---

### 3.4. Update Program

- **PUT** `/programs/{program_id}`
- **Auth:** Admin only

**Response:** `200 OK`

---

### 3.5. Delete Program

- **DELETE** `/programs/{program_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

---

## 4. Subjects

### 4.1. List Subjects

- **GET** `/subjects`
- **Auth:** Public
- **Query Params:**
  - `department_id` (uuid): Filter by department
  - `subject_type` (string): theory, practical, project
  - `credits` (int): Filter by credit hours
  - `search` (string): Search by name or code
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "subject_id": "uuid",
      "subject_name": "Data Structures and Algorithms",
      "subject_code": "CS201",
      "department": {
        "department_id": "uuid",
        "department_name": "Computer Science and Engineering"
      },
      "credits": 4,
      "subject_type": "theory",
      "is_active": true
    }
  ],
  "pagination": { ... }
}
```

---

### 4.2. Get Subject by ID

- **GET** `/subjects/{subject_id}`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "subject_id": "uuid",
  "subject_name": "Data Structures and Algorithms",
  "subject_code": "CS201",
  "department": {
    "department_id": "uuid",
    "department_name": "Computer Science and Engineering"
  },
  "credits": 4,
  "subject_type": "theory",
  "description": "Fundamental data structures and algorithms",
  "syllabus": "Week 1: Arrays and Lists...",
  "prerequisites": [
    {
      "subject_id": "uuid",
      "subject_code": "CS101",
      "subject_name": "Introduction to Programming",
      "is_mandatory": true
    }
  ],
  "corequisites": [],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### 4.3. Create Subject

- **POST** `/subjects`
- **Auth:** Admin, Faculty (own department)

**Request:**

```json
{
  "subject_name": "Data Structures and Algorithms",
  "subject_code": "CS201",
  "department_id": "uuid",
  "credits": 4,
  "subject_type": "theory",
  "description": "Fundamental data structures",
  "syllabus": "Week 1: Arrays...",
  "prerequisites": [
    { "subject_id": "uuid", "is_mandatory": true }
  ],
  "corequisites": ["uuid"]
}
```

**Response:** `201 Created`

---

### 4.4. Update Subject

- **PUT** `/subjects/{subject_id}`
- **Auth:** Admin, Faculty (own department)

**Response:** `200 OK`

---

### 4.5. Delete Subject

- **DELETE** `/subjects/{subject_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

---

## 5. Semesters

### 5.1. List Semesters

- **GET** `/semesters`
- **Auth:** Public
- **Query Params:**
  - `academic_year` (int): Filter by year
  - `is_current` (boolean): Get current semester only
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "semester_id": "uuid",
      "semester_name": "Fall 2024",
      "semester_code": "F2024",
      "academic_year": 2024,
      "start_date": "2024-08-01",
      "end_date": "2024-12-15",
      "registration_start": "2024-07-15",
      "registration_end": "2024-08-05",
      "is_current": true
    }
  ],
  "pagination": { ... }
}
```

---

### 5.2. Get Current Semester

- **GET** `/semesters/current`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "semester_id": "uuid",
  "semester_name": "Fall 2024",
  "semester_code": "F2024",
  "academic_year": 2024,
  "start_date": "2024-08-01",
  "end_date": "2024-12-15",
  "is_current": true,
  "days_remaining": 45,
  "registration_open": false
}
```

---

### 5.3. Create Semester

- **POST** `/semesters`
- **Auth:** Admin only

**Request:**

```json
{
  "semester_name": "Spring 2025",
  "semester_code": "S2025",
  "academic_year": 2025,
  "start_date": "2025-01-15",
  "end_date": "2025-05-30",
  "registration_start": "2025-01-01",
  "registration_end": "2025-01-20"
}
```

**Response:** `201 Created`

---

### 5.4. Update Semester

- **PUT** `/semesters/{semester_id}`
- **Auth:** Admin only

**Response:** `200 OK`

---

### 5.5. Set Current Semester

- **POST** `/semesters/{semester_id}/set-current`
- **Auth:** Admin only

**Response:** `200 OK`

---

## 6. Courses

### 6.1. List Courses

- **GET** `/courses`
- **Auth:** Public (basic info), Authenticated (full info)
- **Query Params:**
  - `department_id` (uuid): Filter by department
  - `program_id` (uuid): Filter by program
  - `semester_id` (uuid): Filter by semester
  - `semester_number` (int): Filter by target semester (1-8)
  - `subject_id` (uuid): Filter by subject
  - `faculty_id` (uuid): Filter by assigned faculty
  - `status` (string): draft, active, completed, cancelled
  - `search` (string): Search by name or code
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "course_id": "uuid",
      "course_code": "CS201-F2024",
      "course_name": "Data Structures and Algorithms",
      "subject": {
        "subject_id": "uuid",
        "subject_code": "CS201"
      },
      "department": {
        "department_id": "uuid",
        "department_code": "CSE"
      },
      "semester": {
        "semester_id": "uuid",
        "semester_name": "Fall 2024"
      },
      "semester_number": 3,
      "academic_year": 2024,
      "faculty": [
        {
          "faculty_id": "uuid",
          "name": "Dr. Jane Doe",
          "role": "instructor",
          "is_primary": true
        }
      ],
      "max_students": 60,
      "current_enrollment": 45,
      "status": "active",
      "is_active": true
    }
  ],
  "pagination": { ... }
}
```

---

### 6.2. Get Course by ID

- **GET** `/courses/{course_id}`
- **Auth:** Public (basic), Enrolled/Faculty (full)

**Response:** `200 OK`

```json
{
  "course_id": "uuid",
  "course_code": "CS201-F2024",
  "course_name": "Data Structures and Algorithms",
  "subject": {
    "subject_id": "uuid",
    "subject_code": "CS201",
    "subject_name": "Data Structures and Algorithms",
    "credits": 4,
    "subject_type": "theory"
  },
  "department": {
    "department_id": "uuid",
    "department_name": "Computer Science and Engineering"
  },
  "program": {
    "program_id": "uuid",
    "program_name": "B.Tech Computer Science"
  },
  "semester": {
    "semester_id": "uuid",
    "semester_name": "Fall 2024",
    "start_date": "2024-08-01",
    "end_date": "2024-12-15"
  },
  "semester_number": 3,
  "academic_year": 2024,
  "faculty": [
    {
      "faculty_id": "uuid",
      "user_id": "uuid",
      "name": "Dr. Jane Doe",
      "designation": "Associate Professor",
      "role": "instructor",
      "is_primary": true
    }
  ],
  "max_students": 60,
  "current_enrollment": 45,
  "available_seats": 15,
  "status": "active",
  "description": "In-depth study of data structures",
  "prerequisites": [
    {
      "subject_code": "CS101",
      "subject_name": "Introduction to Programming"
    }
  ],
  "created_by": "admin-uuid",
  "is_active": true,
  "created_at": "2024-07-15T00:00:00Z",
  "updated_at": "2024-08-01T00:00:00Z"
}
```

---

### 6.3. Create Course

- **POST** `/courses`
- **Auth:** Admin, Faculty (own department)

**Request:**

```json
{
  "course_code": "CS201-F2024",
  "course_name": "Data Structures and Algorithms",
  "subject_id": "uuid",
  "department_id": "uuid",
  "program_id": "uuid",
  "semester_id": "uuid",
  "semester_number": 3,
  "academic_year": 2024,
  "max_students": 60,
  "description": "In-depth study of data structures"
}
```

**Response:** `201 Created`

```json
{
  "course_id": "uuid",
  "course_code": "CS201-F2024",
  "message": "Course created successfully"
}
```

**Kafka Event Published:** `COURSE_CREATED` to `course.events`

---

### 6.4. Update Course

- **PUT** `/courses/{course_id}`
- **Auth:** Admin, Assigned Faculty

**Request:**

```json
{
  "course_name": "Data Structures and Algorithms (Updated)",
  "max_students": 80,
  "description": "Updated description",
  "status": "active"
}
```

**Response:** `200 OK`

**Kafka Event Published:** `COURSE_UPDATED` to `course.events`

---

### 6.5. Activate Course

- **POST** `/courses/{course_id}/activate`
- **Auth:** Admin, Primary Faculty

**Response:** `200 OK`

**Kafka Event Published:** `COURSE_ACTIVATED` to `course.events`

---

### 6.6. Deactivate Course

- **POST** `/courses/{course_id}/deactivate`
- **Auth:** Admin only

**Response:** `200 OK`

**Kafka Event Published:** `COURSE_DEACTIVATED` to `course.events`

---

### 6.7. Delete Course

Soft delete a course.

- **DELETE** `/courses/{course_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

**Kafka Event Published:** `COURSE_DELETED` to `course.events`

---

### 6.8. Get Course Students

- **GET** `/courses/{course_id}/students`
- **Auth:** Admin, Assigned Faculty

**Query Params:**
- `status` (string): enrolled, waitlisted, dropped, completed, failed
- `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "enrollment_id": "uuid",
      "student": {
        "student_id": "uuid",
        "user_id": "uuid",
        "name": "John Student",
        "registration_number": "2022CS001",
        "email": "john@university.edu"
      },
      "enrollment_status": "enrolled",
      "enrollment_date": "2024-07-20T10:30:00Z",
      "grade": null
    }
  ],
  "summary": {
    "total_enrolled": 45,
    "total_waitlisted": 5,
    "total_dropped": 2,
    "total_completed": 0
  },
  "pagination": { ... }
}
```

---

## 7. Faculty Assignments

### 7.1. Assign Faculty to Course

- **POST** `/courses/{course_id}/faculty`
- **Auth:** Admin only

**Request:**

```json
{
  "faculty_id": "uuid",
  "role": "instructor",
  "is_primary": true
}
```

**Response:** `201 Created`

```json
{
  "faculty_course_id": "uuid",
  "message": "Faculty assigned successfully"
}
```

**Kafka Event Published:** `FACULTY_ASSIGNED` to `course.events`

---

### 7.2. List Course Faculty

- **GET** `/courses/{course_id}/faculty`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "data": [
    {
      "faculty_course_id": "uuid",
      "faculty": {
        "faculty_id": "uuid",
        "user_id": "uuid",
        "name": "Dr. Jane Doe",
        "designation": "Associate Professor",
        "email": "jane@university.edu"
      },
      "role": "instructor",
      "is_primary": true,
      "assigned_at": "2024-07-15T00:00:00Z"
    }
  ]
}
```

---

### 7.3. Update Faculty Assignment

- **PUT** `/courses/{course_id}/faculty/{faculty_id}`
- **Auth:** Admin only

**Request:**

```json
{
  "role": "co-instructor",
  "is_primary": false
}
```

**Response:** `200 OK`

---

### 7.4. Remove Faculty from Course

- **DELETE** `/courses/{course_id}/faculty/{faculty_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

**Kafka Event Published:** `FACULTY_UNASSIGNED` to `course.events`

---

## 8. Enrollments

### 8.1. Enroll in Course

- **POST** `/courses/{course_id}/enroll`
- **Auth:** Student (self), Admin (any student)

**Request:**

```json
{
  "student_id": "uuid"
}
```

> Note: `student_id` is optional for self-enrollment (derived from token)

**Response:** `201 Created`

```json
{
  "enrollment_id": "uuid",
  "enrollment_status": "enrolled",
  "message": "Successfully enrolled in course"
}
```

OR (if course is full):

```json
{
  "enrollment_id": "uuid",
  "enrollment_status": "waitlisted",
  "waitlist_position": 3,
  "message": "Added to waitlist"
}
```

**Kafka Event Published:** `STUDENT_ENROLLED` or `WAITLIST_ADDED` to `enrollment.events`

---

### 8.2. Drop Course

- **POST** `/courses/{course_id}/drop`
- **Auth:** Student (self), Admin (any student)

**Request:**

```json
{
  "student_id": "uuid",
  "reason": "Schedule conflict"
}
```

**Response:** `200 OK`

```json
{
  "message": "Successfully dropped from course"
}
```

**Kafka Event Published:** `STUDENT_DROPPED` to `enrollment.events`

---

### 8.3. Get Student Enrollments

- **GET** `/students/{student_id}/enrollments`
- **Auth:** Student (self), Faculty (if teaching), Admin

**Query Params:**
- `semester_id` (uuid): Filter by semester
- `status` (string): enrolled, waitlisted, dropped, completed, failed
- `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "enrollment_id": "uuid",
      "course": {
        "course_id": "uuid",
        "course_code": "CS201-F2024",
        "course_name": "Data Structures and Algorithms",
        "credits": 4
      },
      "semester": {
        "semester_id": "uuid",
        "semester_name": "Fall 2024"
      },
      "enrollment_status": "enrolled",
      "enrollment_date": "2024-07-20T10:30:00Z",
      "grade": null,
      "grade_points": null
    }
  ],
  "summary": {
    "total_credits_enrolled": 18,
    "courses_count": 5
  },
  "pagination": { ... }
}
```

---

### 8.4. Update Enrollment (Grade)

- **PUT** `/enrollments/{enrollment_id}`
- **Auth:** Faculty (assigned), Admin

**Request:**

```json
{
  "enrollment_status": "completed",
  "grade": "A",
  "grade_points": 9.0
}
```

**Response:** `200 OK`

**Kafka Event Published:** `ENROLLMENT_COMPLETED` or `ENROLLMENT_FAILED` to `enrollment.events`

---

### 8.5. Bulk Enroll Students

- **POST** `/courses/{course_id}/bulk-enroll`
- **Auth:** Admin only

**Request:**

```json
{
  "student_ids": ["uuid1", "uuid2", "uuid3"],
  "skip_prerequisites": false
}
```

**Response:** `200 OK`

```json
{
  "successful": [
    { "student_id": "uuid1", "enrollment_id": "uuid", "status": "enrolled" }
  ],
  "failed": [
    { "student_id": "uuid2", "error": "Prerequisites not met" }
  ],
  "summary": {
    "total": 3,
    "successful": 2,
    "failed": 1
  }
}
```

---

## 9. Faculty Profiles

### 9.1. List Faculty

- **GET** `/faculty`
- **Auth:** Public
- **Query Params:**
  - `department_id` (uuid): Filter by department
  - `designation` (string): Filter by designation
  - `search` (string): Search by name
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "faculty_id": "uuid",
      "user_id": "uuid",
      "employee_id": "EMP001",
      "name": "Dr. Jane Doe",
      "email": "jane@university.edu",
      "department": {
        "department_id": "uuid",
        "department_name": "Computer Science"
      },
      "designation": "Associate Professor",
      "specialization": "Machine Learning, AI",
      "is_active": true
    }
  ],
  "pagination": { ... }
}
```

---

### 9.2. Get Faculty by ID

- **GET** `/faculty/{faculty_id}`
- **Auth:** Public

**Response:** `200 OK`

```json
{
  "faculty_id": "uuid",
  "user_id": "uuid",
  "employee_id": "EMP001",
  "name": "Dr. Jane Doe",
  "email": "jane@university.edu",
  "phone": "+1234567890",
  "department": {
    "department_id": "uuid",
    "department_name": "Computer Science and Engineering"
  },
  "designation": "Associate Professor",
  "qualification": "Ph.D. Computer Science, MIT",
  "specialization": "Machine Learning, Artificial Intelligence",
  "joining_date": "2015-08-01",
  "office_room": "CSE-301",
  "office_hours": "Mon, Wed 2-4 PM",
  "current_courses": [
    {
      "course_id": "uuid",
      "course_code": "CS201-F2024",
      "course_name": "Data Structures",
      "role": "instructor"
    }
  ],
  "is_active": true
}
```

---

### 9.3. Create Faculty Profile

- **POST** `/faculty`
- **Auth:** Admin only

**Request:**

```json
{
  "user_id": "uuid",
  "employee_id": "EMP001",
  "department_id": "uuid",
  "designation": "Assistant Professor",
  "qualification": "Ph.D. Computer Science",
  "specialization": "Data Mining",
  "joining_date": "2024-01-01",
  "office_room": "CSE-205"
}
```

**Response:** `201 Created`

---

### 9.4. Update Faculty Profile

- **PUT** `/faculty/{faculty_id}`
- **Auth:** Faculty (self), Admin

**Response:** `200 OK`

---

### 9.5. Get Faculty Courses

- **GET** `/faculty/{faculty_id}/courses`
- **Auth:** Public

**Query Params:**
- `semester_id` (uuid): Filter by semester
- `role` (string): instructor, co-instructor, teaching_assistant
- `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "course": {
        "course_id": "uuid",
        "course_code": "CS201-F2024",
        "course_name": "Data Structures and Algorithms",
        "current_enrollment": 45
      },
      "role": "instructor",
      "is_primary": true,
      "assigned_at": "2024-07-15T00:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

## 10. Student Profiles

### 10.1. List Students

- **GET** `/students`
- **Auth:** Faculty, Admin
- **Query Params:**
  - `department_id` (uuid): Filter by department
  - `program_id` (uuid): Filter by program
  - `current_semester` (int): Filter by semester
  - `batch_year` (int): Filter by batch
  - `search` (string): Search by name or registration number
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "student_id": "uuid",
      "user_id": "uuid",
      "registration_number": "2022CS001",
      "name": "John Student",
      "email": "john@university.edu",
      "department": {
        "department_id": "uuid",
        "department_name": "Computer Science"
      },
      "program": {
        "program_id": "uuid",
        "program_code": "BTCS"
      },
      "current_semester": 5,
      "batch_year": 2022,
      "current_cgpa": 8.5,
      "is_active": true
    }
  ],
  "pagination": { ... }
}
```

---

### 10.2. Get Student by ID

- **GET** `/students/{student_id}`
- **Auth:** Student (self), Faculty (if teaching), Admin

**Response:** `200 OK`

```json
{
  "student_id": "uuid",
  "user_id": "uuid",
  "registration_number": "2022CS001",
  "roll_number": "22CS01",
  "name": "John Student",
  "email": "john@university.edu",
  "phone": "+1234567890",
  "department": {
    "department_id": "uuid",
    "department_name": "Computer Science and Engineering"
  },
  "program": {
    "program_id": "uuid",
    "program_name": "B.Tech Computer Science",
    "duration_years": 4
  },
  "current_semester": 5,
  "batch_year": 2022,
  "admission_date": "2022-08-01",
  "current_cgpa": 8.5,
  "total_credits_earned": 80,
  "current_enrollments": [
    {
      "course_id": "uuid",
      "course_code": "CS301-F2024",
      "course_name": "Operating Systems",
      "credits": 4
    }
  ],
  "is_active": true
}
```

---

### 10.3. Create Student Profile

- **POST** `/students`
- **Auth:** Admin only

**Request:**

```json
{
  "user_id": "uuid",
  "registration_number": "2024CS001",
  "roll_number": "24CS01",
  "department_id": "uuid",
  "program_id": "uuid",
  "current_semester": 1,
  "batch_year": 2024,
  "admission_date": "2024-08-01"
}
```

**Response:** `201 Created`

---

### 10.4. Update Student Profile

- **PUT** `/students/{student_id}`
- **Auth:** Student (limited fields), Admin (all fields)

**Response:** `200 OK`

---

### 10.5. Update Student Semester

Promote student to next semester.

- **POST** `/students/{student_id}/promote`
- **Auth:** Admin only

**Request:**

```json
{
  "new_semester": 6,
  "update_cgpa": 8.7,
  "credits_earned": 18
}
```

**Response:** `200 OK`

---

## 11. Academic Calendar

### 11.1. List Calendar Events

- **GET** `/calendar`
- **Auth:** Public
- **Query Params:**
  - `semester_id` (uuid): Filter by semester
  - `event_type` (string): holiday, exam, registration, deadline
  - `start_date` (date): Events starting from
  - `end_date` (date): Events until
  - `is_holiday` (boolean): Only holidays
  - `page`, `limit`: Pagination

**Response:** `200 OK`

```json
{
  "data": [
    {
      "event_id": "uuid",
      "event_name": "Mid-Semester Examination",
      "event_type": "exam",
      "start_date": "2024-10-01",
      "end_date": "2024-10-10",
      "description": "Mid-semester examinations for all courses",
      "is_holiday": false,
      "semester": {
        "semester_id": "uuid",
        "semester_name": "Fall 2024"
      }
    }
  ],
  "pagination": { ... }
}
```

---

### 11.2. Create Calendar Event

- **POST** `/calendar`
- **Auth:** Admin only

**Request:**

```json
{
  "semester_id": "uuid",
  "event_name": "Diwali Holiday",
  "event_type": "holiday",
  "start_date": "2024-11-01",
  "end_date": "2024-11-05",
  "description": "Diwali festival holidays",
  "is_holiday": true
}
```

**Response:** `201 Created`

---

### 11.3. Update Calendar Event

- **PUT** `/calendar/{event_id}`
- **Auth:** Admin only

**Response:** `200 OK`

---

### 11.4. Delete Calendar Event

- **DELETE** `/calendar/{event_id}`
- **Auth:** Admin only

**Response:** `204 No Content`

---

## 12. Kafka Events

The Course Service integrates with Apache Kafka for event-driven communication.

### Topics Published

| Topic | Events |
|-------|--------|
| `course.events` | COURSE_CREATED, COURSE_UPDATED, COURSE_DELETED, COURSE_ACTIVATED, COURSE_DEACTIVATED, FACULTY_ASSIGNED, FACULTY_UNASSIGNED |
| `enrollment.events` | STUDENT_ENROLLED, STUDENT_DROPPED, ENROLLMENT_COMPLETED, ENROLLMENT_FAILED, WAITLIST_ADDED, WAITLIST_PROMOTED |

### Topics Consumed

| Topic | Consumer Group | Events Handled |
|-------|----------------|----------------|
| `user.events` | `course-service-group` | USER_CREATED, USER_UPDATED, USER_DELETED |

### Event Schema Example

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "COURSE_CREATED",
  "event_version": "1.0",
  "timestamp": "2024-12-27T10:30:00Z",
  "service_name": "course-service",
  "correlation_id": "req-12345-abcde",
  "metadata": {
    "user_id": "admin-uuid",
    "ip_address": "192.168.1.100"
  },
  "payload": {
    "course_id": "uuid",
    "course_code": "CS201-F2024",
    "course_name": "Data Structures and Algorithms",
    "subject_id": "uuid",
    "department_id": "uuid",
    "program_id": "uuid",
    "semester": 3,
    "academic_year": 2024,
    "max_students": 60,
    "created_by": "admin-uuid"
  }
}
```

---

## 13. Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }
  }
}
```

### Common Error Codes

| HTTP Status | Code | Description |
|-------------|------|-------------|
| 400 | `VALIDATION_ERROR` | Invalid request data |
| 400 | `PREREQUISITES_NOT_MET` | Course prerequisites not satisfied |
| 400 | `COURSE_FULL` | Course has reached maximum enrollment |
| 400 | `REGISTRATION_CLOSED` | Course registration period ended |
| 401 | `UNAUTHORIZED` | Authentication required |
| 403 | `FORBIDDEN` | Insufficient permissions |
| 404 | `NOT_FOUND` | Resource not found |
| 409 | `DUPLICATE_ENTRY` | Resource already exists |
| 409 | `ALREADY_ENROLLED` | Student already enrolled in course |
| 422 | `UNPROCESSABLE_ENTITY` | Business rule violation |
| 500 | `INTERNAL_ERROR` | Server error |

### Error Examples

**Validation Error:**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "fields": [
        { "field": "course_code", "message": "Course code is required" },
        { "field": "semester_number", "message": "Must be between 1 and 8" }
      ]
    }
  }
}
```

**Prerequisites Not Met:**

```json
{
  "error": {
    "code": "PREREQUISITES_NOT_MET",
    "message": "Cannot enroll: prerequisites not satisfied",
    "details": {
      "missing_prerequisites": [
        { "subject_code": "CS101", "subject_name": "Introduction to Programming" }
      ]
    }
  }
}
```

---

## 14. Permissions (RBAC)

### Role Permissions Matrix

| Resource | Action | Student | Faculty | Admin |
|----------|--------|:-------:|:-------:|:-----:|
| **Departments** | Read | Yes | Yes | Yes |
| | Create/Update/Delete | - | - | Yes |
| **Programs** | Read | Yes | Yes | Yes |
| | Create/Update/Delete | - | - | Yes |
| **Subjects** | Read | Yes | Yes | Yes |
| | Create/Update | - | Own Dept | Yes |
| | Delete | - | - | Yes |
| **Semesters** | Read | Yes | Yes | Yes |
| | Create/Update/Delete | - | - | Yes |
| **Courses** | Read (Basic) | Yes | Yes | Yes |
| | Read (Full) | Enrolled | Assigned | Yes |
| | Create | - | Own Dept | Yes |
| | Update | - | Assigned | Yes |
| | Delete | - | - | Yes |
| | Activate/Deactivate | - | Primary | Yes |
| **Faculty Assignment** | Read | Yes | Yes | Yes |
| | Create/Update/Delete | - | - | Yes |
| **Enrollments** | Enroll Self | Yes | - | Yes |
| | Enroll Others | - | - | Yes |
| | Drop Self | Yes | - | Yes |
| | Drop Others | - | - | Yes |
| | Update Grade | - | Assigned | Yes |
| | View Own | Yes | - | - |
| | View Course Students | - | Assigned | Yes |
| **Faculty Profiles** | Read | Yes | Yes | Yes |
| | Create | - | - | Yes |
| | Update | - | Self | Yes |
| **Student Profiles** | Read Own | Yes | - | - |
| | Read Others | - | Teaching | Yes |
| | Create | - | - | Yes |
| | Update | Limited | - | Yes |
| **Academic Calendar** | Read | Yes | Yes | Yes |
| | Create/Update/Delete | - | - | Yes |

### Permission Definitions

```
course:create    - Create new courses
course:read      - View course details
course:update    - Update course metadata
course:delete    - Delete/archive courses
course:activate  - Activate a course

enrollment:create     - Enroll students
enrollment:read       - View enrollments
enrollment:update     - Update enrollment (grades)
enrollment:delete     - Drop enrollments

faculty:create   - Create faculty profiles
faculty:read     - View faculty details
faculty:update   - Update faculty profiles
faculty:assign   - Assign faculty to courses

student:create   - Create student profiles
student:read     - View student details
student:update   - Update student profiles
student:promote  - Promote to next semester
```

---

## 15. Rate Limiting

| Endpoint Category | Rate Limit |
|-------------------|------------|
| Public Read | 100 requests/minute |
| Authenticated Read | 200 requests/minute |
| Write Operations | 50 requests/minute |
| Bulk Operations | 10 requests/minute |

---

## 16. Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024-12-27 | Initial API specification |
| 2.0 | 2024-12-28 | Complete rewrite aligned with microservice architecture |
