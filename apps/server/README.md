# VaultHub Server

This directory contains the main server application for VaultHub.

## Overview

The VaultHub server is a Go-based web application built with the Fiber framework that provides:

- Secure vault management with AES-256-GCM encryption
- JWT and OIDC authentication
- API key management for programmatic access
- RESTful API for vault operations
- Audit logging for all vault operations

## Running the Server

### Local Development

```bash
# Set required environment variables
export JWT_SECRET=your-jwt-secret
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Run the server
go run ./apps/server/main.go
```

### Production Build

```bash
# Build the server binary
go build -o vault-hub-server ./apps/server/main.go

# Run the binary
./vault-hub-server
```

## Configuration

The server requires the following environment variables:

- `JWT_SECRET` - Secret for JWT token signing (required)
- `ENCRYPTION_KEY` - AES-256 encryption key (required)
- `APP_PORT` - Server port (default: 3000)
- `DATABASE_TYPE` - Database type: sqlite|mysql|postgres (default: sqlite)
- `DATABASE_URL` - Database connection string (default: data.db)

Optional OIDC configuration:
- `OIDC_CLIENT_ID` - OIDC client ID
- `OIDC_CLIENT_SECRET` - OIDC client secret
- `OIDC_ISSUER` - OIDC issuer URL

## Architecture

The server follows a clean architecture pattern:

- **Entry Point**: `main.go` - Server initialization and startup
- **Routes**: `../../route/` - HTTP routing and middleware
- **Handlers**: `../../handler/` - Request/response handling
- **Models**: `../../model/` - Database entities and operations
- **Internal**: `../../internal/` - Configuration, authentication, encryption

## API Documentation

The server exposes a RESTful API documented in OpenAPI 3.0 format. See `../../api/openapi/api.yaml` for the complete specification.