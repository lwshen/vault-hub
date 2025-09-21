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

# PostgreSQL configuration
export DATABASE_TYPE=postgres
export DATABASE_URL="postgres://user:password@localhost:5432/vaulthub?sslmode=disable"

# OIDC configuration (optional)
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
> Your encryption key should be kept secure and backed up safely.

## Production Deployment

### Docker Deployment

```bash
# Build the server
go build -o vault-hub-server ./apps/server/main.go

# Create a Dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY vault-hub-server .
CMD ["./vault-hub-server"]

# Build and run
docker build -t vault-hub .
docker run -p 3000:3000 \
  -e JWT_SECRET="your-jwt-secret" \
  -e ENCRYPTION_KEY="your-encryption-key" \
  vault-hub
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vault-hub
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
        image: vault-hub:latest
        ports:
        - containerPort: 3000
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
```

### Reverse Proxy Setup

```nginx
# Nginx configuration
server {
    listen 80;
    server_name vaulthub.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```