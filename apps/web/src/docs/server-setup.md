# Server Setup

Set up and configure the VaultHub server for your team or organization. The server provides the web interface and API endpoints for vault management.

## Installation

VaultHub consists of a backend server and a web interface. You can run it locally or deploy it to your infrastructure.

### Prerequisites

- Go 1.24+ for the backend server
- Node.js 22+ and pnpm for the web interface (optional)
- Database: SQLite (default), MySQL, or PostgreSQL

### Quick Start

```bash
# Clone the repository
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Set required environment variables
export JWT_SECRET=your-jwt-secret-here
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Run the server
go run ./apps/server/main.go
```

## Configuration

VaultHub can be configured using environment variables. Here are the essential settings:

### Required Variables

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | Secret key for JWT token signing |
| `ENCRYPTION_KEY` | AES-256 encryption key for vault data |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | 3000 | Server port |
| `DATABASE_TYPE` | sqlite | Database type: sqlite, mysql, postgres |
| `DATABASE_URL` | data.db | Database connection string |
| `OIDC_CLIENT_ID` | - | OIDC client ID (optional) |
| `OIDC_CLIENT_SECRET` | - | OIDC client secret (optional) |
| `OIDC_ISSUER` | - | OIDC issuer URL (optional) |

### Example Configuration

```bash
# Basic configuration
export JWT_SECRET="your-super-secret-jwt-key"
export ENCRYPTION_KEY="$(openssl rand -base64 32)"
export APP_PORT=3000
export DATABASE_TYPE=sqlite
export DATABASE_URL=./data.db

# MySQL configuration
export DATABASE_TYPE=mysql
export DATABASE_URL="user:password@tcp(localhost:3306)/vaulthub?charset=utf8mb4&parseTime=True&loc=Local"

# PostgreSQL configuration
export DATABASE_TYPE=postgres
export DATABASE_URL="postgres://user:password@localhost:5432/vaulthub?sslmode=disable"

# OIDC configuration (optional - all three must be provided if using OIDC)
export OIDC_CLIENT_ID="your-oidc-client-id"
export OIDC_CLIENT_SECRET="your-oidc-client-secret"
export OIDC_ISSUER="https://your-oidc-provider.com"
```

## Creating Your First Vault

Once VaultHub is running, you can create your first vault through the web interface:

1. Navigate to `http://localhost:3000`
2. Register a new account or log in
3. Go to the Dashboard and click "Create Vault"
4. Enter a name and key-value pairs for your environment variables
5. Save your vault - all values are automatically encrypted

> **🔒 Security Note**  
> All vault values are encrypted with AES-256-GCM before being stored in the database.  
> Your encryption key should be kept secure and backed up safely. If you lose this key,  
> all encrypted vault data will be permanently inaccessible.

## Health Monitoring

VaultHub provides a comprehensive health check endpoint:

```bash
# Check system status
curl http://localhost:3000/api/status

# Response includes:
# - Application version and commit
# - Database health and performance metrics
# - System resource utilization
# - Overall health status: healthy, degraded, or unavailable
```

Use this endpoint for:
- Load balancer health checks
- Monitoring system alerts
- Container orchestration health probes
- CI/CD deployment verification

## Production Deployment

### Docker Deployment

#### Using Pre-built Images

```bash
# Run the server with Docker
docker run -d \
  --name vault-hub-server \
  -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  -e ENCRYPTION_KEY="$(openssl rand -base64 32)" \
  -e DATABASE_URL="/app/data/vault.db" \
  --restart unless-stopped \
  vaulthub/server:latest

# Health check
curl http://localhost:3000/api/status
```

#### Multi-stage Build (Custom)

```dockerfile
# Use the project's Dockerfile for a complete build
ARG NODE_VERSION=22
ARG GO_VERSION=1.24
ARG VERSION=dev
ARG COMMIT=unknown

# Frontend build stage
FROM node:${NODE_VERSION}-alpine AS frontend-builder
WORKDIR /app
COPY apps/web ./
RUN corepack enable && pnpm install --frozen-lockfile && pnpm build

# Backend build stage
FROM golang:${GO_VERSION}-alpine AS backend-builder
ARG VERSION
ARG COMMIT
WORKDIR /app
RUN apk add --no-cache gcc libc-dev
ENV GO111MODULE=on CGO_ENABLED=1 GOOS=linux
COPY . .
COPY --from=frontend-builder /internal/embed/dist ./internal/embed/dist
RUN go mod download
RUN go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=${VERSION} -X github.com/lwshen/vault-hub/internal/version.Commit=${COMMIT}" -o vault-hub-server apps/server/main.go

# Final stage
FROM alpine:3.22
WORKDIR /app
RUN addgroup -g 1001 -S vaultuser && adduser -u 1001 -S vaultuser -G vaultuser
COPY --from=backend-builder /app/vault-hub-server ./
RUN chown -R vaultuser:vaultuser /app
USER vaultuser
EXPOSE 3000
CMD ["./vault-hub-server"]
```

```bash
# Build and run custom image
docker build -t vault-hub-custom .
docker run -d \
  --name vault-hub-server \
  -p 3000:3000 \
  -v vault-hub-data:/app/data \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  -e ENCRYPTION_KEY="$(openssl rand -base64 32)" \
  -e DATABASE_URL="/app/data/vault.db" \
  --restart unless-stopped \
  vault-hub-custom
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  vault-hub:
    image: vaulthub/server:latest
    ports:
      - "3000:3000"
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
      - DATABASE_TYPE=postgres
      - DATABASE_URL=postgres://vaulthub:${DB_PASSWORD}@db:5432/vaulthub?sslmode=disable
    volumes:
      - vault-data:/app/data
    depends_on:
      - db
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/status"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=vaulthub
      - POSTGRES_USER=vaulthub
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U vaulthub"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  vault-data:
  postgres-data:
```

