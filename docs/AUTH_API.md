# Authentication API Documentation

## Overview

The Mizon authentication service provides secure user authentication with enterprise-grade security features including password hashing, rate limiting, and secure token management.

**Base URL:** `http://localhost:8001`

## Security Features

- ✅ Bcrypt password hashing (cost factor: 12)
- ✅ Strong password requirements (min 8 chars, uppercase, lowercase, number, special char)
- ✅ Rate limiting on sensitive endpoints
- ✅ Input sanitization to prevent XSS
- ✅ Secure random token generation for password resets
- ✅ Email enumeration protection
- ✅ Session management with configurable expiry
- ✅ Remember me functionality

## Endpoints

### 1. User Registration

**POST** `/api/auth/signup`

Create a new user account with validated credentials.

**Request Body:**
```json
{
  "username": "string (required, min 3 chars, alphanumeric + _ -)",
  "email": "string (required, valid email format)",
  "password": "string (required, min 8 chars with complexity requirements)",
  "first_name": "string (optional)",
  "last_name": "string (optional)"
}
```

**Password Requirements:**
- Minimum 8 characters
- Maximum 128 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character (!@#$%^&*(),.?":{}|<>)

**Response (201 Created):**
```json
{
  "user_id": "uuid",
  "username": "string",
  "email": "string",
  "message": "User created successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or password requirements not met
- `409 Conflict` - Username or email already exists
- `429 Too Many Requests` - Rate limit exceeded (max 3 signups per hour per IP)

---

### 2. User Login

**POST** `/api/auth/login`

Authenticate a user and create a session.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)",
  "remember_me": "boolean (optional, default: false)"
}
```

**Response (200 OK):**
```json
{
  "token": "uuid",
  "user_id": "uuid",
  "username": "string"
}
```

**Session Expiry:**
- Normal login: 24 hours
- Remember me: 30 days

**Error Responses:**
- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Invalid credentials
- `429 Too Many Requests` - Rate limit exceeded (max 5 attempts per minute per IP)

---

### 3. Validate Token

**GET** `/api/auth/validate`

Validate an authentication token and check if it's still active.

**Headers:**
```
Authorization: <token>
```

**Response (200 OK):**
```json
{
  "valid": true,
  "user_id": "uuid"
}
```

or

```json
{
  "valid": false
}
```

---

### 4. User Logout

**POST** `/api/auth/logout`

Invalidate the current session token.

**Headers:**
```
Authorization: <token>
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

**Error Responses:**
- `401 Unauthorized` - No or invalid token provided

---

### 5. Get User Profile

**GET** `/api/auth/profile`

Retrieve the authenticated user's profile information.

**Headers:**
```
Authorization: <token>
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "created_at": "timestamp"
}
```

**Error Responses:**
- `401 Unauthorized` - No or invalid token provided

---

### 6. Update User Profile

**PUT** `/api/auth/profile`

Update the authenticated user's profile information.

**Headers:**
```
Authorization: <token>
```

**Request Body:**
```json
{
  "email": "string (optional)",
  "first_name": "string (optional)",
  "last_name": "string (optional)"
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email format
- `401 Unauthorized` - No or invalid token provided
- `409 Conflict` - Email already in use by another user

---

### 7. Forgot Password

**POST** `/api/auth/forgot-password`

Initiate a password reset process by generating a secure reset token.

**Request Body:**
```json
{
  "email": "string (required)"
}
```

**Response (200 OK):**
```json
{
  "message": "If the email exists, a password reset link has been sent"
}
```

**Note:** This endpoint always returns 200 OK to prevent email enumeration attacks. The reset token is logged server-side for development purposes. In production, it should be sent via email.

**Token Expiry:** 1 hour

---

### 8. Reset Password

**POST** `/api/auth/reset-password`

Complete the password reset process using a valid reset token.

**Request Body:**
```json
{
  "token": "string (required)",
  "password": "string (required, must meet password requirements)"
}
```

**Response (200 OK):**
```json
{
  "message": "Password reset successful"
}
```

**Side Effects:**
- All existing user sessions are invalidated
- The reset token is deleted

**Error Responses:**
- `400 Bad Request` - Invalid token, expired token, or password requirements not met

---

## Rate Limiting

The following rate limits are enforced:

| Endpoint | Limit | Window |
|----------|-------|--------|
| `/api/auth/signup` | 3 requests | 1 hour |
| `/api/auth/login` | 5 requests | 1 minute |

Rate limits are applied per IP address.

---

## Example Usage

### Register a new user
```bash
curl -X POST http://localhost:8001/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!",
    "remember_me": true
  }'
```
