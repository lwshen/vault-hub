# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Go Backend

- **Build**: `go build -o tmp/main ./cmd/main.go`
- **Run**: `go run ./cmd/main.go`
- **Test**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` (run all tests with required env vars)
- **Test specific package**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./model -v`
- **Generate API code**: `go generate api/tool.go` (run after modifying files in `api/openapi/*`)

### React Frontend (web/)

- **Install dependencies**: `pnpm install` (uses pnpm as package manager)
- **Development server**: `pnpm run dev`
- **Build production**: `pnpm run build`
- **Lint**: `pnpm run lint`
- **Type check**: `pnpm run typecheck`
- **Preview build**: `pnpm run preview`

## Architecture Overview

VaultHub is a secure environment variable and API key management system with AES-256-GCM encryption.

### Backend (Go + Fiber)

- **Entry point**: `cmd/main.go` - Sets up Fiber web server
- **Database**: GORM with support for SQLite, MySQL, PostgreSQL
- **API**: OpenAPI 3.0 spec in `api/openapi/api.yaml`, generated code in `api/generated.go`
- **Models**: `model/` - Database entities (User, Vault, AuditLog, APIKey)
- **Routes**: `route/` - HTTP routing and middleware
- **Handlers**: `handler/` - Request/response handling
- **Internal packages**:
  - `internal/config/` - Environment configuration
  - `internal/auth/` - JWT and OIDC authentication
  - `internal/encryption/` - AES-256-GCM encryption for vault values

### Frontend (React + TypeScript + Vite)

- **Framework**: React 19 with TypeScript
- **Build tool**: Vite
- **Routing**: Wouter (lightweight router)
- **UI**: Tailwind CSS + Radix UI components
- **API client**: Custom generated TypeScript client (`@lwshen/vault-hub-ts-fetch-client`)
- **State**: React Context for auth and theme management

### Key Security Features

- All vault values encrypted with AES-256-GCM before database storage
- JWT-based authentication with optional OIDC support
- API key authentication for programmatic access
- Transparent encryption/decryption at model layer
- Audit logging for all vault operations
- Strict authentication middleware with route-based credential enforcement

## Required Environment Variables

For the backend to start, you must set:

- `JWT_SECRET` - Secret for JWT token signing
- `ENCRYPTION_KEY` - AES-256 encryption key (generate with `openssl rand -base64 32`)

Optional configuration:

- `APP_PORT` (default: 3000)
- `DATABASE_TYPE` (sqlite|mysql|postgres, default: sqlite)
- `DATABASE_URL` (default: data.db)
- OIDC settings: `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_ISSUER`

## Database Models

- **User**: User accounts with email/password or OIDC
- **Vault**: Encrypted key-value pairs for environment variables
- **AuditLog**: Audit trail of vault operations
- **APIKey**: API key management for programmatic access

## API Generation

The project uses OpenAPI 3.0 specification (`api/openapi/api.yaml`) with `oapi-codegen` to generate:

- Go server stubs (`api/generated.go`)
- TypeScript client library (published as npm package)

**Important**: Always run `go generate api/tool.go` after modifying files in `api/openapi/*` to regenerate the Go types and interfaces. The API spec uses camelCase naming convention for all properties (e.g., `uniqueId`, `createdAt`, `isActive`).

**NEVER EDIT**: Do not modify `api/generated.go` directly as it is auto-generated code. All API changes must be made in the OpenAPI specification files in `api/openapi/*`.

## Authentication & Authorization

### Authentication Middleware Rules

The application enforces strict authentication rules via middleware (`route/middleware.go`):

**Public Routes (No Authentication Required):**
- `/api/auth/login` - User login
- `/api/auth/register` - User registration  
- `/api/auth/login/oidc` - OIDC login
- `/api/auth/callback/oidc` - OIDC callback
- Static web assets (`/`, `/*`)

**API Key Only Routes:**
- `/api/api-key/*` - Vault access via API keys (e.g., `/api/api-key/vaults`, `/api/api-key/vault/{id}`)
- Must use `Authorization: Bearer vhub_xxx` header
- Rejects JWT tokens with error message

**JWT Only Routes:**
- All other `/api/*` routes - User management, API key management, vault management via web UI
- Must use `Authorization: Bearer <jwt_token>` header  
- Rejects API keys with error message

### Context Variables

- **API Key Auth**: Sets `c.Locals("user_id", &key.UserID)` and `c.Locals("api_key", key)`
- **JWT Auth**: Sets `c.Locals("user", &user)` (full User object)

### API Endpoints

**API Key Vault Access:**
- `GET /api/api-key/vaults` - List accessible vaults (VaultLite format, no decrypted values)
- `GET /api/api-key/vault/{uniqueId}` - Get specific vault (full Vault format with decrypted value)
- Implements proper access control via `APIKey.HasVaultAccess()`
- Includes audit logging for vault read operations
- **Enhanced Security**: Supports optional client-side encryption via `X-Enable-Client-Encryption: true` header
  - Uses PBKDF2 key derivation from API key + vault unique ID as salt
  - Provides per-vault encryption keys without key exchange complexity

## Testing Strategy

- Go unit tests for encryption (`internal/encryption/encryption_test.go`)
- Database model tests (`model/db_test.go`)
- Configuration tests (`internal/config/config_test.go`)
- Frontend uses standard React testing patterns

## Frontend Code Style

ESLint configuration enforces:

- 2-space indentation
- Single quotes
- Semicolons required
- Stylistic rules from `@stylistic/eslint-plugin`
- React-specific rules and hooks validation
- TypeScript strict mode
