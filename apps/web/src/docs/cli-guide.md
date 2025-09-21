# CLI Guide

The VaultHub CLI is the primary way to interact with VaultHub. It provides secure, programmatic access to your vaults and is perfect for development workflows, CI/CD pipelines, and automation. Start here to quickly access your secrets without needing to set up a server.

## CLI Installation

### Download Pre-built Binaries

Download the latest CLI binary for your platform from the [GitHub releases page](https://github.com/lwshen/vault-hub/releases/latest).

**Latest Release**

Download the latest version of VaultHub CLI from our GitHub releases page:

[👉 Download Latest Release](https://github.com/lwshen/vault-hub/releases/latest)

**Supported Platforms:**
- **Linux** - amd64, arm64
- **Windows** - amd64  
- **macOS** - amd64, arm64

### Build from Source

```bash
# Clone the repository
git clone https://github.com/lwshen/vault-hub.git
cd vault-hub

# Build the CLI
go build -o vault-hub-cli ./apps/cli/main.go

# Make it executable and move to PATH (Linux/macOS)
chmod +x vault-hub-cli
sudo mv vault-hub-cli /usr/local/bin/vault-hub
```

## Authentication

The CLI uses API keys for authentication. First, create an API key in the web interface:

1. Log into the VaultHub web interface
2. Navigate to Dashboard → API Keys
3. Click "Create API Key" and give it a name
4. Copy the generated API key (starts with `vhub_`)

### Setting Up Authentication

```bash
# Set the API key as an environment variable
export VAULT_HUB_API_KEY=vhub_your_api_key_here

# Or pass it directly to commands
vault-hub --api-key vhub_your_api_key_here list
```

## Commands

### List Vaults

```bash
# List all accessible vaults
vault-hub list

# Short form
vault-hub ls
```

### Get Vault Contents

```bash
# Get vault by name
vault-hub get --name production-secrets

# Get vault by ID
vault-hub get --id vault-uuid-here

# Export to .env file
vault-hub get --name production-secrets --output .env

# Execute command with environment variables
vault-hub get --name production-secrets --exec "npm start"
```

### Version Information

```bash
# Show version and build information
vault-hub version
```

## Example Workflows

### Development Workflow

```bash
# Get development secrets and start your app
vault-hub get --name dev-secrets --exec "npm run dev"

# Export secrets to .env file for local development
vault-hub get --name dev-secrets --output .env
```

### CI/CD Pipeline

```bash
# In your CI/CD pipeline
export VAULT_HUB_API_KEY=${{ secrets.VAULT_HUB_API_KEY }}
vault-hub get --name production-secrets --exec "docker build -t myapp ."
```