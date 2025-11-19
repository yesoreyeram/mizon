# User Simulator Service

A background service that automatically generates and registers random users in the Mizon authentication system at configurable intervals.

## Features

- Generates realistic user data with random names
- Creates unique usernames and email addresses
- Generates secure passwords meeting all authentication requirements
- Configurable generation rate (users per minute)
- Spreads requests evenly across time to avoid rate limiting
- Automatic retry and error handling
- Rate limit detection and warnings

## Configuration

The service is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_SERVICE_URL` | `http://localhost:8001` | URL of the authentication service |
| `MIN_USERS_PER_MINUTE` | `2` | Minimum number of users to generate per minute |
| `MAX_USERS_PER_MINUTE` | `7` | Maximum number of users to generate per minute |
| `LOG_LEVEL` | `INFO` | Logging level (DEBUG, INFO, WARN, ERROR) |
| `LOG_FORMAT` | `json` | Log format (json or text) |

## Usage

### With Docker Compose

The service is included in the main docker-compose.yml and will start automatically:

```bash
docker-compose up -d user-simulator
```

### Standalone

```bash
# Set environment variables
export AUTH_SERVICE_URL=http://localhost:8001
export MIN_USERS_PER_MINUTE=2
export MAX_USERS_PER_MINUTE=7

# Run the service
go run main.go
```

## How It Works

1. **User Generation**: Every minute, the service generates a random number of users between `MIN_USERS_PER_MINUTE` and `MAX_USERS_PER_MINUTE`

2. **User Data**: Each user is created with:
   - Unique username (firstname + lastname + timestamp)
   - Valid email address (firstname.lastname@example.com)
   - Secure password (12+ characters with uppercase, lowercase, numbers, and special characters)
   - Random first and last name from predefined lists

3. **Request Spacing**: Requests are evenly distributed across the minute to avoid bursts and rate limiting

4. **Error Handling**:
   - Rate limit errors (429) are tracked separately
   - Duplicate user errors (409) are logged but don't stop processing
   - Network errors are logged and retried in the next batch

## Generated User Example

```json
{
  "username": "EmmaSmit42315",
  "email": "Emma.Smith4231@example.com",
  "password": "aB3!xY9@mN2$",
  "first_name": "Emma",
  "last_name": "Smith"
}
```

## Rate Limiting

The auth service has rate limits:
- **Signup**: 3 requests per hour per IP

If you set `MAX_USERS_PER_MINUTE` too high, you may hit these limits. The service will detect and warn about rate limiting. In production, consider:
- Using a distributed rate limiter (Redis)
- Adjusting the generation rate
- Implementing IP rotation if needed

## Monitoring

The service logs:
- Number of users generated per batch
- Success/failure counts
- Rate limit warnings
- Individual user creation confirmations

Example log output:
```
INFO User simulator started
INFO Target: 2-7 users per minute
INFO Auth service URL: http://auth-service:8001
INFO Generating 5 users...
INFO Created user: EmmaSmit42315 (Emma.Smith4231@example.com)
INFO Created user: LiamJohn98765 (Liam.Johnson9876@example.com)
...
INFO Batch complete: 5 succeeded, 0 failed
```

## Building

```bash
go build -o user-simulator
```

Or with Docker:

```bash
docker build -t mizon-user-simulator -f Dockerfile ../..
```

## Testing

The service can be tested by:

1. Running it locally with low generation rates
2. Checking the auth service logs for signup requests
3. Querying the database for created users:

```sql
SELECT username, email, created_at 
FROM users 
WHERE username LIKE '%[0-9]%' 
ORDER BY created_at DESC 
LIMIT 10;
```

## Security Considerations

- Generated users have secure passwords but use predictable email domains
- In production, consider adding email verification
- The service should run in a trusted environment
- Generated users are for testing/demo purposes only
