# VaultHub Go Client

A Go client library for interacting with the VaultHub API.

## Installation

```bash
go get github.com/lwshen/vault-hub-go-client
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lwshen/vault-hub-go-client"
)

func main() {
    // Create a new client
    client, err := client.NewClient(
        "https://your-vault-hub-server.com",
        "your-api-key-here",
    )
    if err != nil {
        log.Fatal(err)
    }

    // List all accessible vaults
    vaults, err := client.ListVaults(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for _, vault := range vaults {
        fmt.Printf("Vault: %s (%s)\n", vault.Name, vault.UniqueID)
    }
}
```

### Advanced Configuration

```go
client, err := client.NewClient(
    "https://your-vault-hub-server.com",
    "your-api-key-here",
    client.WithTimeout(60*time.Second),
    client.WithHTTPClient(customHTTPClient),
)
```

### Available Methods

- `ListVaults(ctx context.Context) ([]VaultLite, error)` - List all accessible vaults
- `GetVault(ctx context.Context, uniqueID string) (*Vault, error)` - Get vault by unique ID
- `GetVaultByName(ctx context.Context, name string) (*Vault, error)` - Get vault by name
- `Health(ctx context.Context) error` - Check server health

## API Endpoints

The client interacts with the following VaultHub API endpoints:

- `GET /api/cli/vaults` - List vaults
- `GET /api/cli/vault/{uniqueId}` - Get vault by ID
- `GET /api/cli/vault/name/{name}` - Get vault by name
- `GET /health` - Health check

## Authentication

All API requests require an API key passed in the `Authorization` header as a Bearer token.