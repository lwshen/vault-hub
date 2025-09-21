# API Reference

VaultHub provides a RESTful API with OpenAPI 3.0 specification. All API endpoints use JSON for data exchange and require proper authentication.

## Authentication

VaultHub supports two authentication methods depending on the endpoint:

### JWT Authentication
Used for web interface and user management endpoints.

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### API Key Authentication
Used for CLI and programmatic access to vault data.

```http
GET /api/cli/vaults
Authorization: Bearer vhub_your_api_key_here
```

## Vault Operations

### List Vaults (CLI)

```http
GET /api/cli/vaults
Authorization: Bearer vhub_your_api_key_here
```

**Response:**
```json
{
  "vaults": [
    {
      "uniqueId": "vault-uuid",
      "name": "production-secrets",
      "description": "Production environment variables",
      "createdAt": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### Get Vault by Name (CLI)

```http
GET /api/cli/vault/name/{name}
Authorization: Bearer vhub_your_api_key_here
```

**Response:**
```json
{
  "uniqueId": "vault-uuid",
  "name": "production-secrets",
  "description": "Production environment variables",
  "value": {
    "API_KEY": "secret-api-key",
    "DATABASE_URL": "postgresql://...",
    "REDIS_URL": "redis://..."
  },
  "createdAt": "2025-01-01T00:00:00Z"
}
```

### Get Vault by ID (CLI)

```http
GET /api/cli/vault/{uniqueId}
Authorization: Bearer vhub_your_api_key_here
```

### Create Vault (Web)

```http
POST /api/vaults
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "staging-secrets",
  "description": "Staging environment variables",
  "value": {
    "API_KEY": "staging-api-key",
    "DATABASE_URL": "postgresql://staging..."
  }
}
```

### Update Vault (Web)

```http
PUT /api/vaults/{id}
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "updated-name",
  "description": "Updated description",
  "value": {
    "API_KEY": "new-api-key",
    "DATABASE_URL": "postgresql://new..."
  }
}
```

### Delete Vault (Web)

```http
DELETE /api/vaults/{id}
Authorization: Bearer jwt_token_here
```

## API Key Management

### Create API Key

```http
POST /api/api-keys
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "CI/CD Pipeline Key",
  "description": "Key for accessing production secrets in CI"
}
```

**Response:**
```json
{
  "id": "key-uuid",
  "name": "CI/CD Pipeline Key",
  "key": "vhub_generated_api_key_here",
  "createdAt": "2025-01-01T00:00:00Z"
}
```

### List API Keys

```http
GET /api/api-keys
Authorization: Bearer jwt_token_here
```

**Response:**
```json
{
  "apiKeys": [
    {
      "id": "key-uuid",
      "name": "CI/CD Pipeline Key",
      "createdAt": "2025-01-01T00:00:00Z",
      "lastUsed": "2025-01-02T12:00:00Z"
    }
  ]
}
```

### Delete API Key

```http
DELETE /api/api-keys/{id}
Authorization: Bearer jwt_token_here
```

## User Management

### Register User

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

### Get Current User

```http
GET /api/auth/me
Authorization: Bearer jwt_token_here
```

### Update User Profile

```http
PUT /api/auth/me
Authorization: Bearer jwt_token_here
Content-Type: application/json

{
  "name": "Updated Name"
}
```

## System Status

### Get System Status

```http
GET /api/status
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.2.9",
  "database": {
    "status": "healthy",
    "responseTime": "2ms"
  },
  "system": {
    "memoryUsage": 45.2,
    "diskUsage": 23.1
  }
}
```

## Error Responses

All API endpoints return consistent error responses:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

### Common Error Codes

| Status Code | Error Code | Description |
|-------------|------------|-------------|
| 400 | `INVALID_REQUEST` | Invalid request format or parameters |
| 401 | `UNAUTHORIZED` | Missing or invalid authentication |
| 403 | `FORBIDDEN` | Insufficient permissions |
| 404 | `NOT_FOUND` | Resource not found |
| 409 | `CONFLICT` | Resource already exists |
| 500 | `INTERNAL_ERROR` | Internal server error |

## Rate Limiting

API requests are rate limited to prevent abuse:

- **Web endpoints**: 100 requests per minute per user
- **CLI endpoints**: 1000 requests per minute per API key

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1641024000
```