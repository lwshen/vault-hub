# VaultHub

A secure solution for managing environment variables and API keys with web interface and command-line access.

## Features

- 🔐 **AES-256-GCM encryption** for all stored values
- 🌐 **Web interface** - Modern React application
- ⌨️ **Command-line interface** - Cross-platform CLI tool
- 🔑 **Multiple authentication** - JWT and API keys
- 📊 **Audit logging** - Complete operation history
- 🗄️ **Multiple databases** - SQLite, MySQL, PostgreSQL

## Quick Start

### Prerequisites
- Go 1.24+
- Node.js 22+ and pnpm (for web interface)

### Setup
```bash
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Required environment variables
export JWT_SECRET=your-jwt-secret
export ENCRYPTION_KEY=$(openssl rand -base64 32)
```

### Run Server
```bash
go run ./apps/server/main.go
# Server starts at http://localhost:3000
```

### Run Web Interface
```bash
cd apps/web
pnpm install && pnpm run dev
# Web interface at http://localhost:5173
```

### Use CLI
```bash
# Build CLI
go build -o vault-hub-cli ./apps/cli/main.go

# List vaults
./vault-hub-cli list

# Get specific vault
./vault-hub-cli get my-secrets
```

## Architecture

- **Backend**: Go + Fiber web server with OpenAPI 3.0 spec
- **Frontend**: React 19 + TypeScript + Vite + Tailwind CSS
- **CLI**: Go + Cobra with API key authentication
- **Database**: GORM with SQLite/MySQL/PostgreSQL support

## Security

- All vault values encrypted with AES-256-GCM before database storage
- JWT-based web authentication with optional OIDC support
- API key authentication for CLI access with `vhub_` prefix
- Complete audit trail of all operations
- Transparent encryption/decryption at model layer

## Environment Variables

**Required**:
- `JWT_SECRET` - JWT token signing secret
- `ENCRYPTION_KEY` - AES encryption key

**Optional**:
- `APP_PORT` - Server port (default: 3000)
- `DATABASE_TYPE` - sqlite|mysql|postgres (default: sqlite)
- `DATABASE_URL` - Database connection (default: data.db)

## License

Apache License 2.0 - see [LICENSE](LICENSE) file.