```bash
# Create .env file
echo "JWT_SECRET=$(openssl rand -base64 32)" > .env
echo "ENCRYPTION_KEY=$(openssl rand -base64 32)" >> .env
echo "DB_PASSWORD=$(openssl rand -base64 32)" >> .env

# Start services
docker-compose up -d

# Check health
docker-compose ps
curl http://localhost:3000/api/status
```

### Kubernetes Deployment

```yaml
# Create secrets first
apiVersion: v1
kind: Secret
metadata:
  name: vault-hub-secrets
type: Opaque
data:
  jwt-secret: <base64-encoded-jwt-secret>
  encryption-key: <base64-encoded-encryption-key>
  db-password: <base64-encoded-db-password>
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-hub-config
data:
  APP_PORT: "3000"
  DATABASE_TYPE: "postgres"
  DATABASE_URL: "postgres://vaulthub:$(DB_PASSWORD)@postgres:5432/vaulthub?sslmode=disable"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vault-hub
  labels:
    app: vault-hub
spec:
  replicas: 3
  selector:
    matchLabels:
      app: vault-hub
  template:
    metadata:
      labels:
        app: vault-hub
    spec:
      containers:
      - name: vault-hub
        image: vaulthub/server:latest
        ports:
        - containerPort: 3000
          name: http
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: vault-hub-secrets
              key: jwt-secret
        - name: ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: vault-hub-secrets
              key: encryption-key
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: vault-hub-secrets
              key: db-password
        envFrom:
        - configMapRef:
            name: vault-hub-config
        livenessProbe:
          httpGet:
            path: /api/status
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /api/status
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: vault-hub
  labels:
    app: vault-hub
spec:
  selector:
    app: vault-hub
  ports:
  - port: 80
    targetPort: 3000
    name: http
  type: ClusterIP
```

### Reverse Proxy Setup

```nginx
# Nginx configuration
server {
    listen 80;
    server_name vaulthub.yourdomain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name vaulthub.yourdomain.com;
    
    # SSL configuration
    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;
        
        # Health check endpoint
        proxy_connect_timeout 5s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health check for load balancer
    location /health {
        proxy_pass http://localhost:3000/api/status;
        access_log off;
    }
}
```

## Database Setup

### SQLite (Default)

```bash
# SQLite requires no additional setup
# Database file will be created automatically
export DATABASE_TYPE=sqlite
export DATABASE_URL=./data/vault.db

# Ensure data directory exists
mkdir -p ./data
```

### PostgreSQL

```bash
# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# Create database and user
sudo -u postgres psql << EOF
CREATE DATABASE vaulthub;
CREATE USER vaulthub WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE vaulthub TO vaulthub;
\q
EOF

# Configure VaultHub
export DATABASE_TYPE=postgres
export DATABASE_URL="postgres://vaulthub:your-secure-password@localhost:5432/vaulthub?sslmode=require"
```

### MySQL

```bash
# Install MySQL
sudo apt-get install mysql-server

# Create database and user
sudo mysql << EOF
CREATE DATABASE vaulthub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'vaulthub'@'localhost' IDENTIFIED BY 'your-secure-password';
GRANT ALL PRIVILEGES ON vaulthub.* TO 'vaulthub'@'localhost';
FLUSH PRIVILEGES;
EOF

# Configure VaultHub
export DATABASE_TYPE=mysql
export DATABASE_URL="vaulthub:your-secure-password@tcp(localhost:3306)/vaulthub?charset=utf8mb4&parseTime=True&loc=Local"
```

## OIDC Authentication Setup

VaultHub supports OpenID Connect (OIDC) for enterprise authentication.

### Auth0 Configuration

```bash
# Auth0 settings
export OIDC_CLIENT_ID="your-auth0-client-id"
export OIDC_CLIENT_SECRET="your-auth0-client-secret"
export OIDC_ISSUER="https://your-domain.auth0.com/"
```

### Okta Configuration

```bash
# Okta settings
export OIDC_CLIENT_ID="your-okta-client-id"
export OIDC_CLIENT_SECRET="your-okta-client-secret"
export OIDC_ISSUER="https://your-domain.okta.com/"
```

### Azure AD Configuration

```bash
# Azure AD settings
export OIDC_CLIENT_ID="your-azure-application-id"
export OIDC_CLIENT_SECRET="your-azure-client-secret"
export OIDC_ISSUER="https://login.microsoftonline.com/your-tenant-id/v2.0"
```

### OIDC Validation

VaultHub validates OIDC configuration on startup:
- All three OIDC variables must be set if any one is provided
- Server will exit with error if OIDC is partially configured
- OIDC discovery endpoint must be accessible

## Security Best Practices

### Environment Variables

```bash
# Generate secure secrets
JWT_SECRET=$(openssl rand -base64 32)
ENCRYPTION_KEY=$(openssl rand -base64 32)

# Store in secure location (never commit to version control)
echo "JWT_SECRET=$JWT_SECRET" > /etc/vaulthub/.env
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY" >> /etc/vaulthub/.env
chmod 600 /etc/vaulthub/.env
chown vaulthub:vaulthub /etc/vaulthub/.env
```

### Backup Strategy

```bash
# Backup encryption key (CRITICAL)
echo "$ENCRYPTION_KEY" | base64 > encryption-key.backup
# Store this backup in a secure, separate location

# Backup database
# SQLite
cp /path/to/vault.db vault-backup-$(date +%Y%m%d).db

# PostgreSQL
pg_dump -h localhost -U vaulthub vaulthub > vault-backup-$(date +%Y%m%d).sql

# MySQL
mysqldump -u vaulthub -p vaulthub > vault-backup-$(date +%Y%m%d).sql
```
