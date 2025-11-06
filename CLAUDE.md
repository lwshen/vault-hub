# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Go Backend (apps/server/)

- **Build**: `go build -o tmp/main ./apps/server/main.go`
- **Build with version**: `go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') -X github.com/lwshen/vault-hub/internal/version.Commit=$(git rev-parse --short HEAD)" -o tmp/main ./apps/server/main.go`
- **Run**: `go run ./apps/server/main.go` (launches API at http://localhost:3000 for local dev)
- **Test**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` (run all tests with required env vars)
- **Test specific package**: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./model -v`
- **Generate API code**: `go generate packages/api/tool.go` (run after modifying files in `packages/api/openapi/*`)

### Go CLI (apps/cli/)

- **Build**: `go build -o tmp/vault-hub-cli ./apps/cli/main.go` or `go build -o vault-hub-cli ./apps/cli/main.go`
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

- **Install dependencies**: `pnpm --dir apps/web install` (run from repo root)
- **Development server**: `pnpm --dir apps/web run dev`
- **Build production**: `pnpm --dir apps/web run build`
- **Lint**: `pnpm --dir apps/web run lint`
- **Type check**: `pnpm --dir apps/web run typecheck`
- **Preview build**: `pnpm --dir apps/web run preview`
- **Note**: `apps/web` is sourced from an external package; avoid editing its files unless coordinating with the frontend owners.
- **Sync frontend before builds**: `git submodule update --init --remote apps/web`

### Live Reload (Air)

- **Run watcher**: `air -c .air.toml`
- Air rebuilds the Go server and triggers `pnpm --dir apps/web run build -- --mode development --outDir ../../internal/embed/dist` so embedded static assets stay fresh.

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
- **Framework**: React 19.1.1 with TypeScript 5.9.2
- **Build tool**: Vite 7.1.5 with Tailwind CSS 4.1.13 (Lightning CSS)
- **Package manager**: pnpm 10.15.1
- **Routing**: Wouter (lightweight router)
- **UI**: Tailwind CSS 4.x + Radix UI components + Framer Motion for animations
- **API client**: Custom generated TypeScript client (`@lwshen/vault-hub-ts-fetch-client`)
- **State**: Zustand stores for component state, React Context for auth and theme management
- **Components**: Organized into dashboard, layout, modals, and UI components
- **Documentation System**: Built-in markdown-based documentation with TOC configuration
- **Features Page**: Marketing page showcasing VaultHub capabilities with advanced animations
- **Markdown Rendering**: `react-markdown` v10.1.0 with `remark-gfm` for GitHub Flavored Markdown
- **Typography**: `@tailwindcss/typography` v0.5.18 for prose styling
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

### Security & Configuration Tips

- **Keep secrets secure**: Store `JWT_SECRET`, `ENCRYPTION_KEY`, and database credentials in environment variables or secret storage systems
- **Never commit sensitive data**: Avoid committing real credentials, `data.db`, or any files containing secrets
- **Document configuration changes**: When adding new OIDC/database variables, document them in PRs
- **Sanitize shared snapshots**: Scrub sensitive rows from shared `data.db` snapshots
- **Use ephemeral test secrets**: Ensure secrets used in tests are ephemeral (e.g., `openssl rand`)

## Database Models

- **User**: User accounts with email/password or OIDC
- **Vault**: Encrypted key-value pairs for environment variables
- **AuditLog**: Audit trail of vault operations
- **APIKey**: API key management for programmatic access

## API Generation

The project uses OpenAPI 3.0 specification (`packages/api/openapi/api.yaml`) with `oapi-codegen` to generate:

- Go server stubs (`packages/api/generated.go`)
- TypeScript client library (published as npm package)

**Important**: After modifying files in `packages/api/openapi/*`:
1. **Bump the API version** in `packages/api/openapi/api.yaml` (update the `version` field in the `info` section) unless this branch has already updated the version relative to `main`
2. Run `go generate packages/api/tool.go` to regenerate the Go types and interfaces

The API spec uses camelCase naming convention for all properties (e.g., `uniqueId`, `createdAt`, `isActive`).

**NEVER EDIT**: Do not modify `packages/api/generated.go` or `packages/api/api.bundled.yaml` directly as they are auto-generated code. All API changes must be made in the OpenAPI specification files in `packages/api/openapi/*`.

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

### Backend Testing

- **Test files**: Place Go tests in `*_test.go` files next to their implementations
- **Coverage areas**: Config, encryption, and database flows must be tested
- **Go unit tests**:
  - Encryption tests: `internal/encryption/encryption_test.go`
  - Database model tests: `model/db_test.go`
  - Configuration tests: `internal/config/config_test.go`
- **Test secrets**: Ensure secrets used in tests are ephemeral (`openssl rand`); never commit real credentials or `data.db`
- **Environment variables**: All tests require `JWT_SECRET` and `ENCRYPTION_KEY` to be set

### Frontend Testing

- **Current**: Standard React testing patterns
- **Future**: Add Vitest + Testing Library with a `pnpm run test` script when introducing UI coverage

## Documentation System

The project includes a comprehensive built-in documentation system:

### Documentation Structure
- **Location**: `apps/web/src/docs/`
- **Format**: Markdown files with TypeScript TOC configuration
- **Sections**: CLI Guide, Server Setup, API Reference, Security
- **Navigation**: Hash-based routing with browser history support (e.g., `/docs#cli-guide`)

### Key Components
- **MarkdownContent**: Reusable component (`src/components/ui/markdown-content.tsx`) with configurable prose sizes
- **TOC Configuration**: Type-safe table of contents in `src/docs/toc.ts`
- **Markdown Rendering**: Uses `react-markdown` with `remark-gfm` for GitHub Flavored Markdown
- **Typography**: Tailwind CSS Typography plugin for consistent prose styling

### Documentation Files
- `cli-guide.md` - CLI installation, authentication, and usage examples
- `server-setup.md` - Server configuration and deployment
- `api-reference.md` - API endpoint documentation with OpenAPI references
- `security.md` - Security features, encryption, and best practices

### Features
- **URL-based Navigation**: Direct linking to sections with `/docs#section-id`
- **Dark Mode Support**: Automatic theme switching with `prose-invert`
- **Mobile Responsive**: Optimized for all screen sizes
- **Search Friendly**: Semantic HTML with proper heading structure

## Frontend State Management

The frontend uses Zustand for component-level state management:

- **Zustand stores**: Located in `src/stores/` for audit logs, API keys, and vaults
- **Store pattern**: Each store contains state, actions, and loading states with comprehensive error handling
- **Input validation**: All user inputs (pagination, deletion) include validation and error boundaries
- **API integration**: Stores directly use generated API clients with proper error handling
- **React Context**: Still used for global auth and theme state

## Coding Style & Naming Conventions

### Go Code Style

- **Formatting**: Always format Go code with `gofmt -w <files>` before committing
- **Exported types**: Use PascalCase for exported types and functions
- **Internal helpers**: Keep internal/private functions and types unexported (lowercase first letter)
- **Linting**: Run `golangci-lint run ./...` after editing Go code to ensure quality standards

### CLI Code Style

- **File naming**: CLI command files adopt hyphenated filenames (e.g., `list.go`)
- **Flag naming**: Use snake_case for flags
- **Implementation**: Logic resides in `internal/cli` with encryption utilities in `internal/encryption`

### Frontend Code Style

**ESLint configuration** enforces:

- 2-space indentation
- Single quotes
- Semicolons required
- Stylistic rules from `@stylistic/eslint-plugin`
- React-specific rules and hooks validation
- TypeScript strict mode

**Component conventions**:

- **Components**: Use PascalCase for React components
- **Hooks/utilities**: Use camelCase for hooks and utility functions
- **CSS**: Apply Tailwind classes inline; keep global CSS minimal
- **Location**: Components organized in `src/pages`, `src/components`, `src/stores`

### Generated Code

- **Commit policy**: Only commit generated artifacts when necessary
- **Regeneration**: Regenerate `packages/api` outputs after spec edits using `go generate packages/api/tool.go`

## Tailwind CSS 4.x Configuration

The project uses Tailwind CSS 4.x with the new CSS-first configuration approach:

### Configuration Method
- **No `tailwind.config.js`**: Uses CSS-first approach via `@import` and `@plugin` directives
- **Main CSS file**: `apps/web/src/index.css` contains all Tailwind configuration
- **Typography Plugin**: Added via `@plugin "@tailwindcss/typography";` directive
- **Vite Integration**: Uses `@tailwindcss/vite` plugin for seamless integration

### Important CSS Directives
```css
@import "tailwindcss";
@import "tw-animate-css";
@plugin "@tailwindcss/typography";
```

### Theme Configuration
- **CSS Custom Properties**: Extensive design tokens defined in `:root` and `.dark`
- **OKLCH Color Space**: Modern color system for better perceptual uniformity
- **Custom Variants**: Dark mode via `@custom-variant dark (&:is(.dark *))`

## CI/CD Pipeline

### GitHub Actions Workflows

The project uses multiple GitHub Actions workflows for comprehensive CI/CD:

#### Main CI Workflow (`.github/workflows/ci.yml`)
- **Triggers**: Push to main, pull requests to main
- **Go Version**: 1.24.2 with module caching
- **Frontend**: pnpm 10.15.1 with Node.js 22
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

## Vault Detail Page Implementation

### Recent UX Improvement (January 2025)

The vault viewing and editing experience was significantly improved by replacing modal dialogs with dedicated full-page views:

#### Previous Implementation (Modal-based)
- Used `ViewVaultValueModal` and `EditVaultValueModal` components
- Limited screen real estate, especially on mobile devices
- Cramped editing experience with small text areas

#### Current Implementation (Full-page)
- **Dedicated Route**: `/dashboard/vaults/:vaultId` with URL-based mode switching
- **Responsive Design**: Mobile-first approach with sticky action bar for better thumb access
- **Components**:
  - `VaultDetail` page wrapper using `DashboardLayout`
  - `VaultDetailContent` component containing all vault logic
- **Layout Structure**: Proper height management without scroll bar issues

#### Key Implementation Details

**Route Configuration** (`src/routes.tsx`):
```tsx
<Route path={PATH.VAULT_DETAIL}>
  {(params: { vaultId: string; }) => (
    <ProtectedRoute>
      <VaultDetail vaultId={params.vaultId} />
    </ProtectedRoute>
  )}
</Route>
```

**Mode Switching**: Uses URL query parameters (`?mode=edit`) for view/edit state
**Navigation Pattern**: `navigate(\`/dashboard/vaults/\${vault.uniqueId}\`)` from vault list
**Mobile UX**: Dedicated sticky action bar at bottom for better mobile interaction

#### Responsive Features
- **Desktop**: Header actions with text labels
- **Mobile**: Icon-only header actions + sticky bottom action bar
- **Textarea**: Responsive height (6 rows mobile, 8 rows tablet, 12 rows desktop)
- **Warnings**: Context-aware messages for edit/view modes

#### Files Modified
- **Deleted**: `view-vault-value-modal.tsx`, `edit-vault-value-modal.tsx`
- **Modified**: `vaults-content.tsx`, `dashboard-content.tsx`, `routes.tsx`, `path.ts`
- **Created**: `vault-detail.tsx`, `vault-detail-content.tsx`

This implementation provides a much better user experience with proper responsive design and eliminates the cramped modal limitations.

## Commit & Pull Request Guidelines

### Commit Message Format

- **Follow Conventional Commits**: Use prefixes like `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, etc.
- **Scope**: Optional but recommended for clarity (e.g., `chore(deps):`, `feat(cli):`, `fix(auth):`)
- **Examples**:
  - `feat(cli): add vault export command`
  - `fix(auth): resolve JWT token expiration issue`
  - `chore(deps): bump github.com/gofiber/fiber to v2.52.0`
  - `docs: update API reference with new endpoints`

### Pull Request Process

1. **Rebase**: Rebase onto `main` before opening a PR
2. **CI verification**: Ensure `.github/workflows/ci.yml` passes (all checks green)
3. **PR content**: Include:
   - Clear summary of changes and scope
   - Notes on schema or environment variable changes
   - Links to related tracking issues
   - CLI output or screenshots for UI updates
4. **Code quality**: Verify linting and tests pass before submitting

### Pre-PR Checklist

- [ ] Code follows conventional commit format
- [ ] All tests pass (`go test ./...`)
- [ ] Go linting passes (`golangci-lint run ./...`)
- [ ] Frontend linting passes (`pnpm --dir apps/web run lint`)
- [ ] Frontend type checking passes (`pnpm --dir apps/web run typecheck`)
- [ ] API version bumped if OpenAPI spec modified
- [ ] Documentation updated for new features/changes
- [ ] No sensitive data (credentials, `data.db`) committed

## Post-change Checklist

After making code changes, ensure you complete these steps:

### API Changes

- [ ] When modifying `packages/api/openapi/api.yaml`, bump the version in the `info.version` field (unless this branch has already updated it relative to `main`)
- [ ] After editing files in `packages/api/openapi/*`, run `go generate packages/api/tool.go` to regenerate code
- [ ] Do not manually modify `packages/api/api.bundled.yaml` or `packages/api/generated.go`

### Backend Changes

- [ ] Run `golangci-lint run ./...` after Go code changes
- [ ] Run `gofmt -w <files>` to format Go files
- [ ] Verify all tests pass: `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...`

### Frontend Changes

- [ ] Run `pnpm --dir apps/web run typecheck` to catch TypeScript issues
- [ ] Run `pnpm --dir apps/web run lint` to catch lint issues
- [ ] Sync frontend if needed: `git submodule update --init --remote apps/web`

### Pre-commit Verification

- [ ] No `data.db` or sensitive files staged for commit
- [ ] Conventional commit message format used
- [ ] All relevant tests and quality checks pass

## Module Organization

### Project Component Locations

- **Backend API**: `apps/server/` - Go Fiber API with routes in `handler/`, middleware in `route/`, and shared config in `internal/`
- **CLI**: `apps/cli/` - Cobra CLI backed by `internal/cli` logic and `internal/encryption` utilities
- **Frontend**: `apps/web/` - Vite + React UI (`src/pages`, `src/components`, `src/stores`); run UI assets through `pnpm`
  - **Important**: Do not edit files under `apps/web`; managed as an external codebase
- **Scheduled Jobs**: `apps/cron/` and `scripts/` - Supply scheduled jobs and release chores; keep them idempotent
- **API Specification**: `packages/api/` - Shared OpenAPI specs; regenerate clients with `go generate packages/api/tool.go`
- **Database Models**: `model/` - Reusable GORM models
- **Container Assets**: `docker/` - Docker build files

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
│       │   ├── docs/     # Documentation system
│       │   │   ├── cli-guide.md     # CLI installation and usage
│       │   │   ├── server-setup.md  # Server configuration
│       │   │   ├── api-reference.md # API endpoint documentation
│       │   │   ├── security.md      # Security features and best practices
│       │   │   └── toc.ts           # Table of contents configuration
│       │   ├── pages/    # Page components including features and documentation
│       │   │   └── dashboard/vault-detail.tsx # Vault detail page wrapper
│       │   ├── components/
│       │   │   ├── dashboard/vault-detail-content.tsx # Main vault detail logic
│       │   │   └── ui/markdown-content.tsx # Reusable markdown renderer
│       │   └── stores/   # Zustand state management
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
