# VaultHub Docker Usage

This document describes Docker usage for VaultHub, including the server, CLI, and base development images.

## Docker Files Overview

VaultHub includes three Dockerfile configurations:

- **`Dockerfile`** - Multi-stage build for the full VaultHub server with React frontend
- **`Dockerfile-cli`** - Lightweight CLI container with cronjob and one-time execution support  
- **`docker/Dockerfile-base`** - Development base image with Go, Node.js, and build tools
- **`docs/docker.md`** - This documentation file

## VaultHub Server (Dockerfile)

The main Dockerfile builds a complete VaultHub server with embedded React frontend.

### Build the Server Container

```bash
# Build the server container
docker build -t vault-hub-server .
```

### Run the Server

```bash
docker run -d \
  -p 3000:3000 \
  -e JWT_SECRET=your-jwt-secret \
  -e ENCRYPTION_KEY=$(openssl rand -base64 32) \
  -v vault-hub-data:/app/data \
  --name vault-hub-server \
  vault-hub-server
```

### Server Environment Variables

Required:
- `JWT_SECRET` - Secret for JWT token signing
- `ENCRYPTION_KEY` - AES-256 encryption key

Optional:
- `APP_PORT` (default: 3000)
- `DATABASE_TYPE` (sqlite|mysql|postgres, default: sqlite)
- `DATABASE_URL` (default: data.db)
- OIDC settings: `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_ISSUER`

## VaultHub CLI (Dockerfile-cli)

The CLI Dockerfile creates a lightweight container for running VaultHub CLI commands with support for both one-time execution and scheduled cronjobs.

### Build the CLI Container

```bash
# Build the CLI container
docker build -f Dockerfile-cli -t vault-hub-cli .
```

## Usage Modes

The CLI container supports two run modes via the `RUN_MODE` environment variable:

### 1. One-shot Mode (`RUN_MODE=oneshot`)

Executes the CLI command once and exits. Useful for:
- Manual vault operations
- CI/CD pipelines
- One-time data retrieval

**Environment Variables:**
- `RUN_MODE=oneshot` (default)
- `VAULT_HUB_CLI_ARGS` - CLI arguments to execute (default: `list`)
- `VAULT_HUB_SERVER_URL` - VaultHub server URL
- `VAULT_HUB_API_KEY` - API key for authentication

**Example:**
```bash
docker run --rm \
  -e RUN_MODE=oneshot \
  -e VAULT_HUB_CLI_ARGS="get --name my-secrets" \
  -e VAULT_HUB_SERVER_URL=https://vault-hub.example.com \
  -e VAULT_HUB_API_KEY=vhub_xxx \
  vault-hub-cli
```

### 2. Cron Mode (`RUN_MODE=cron`)

Runs the CLI command on a schedule using cron. Useful for:
- Periodic vault synchronization
- Scheduled backups
- Monitoring and alerting

**Environment Variables:**
- `RUN_MODE=cron`
- `CRON_SCHEDULE` - Cron expression (default: `0 * * * *` - every hour)
- `VAULT_HUB_CLI_ARGS` - CLI arguments to execute (default: `list`)
- `VAULT_HUB_SERVER_URL` - VaultHub server URL
- `VAULT_HUB_API_KEY` - API key for authentication

**Example:**
```bash
docker run -d \
  -e RUN_MODE=cron \
  -e CRON_SCHEDULE="0 */6 * * *" \
  -e VAULT_HUB_CLI_ARGS="list" \
  -e VAULT_HUB_SERVER_URL=https://vault-hub.example.com \
  -e VAULT_HUB_API_KEY=vhub_xxx \
  -v vault-hub-logs:/var/log/cron \
  vault-hub-cli
```

## Advanced Usage

For complex deployments, you can create your own docker-compose configuration or orchestration setup based on the examples above.

## Available CLI Commands

The CLI container supports all VaultHub CLI commands:

- `list` or `ls` - List all accessible vaults
- `get --name <vault-name>` - Get specific vault by name
- `get --id <vault-id>` - Get specific vault by unique ID
- `get --name <vault> --output <file>` - Save vault to file
- `get --name <vault> --exec "command"` - Execute command if vault updated
- `version` - Show version information

## Cron Schedule Examples

- `0 * * * *` - Every hour
- `0 */6 * * *` - Every 6 hours
- `0 9 * * 1-5` - 9 AM, Monday to Friday
- `0 0 * * 0` - Every Sunday at midnight
- `*/15 * * * *` - Every 15 minutes

## Log Management

In cron mode, logs are written to `/var/log/cron/vault-hub.log` inside the container. Mount this as a volume for persistence:

```bash
-v vault-hub-logs:/var/log/cron
```

View logs:
```bash
docker exec -it <container-name> tail -f /var/log/cron/vault-hub.log
```

## Security Considerations

1. **API Keys**: Store API keys securely using Docker secrets or environment files
2. **Network**: Use Docker networks to isolate containers
3. **Volumes**: Mount logs and data to persistent volumes
4. **Updates**: Regularly update the base Alpine image for security patches

## Development Base Image (docker/Dockerfile-base)

The base image provides a pre-configured development environment with all necessary tools for building VaultHub components.

### What's Included

- **Go 1.24** - Latest Go compiler
- **Node.js 22** - JavaScript runtime  
- **pnpm** - Fast package manager
- **golangci-lint** - Go code quality tools
- **OpenAPI Generator** - API client generation
- **Java 17** - Required for OpenAPI Generator
- **Build tools** - gcc, build-essential, git, curl, wget

### Build the Base Image

```bash
# Build locally
docker build -f docker/Dockerfile-base -t vault-hub-base .

# Or use the build script for multi-platform builds
cd docker
./build-base-image.sh
```

### Using in CI/CD

The base image is designed for use in GitLab CI and other CI/CD systems:

```yaml
# .gitlab-ci.yml example
image: registry.gitlab.com/your-group/vault-hub/vault-hub-base:latest

stages:
  - build
  - test

build:
  stage: build
  script:
    - go build ./...
    - pnpm install
    - pnpm build
```

### Multi-Platform Support

The base image supports multiple architectures:
- `linux/amd64` - Intel/AMD 64-bit
- `linux/arm64` - ARM 64-bit (Apple Silicon, ARM servers)
- `linux/arm/v7` - ARM 32-bit (Raspberry Pi)

## Troubleshooting

### Check container logs:
```bash
docker logs <container-name>
```

### Debug cron jobs:
```bash
# Check if cron is running
docker exec -it <container-name> ps aux | grep cron

# Check cron configuration
docker exec -it <container-name> cat /etc/crontabs/root

# Manual CLI execution
docker exec -it <container-name> ./vault-hub-cli list
```

### Environment variables:
```bash
# Check environment inside container
docker exec -it <container-name> env | grep VAULT_HUB
```