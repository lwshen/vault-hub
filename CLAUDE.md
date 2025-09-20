# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Go Backend (apps/server/)

- **Build**: `go build -o tmp/main ./apps/server/main.go`
- **Build with version**: `go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') -X github.com/lwshen/vault-hub/internal/version.Commit=$(git rev-parse --short HEAD)" -o tmp/main ./apps/server/main.go`
- **Run**: `go run ./apps/server/main.go`
- **Test**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` (run all tests with required env vars)
- **Test specific package**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./model -v`
- **Generate API code**: `go generate packages/api/tool.go` (run after modifying files in `packages/api/openapi/*`)

### Go CLI (apps/cli/)

- **Build**: `go build -o tmp/vault-hub-cli ./apps/cli/main.go`
- **Build with version**: `go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') -X github.com/lwshen/vault-hub/internal/version.Commit=$(git rev-parse --short HEAD)" -o tmp/vault-hub-cli ./apps/cli/main.go`
- **Run**: `go run ./apps/cli/main.go`
- **Commands**:
  - `vault-hub list` or `vault-hub ls` - List all accessible vaults
  - `vault-hub get --name/--id <vault-name-or-id>` - Get a specific vault by name or unique ID
    - `--exec` flag: Execute command if vault has been updated since last output
    - Example: `vault-hub get --name my-secrets --output .env --exec "source .env && npm start"`
  - `vault-hub get --name <vault-name>` - Get vault by name using `/api/cli/vault/name/{name}` endpoint
  - `vault-hub version` - Show version and commit information
- **Multi-platform builds**: See CI configuration for cross-compilation examples

### React Frontend (apps/web/)

- **Install dependencies**: `pnpm install` (uses pnpm as package manager)
- **Development server**: `pnpm run dev`
- **Build production**: `pnpm run build`
- **Lint**: `pnpm run lint`
- **Type check**: `pnpm run typecheck`
- **Preview build**: `pnpm run preview`

## Architecture Overview

VaultHub is a comprehensive secure environment variable and API key management system with AES-256-GCM encryption, consisting of three main components:

### Backend (Go + Fiber)

- **Entry point**: `apps/server/main.go` - Sets up Fiber web server
- **Database**: GORM with support for SQLite, MySQL, PostgreSQL
- **API**: OpenAPI 3.0 spec in `packages/api/openapi/api.yaml`, generated code in `packages/api/generated.go`
- **Models**: `model/` - Database entities (User, Vault, AuditLog, APIKey)
- **Routes**: `route/` - HTTP routing and middleware
- **Handlers**: `handler/` - Request/response handling
- **Internal packages**:
  - `internal/config/` - Environment configuration
  - `internal/auth/` - JWT and OIDC authentication
  - `internal/encryption/` - AES-256-GCM encryption for vault values

### Frontend (React + TypeScript + Vite)

- **Location**: `apps/web/`
- **Framework**: React 19 with TypeScript
- **Build tool**: Vite 7.1.3 with Tailwind CSS 4.1.12 (Lightning CSS)
- **Package manager**: pnpm 10.15.0
- **Routing**: Wouter (lightweight router)
- **UI**: Tailwind CSS 4.x + Radix UI components + Framer Motion for animations
- **API client**: Custom generated TypeScript client (`@lwshen/vault-hub-ts-fetch-client`)
- **State**: Zustand stores for component state, React Context for auth and theme management
- **Components**: Organized into dashboard, layout, modals, and UI components
- **Development proxy**: API requests proxied to `http://localhost:3000`
- **Build optimization**: Manual chunking for UI libraries, vendor packages, and API client

### CLI (Go + Cobra)

- **Location**: `apps/cli/`
- **Framework**: Cobra for command-line interface
- **Entry point**: `apps/cli/main.go` - Sets up Cobra CLI with vault management commands
- **Commands**:
  - `list` (alias: `ls`) - List all accessible vaults
  - `get --name/--id <vault-name-or-id>` - Get specific vault by name or unique ID
