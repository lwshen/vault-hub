# VaultHub

A comprehensive secure environment variable and API key management system with AES-256-GCM encryption, featuring a modern web interface and cross-platform CLI tool.

## ✨ Features

### 🔐 Security First
- **AES-256-GCM encryption** for all vault values before database storage
- **JWT-based authentication** with optional OIDC support
- **API key authentication** for programmatic access
- **Complete audit logging** of all operations
- **Enhanced client-side encryption** for CLI with PBKDF2 key derivation

### 🌐 Web Interface
- **Modern React 19** application with TypeScript
- **Real-time dashboard** with vault management
- **System status monitoring** with health metrics
- **Built-in documentation** with interactive guides
- **Dark/light theme** support

### ⌨️ Command-Line Interface
- **Cross-platform binaries** (Linux, Windows, macOS)
- **Simple commands**: `list`, `get` with name/ID support
- **Environment file export** (.env file generation)
- **Command execution** with injected environment variables

### 🗄️ Database Support
- **SQLite** (default, zero-config)
- **MySQL** for production deployments
- **PostgreSQL** for enterprise use

## 🚀 Quick Start

### Prerequisites
- **Go 1.24+** for backend development
- **Node.js 22+ and pnpm** for web interface development

### 1. Clone and Setup
```bash
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Required environment variables
export JWT_SECRET=$(openssl rand -base64 64)
export ENCRYPTION_KEY=$(openssl rand -base64 32)
```

### 2. Run Backend Server
```bash
go run ./apps/server/main.go
# Server starts at http://localhost:3000
```

### 3. Run Web Interface (Development)
```bash
cd apps/web
pnpm install
pnpm run dev
# Web interface at http://localhost:5173
```

### 4. Build and Use CLI
```bash
# Build CLI
go build -o vault-hub-cli ./apps/cli/main.go

# Set API key (create one in web interface first)
export VAULT_HUB_API_KEY=vhub_your_api_key_here

# List all vaults
./vault-hub-cli list

# Get vault by name
./vault-hub-cli get --name production-secrets

# Export to .env file
./vault-hub-cli get --name dev-secrets --output .env

# Execute command with vault environment
./vault-hub-cli get --name dev-secrets --exec "npm start"
```

## 🏗️ Architecture

### Backend (Go)
- **Web Framework**: Fiber v2.52.9
- **Database ORM**: GORM v1.31.0
- **Authentication**: golang-jwt/jwt/v5 + optional OIDC
- **API**: OpenAPI 3.0 specification with auto-generated code
- **CLI**: Cobra v1.10.1 framework

### Frontend (React)
- **React**: 19.1.1 with TypeScript 5.9.2
- **Build Tool**: Vite 7.1.5 with Lightning CSS
- **Styling**: Tailwind CSS 4.1.13 + Radix UI components
- **State Management**: Zustand 5.0.8
- **Routing**: Wouter 3.7.1 (lightweight)
- **Animations**: Framer Motion 12.23.12

### API Architecture
- **Modular OpenAPI**: Separate path and schema files
- **Auto-generated clients**: Go server code + TypeScript client
- **Clear separation**: Web API (JWT) vs CLI API (API keys)
- **Published packages**: `@lwshen/vault-hub-ts-fetch-client` on npm

## 🔧 Development

### Build Commands

**Backend:**
```bash
# Run server
go run ./apps/server/main.go

# Build with version info
go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') -X github.com/lwshen/vault-hub/internal/version.Commit=$(git rev-parse --short HEAD)" -o tmp/main ./apps/server/main.go

# Run tests
JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...

# Generate API code (after modifying OpenAPI files)
go generate packages/api/tool.go
```

**Frontend:**
```bash
cd apps/web
pnpm install          # Install dependencies
pnpm run dev          # Development server
pnpm run build        # Production build
pnpm run lint         # ESLint
pnpm run typecheck    # TypeScript validation
```

**CLI:**
```bash
# Build CLI
go build -o vault-hub-cli ./apps/cli/main.go

# Cross-platform builds available via CI/CD
```

### Project Structure
```
vault-hub/
├── apps/
│   ├── server/           # Go backend (Fiber web server)
│   ├── cli/              # Go CLI (Cobra commands)
│   ├── web/              # React frontend (Vite + TypeScript)
│   └── cron/             # Go cron service
├── packages/api/         # OpenAPI 3.0 spec + generated code
├── internal/             # Internal Go packages
│   ├── auth/            # JWT + OIDC authentication
│   ├── encryption/      # AES-256-GCM encryption
│   ├── config/          # Configuration management
│   └── version/         # Version information
├── model/               # GORM database models
├── handler/             # HTTP request handlers
├── route/               # Routing and middleware
└── .github/workflows/   # CI/CD pipelines
```

## 🔒 Security

### Encryption
- **AES-256-GCM** encryption for all vault values
- **Unique IV** per encryption operation
- **AEAD** (Authenticated Encryption with Associated Data)
- **Client-side encryption** option for CLI with PBKDF2

### Authentication
- **JWT tokens** for web interface access
- **API keys** for CLI and programmatic access (prefix: `vhub_`)
- **Optional OIDC** integration for enterprise SSO
- **Route-based protection** with middleware enforcement

### Audit Trail
- **Complete operation history** in audit logs
- **User and API key attribution** for all actions
- **IP address and user agent** tracking
- **Queryable audit metrics** for compliance

## 🌍 Environment Variables

**Required:**
- `JWT_SECRET` - JWT token signing secret
- `ENCRYPTION_KEY` - AES-256 encryption key

**Optional:**
- `APP_PORT` - Server port (default: 3000)
- `DATABASE_TYPE` - sqlite|mysql|postgres (default: sqlite)
- `DATABASE_URL` - Database connection string
- `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_ISSUER` - OIDC configuration

## 📦 Installation

### Pre-built Binaries
Download the latest releases from [GitHub Releases](https://github.com/lwshen/vault-hub/releases/latest):

- `vault-hub-server-{platform}-{arch}` - Backend server
- `vault-hub-cli-{platform}-{arch}` - CLI tool

### Docker
```bash
# Pull from Docker Hub (when available)
docker pull vaulthub/server:latest
docker pull vaulthub/cli:latest
```

### Package Managers
```bash
# TypeScript client
npm install @lwshen/vault-hub-ts-fetch-client

# Go client (separate repository)
go get github.com/lwshen/vault-hub-go-client
```

## 📖 Documentation

VaultHub includes comprehensive built-in documentation accessible via the web interface:

- **CLI Guide** - Installation and usage examples
- **Server Setup** - Configuration and deployment
- **API Reference** - Complete endpoint documentation
- **Security** - Encryption and best practices

Access documentation at `/docs` in the web interface.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make changes and test: `pnpm run lint && pnpm run typecheck`
4. Commit with conventional commits: `feat: add new feature`
5. Push and create a Pull Request

### Code Quality
- **Go**: Run `golangci-lint run ./...` before committing
- **Frontend**: Use `pnpm run lint --fix` for auto-formatting
- **Tests**: Ensure all tests pass with required environment variables

## 📄 License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Repository**: https://github.com/lwshen/vault-hub
- **Releases**: https://github.com/lwshen/vault-hub/releases
- **Issues**: https://github.com/lwshen/vault-hub/issues
- **TypeScript Client**: https://www.npmjs.com/package/@lwshen/vault-hub-ts-fetch-client