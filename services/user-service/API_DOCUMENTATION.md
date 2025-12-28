# NimbusU User Service API Documentation

## Overview

The User Service provides authentication, user management, and authorization for the NimbusU platform.

**Base URL:** `http://localhost:8081` (Development)
**Swagger UI:** `http://localhost:8081/swagger/index.html`

## Authentication

Authentication is handled via JWT tokens (Access and Refresh tokens).

### Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/login` | Login with email and password | No |
| `POST` | `/auth/refresh` | Refresh access token | No |
| `POST` | `/auth/password/reset-request` | Request password reset email | No |
| `POST` | `/auth/password/reset` | Reset password with token | No |
| `POST` | `/auth/logout` | Logout (revoke session) | Yes |
| `POST` | `/auth/password/change` | Change password | Yes |
| `GET` | `/auth/sessions` | List active sessions | Yes |
| `DELETE` | `/auth/sessions` | Revoke all sessions | Yes |
| `DELETE` | `/auth/sessions/{sessionId}` | Revoke specific session | Yes |

## User Management (Self-Service)

Endpoints for users to manage their own profile.

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/users/me` | Get current user profile | Yes |
| `PUT` | `/users/me` | Update current user profile | Yes |

## Admin Management

Endpoints for administrators to manage users. Requires `admin` or `faculty` role.

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/admin/users` | Create new user | Yes (Admin) |
| `GET` | `/admin/users` | List users (paginated) | Yes (Admin) |
| `GET` | `/admin/users/{id}` | Get user by ID | Yes (Admin) |
| `PUT` | `/admin/users/{id}` | Update user details | Yes (Admin) |
| `DELETE` | `/admin/users/{id}` | Delete user | Yes (Admin) |
| `POST` | `/admin/users/{id}/activate` | Activate user | Yes (Admin) |
| `POST` | `/admin/users/{id}/suspend` | Suspend user | Yes (Admin) |
| `POST` | `/admin/users/bulk-import` | Bulk import users | Yes (Admin) |

## Data Models

### Login Request
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Login Response
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": { ... }
  }
}
```

## Error Handling

Standard error response format:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Running Swagger Locally

1. Start the service: `make run`
2. Open `http://localhost:8081/swagger/index.html` in your browser.