- **API Integration**: Designed to work with `/api/cli/*` endpoints for API key authentication
- **Cross-platform**: Built for Linux, Windows, and macOS (amd64, arm64)

### Key Security Features

- All vault values encrypted with AES-256-GCM before database storage
- JWT-based authentication with optional OIDC support
- API key authentication for programmatic access
- Transparent encryption/decryption at model layer
- Audit logging for all vault operations
- Strict authentication middleware with route-based credential enforcement

### Health Monitoring

The `/api/status` endpoint provides comprehensive system monitoring:

- **Database Health**: Response time, connection pool status, availability checks
- **System Health**: Memory usage, disk space, overall system status  
- **Status Levels**: `healthy`, `degraded`, `unavailable` with specific thresholds
- **Performance Metrics**: Database response times, connection counts, resource utilization
- **Multi-factor Assessment**: System status determined by database health, memory usage, disk space

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

The project uses OpenAPI 3.0 specification (`packages/api/openapi/api.yaml`) with `oapi-codegen` to generate:

- Go server stubs (`packages/api/generated.go`)
- TypeScript client library (published as npm package)

**Important**: Always run `go generate packages/api/tool.go` after modifying files in `packages/api/openapi/*` to regenerate the Go types and interfaces. The API spec uses camelCase naming convention for all properties (e.g., `uniqueId`, `createdAt`, `isActive`).

**NEVER EDIT**: Do not modify `packages/api/generated.go` directly as it is auto-generated code. All API changes must be made in the OpenAPI specification files in `packages/api/openapi/*`.

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
- `/api/cli/*` - Vault access via API keys (e.g., `/api/cli/vaults`, `/api/cli/vault/{id}`)
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

**Public API:**
- `GET /api/status` - Get comprehensive system status including version, health, and performance metrics (no authentication required)

**CLI API Vault Access:**
- `GET /api/cli/vaults` - List accessible vaults (VaultLite format, no decrypted values)
- `GET /api/cli/vault/{uniqueId}` - Get specific vault (full Vault format with decrypted value)
- `GET /api/cli/vault/name/{name}` - Get specific vault by name (full Vault format with decrypted value)
- Implements proper access control via `APIKey.HasVaultAccess()`
- Includes audit logging for vault read operations
- **Enhanced Security**: Supports optional client-side encryption via `X-Enable-Client-Encryption: true` header
  - Uses PBKDF2 key derivation from API key + vault unique ID as salt
  - Provides per-vault encryption keys without key exchange complexity

## Go Code Quality

**IMPORTANT**: Always run `golangci-lint run ./...` after editing Go code to ensure code quality and formatting standards are met. This will check for:

- Formatting issues (gofmt)
- Security vulnerabilities (gosec)
- Code style violations
- Unused variables/parameters
- Other Go best practices

**Format Go code**: Use `gofmt -w <files>` to automatically format Go files before committing.

## Testing Strategy

- Go unit tests for encryption (`internal/encryption/encryption_test.go`)
- Database model tests (`model/db_test.go`)
- Configuration tests (`internal/config/config_test.go`)
- Frontend uses standard React testing patterns

## Frontend State Management

The frontend uses Zustand for component-level state management:

- **Zustand stores**: Located in `src/stores/` for audit logs, API keys, and vaults
- **Store pattern**: Each store contains state, actions, and loading states with comprehensive error handling
- **Input validation**: All user inputs (pagination, deletion) include validation and error boundaries
- **API integration**: Stores directly use generated API clients with proper error handling
- **React Context**: Still used for global auth and theme state

## Frontend Code Style

ESLint configuration enforces:

- 2-space indentation
- Single quotes
- Semicolons required
- Stylistic rules from `@stylistic/eslint-plugin`
- React-specific rules and hooks validation
- TypeScript strict mode

## CI/CD Pipeline

### GitHub Actions Workflows

The project uses multiple GitHub Actions workflows for comprehensive CI/CD:

#### Main CI Workflow (`.github/workflows/ci.yml`)
- **Triggers**: Push to main, pull requests to main
- **Go Version**: 1.24.2 with module caching
- **Frontend**: pnpm 10.15.0 with Node.js 22
- **Quality Checks**: golangci-lint, frontend typecheck and lint
- **Testing**: Go tests with required environment variables
- **Builds**: Cross-platform binaries for both server and CLI (Linux/Windows/macOS, amd64/arm64)
- **Artifacts**: Uploads server, CLI binaries, and frontend build

#### Release Workflow (`.github/workflows/release.yml`)
- **Triggers**: Git tags matching `v*`
- **Client Publishing**: 
  - TypeScript fetch client (`@lwshen/vault-hub-ts-fetch-client`) to npm
  - Go client to separate repository (`vault-hub-go-client`)
- **Changelog Generation**: Uses git-cliff with conventional commits
- **Release Assets**: Uploads binaries to GitHub releases
- **Automated PR**: Creates pull request to update CHANGELOG.md

#### Additional Workflows
- **Database Testing**: `db-test.yml` - Database-specific tests
- **Docker Images**: `build-image.yml`, `build-cli-image.yml` - Container builds
- **Client Publishing**: `publish-ts-client.yml`, `publish-go-client.yml` - Standalone client publishing
- **Mirror**: `mirror.yml` - Repository mirroring
- **Claude Integration**: `claude.yml`, `claude-code-review.yml` - AI-powered code reviews

#### Release Management
- **Changelog**: Automated generation using git-cliff with conventional commits
- **Versioning**: Git tags drive version information in binaries
- **Client Libraries**: Auto-published on releases with OpenAPI generators

### Build Outputs

**Server binaries**:
- `vault-hub-server-linux-{amd64,arm64}`
- `vault-hub-server-windows-amd64.exe`
- `vault-hub-server-darwin-{amd64,arm64}`

**CLI binaries**:
- `vault-hub-cli-linux-{amd64,arm64}`
- `vault-hub-cli-windows-amd64.exe`
- `vault-hub-cli-darwin-{amd64,arm64}`

## Project Structure

```
vault-hub/
├── .github/workflows/   # GitHub Actions CI/CD workflows
├── apps/
│   ├── cli/              # Command-line interface (Go + Cobra)
│   │   ├── main.go       # CLI entry point
│   │   └── README.md     # CLI documentation
│   ├── server/           # Backend server (Go + Fiber)
│   │   └── main.go       # Server entry point
│   └── web/              # Frontend application (React + TypeScript)
│       ├── src/          # React source code
│       ├── dist/         # Build output
│       ├── public/       # Static assets
│       ├── package.json  # Frontend dependencies
│       ├── vite.config.ts # Vite configuration
│       └── tsconfig.json # TypeScript configuration
├── packages/
│   └── api/              # OpenAPI specification and generated code
│       ├── openapi/      # OpenAPI 3.0 specification files
│       │   ├── api.yaml  # Main specification
│       │   ├── paths/    # Endpoint definitions
│       │   └── schemas/  # Data model schemas
│       ├── generated.go  # Auto-generated Go server code
│       ├── tool.go       # Code generation tool
│       └── *.go         # API implementation files
├── model/               # Database models (GORM)
├── handler/             # HTTP request handlers
├── route/               # Routing and middleware
├── internal/            # Internal packages
│   ├── auth/           # Authentication (JWT, OIDC)
│   ├── cli/            # CLI command implementations
│   ├── config/         # Configuration management
│   ├── encryption/     # AES-256-GCM encryption
│   └── version/        # Version information
├── docker/             # Docker build files
├── docs/               # Documentation
├── scripts/            # Build and utility scripts
├── cliff.toml          # Changelog generation configuration
├── go.mod              # Go module definition
└── CLAUDE.md           # AI assistant guidance (this file)
```
