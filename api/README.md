# API Schema Organization

This directory contains the OpenAPI specification split into logical files for better organization while maintaining compatibility with `go generate`.

## Structure

### Main File
- `api.yaml` - Contains the complete OpenAPI spec with all schemas inline (used by oapi-codegen)

### Reference Files (for organization)
The following files contain the same schemas as the main file but are organized for easier maintenance:

#### Schema Files (`schemas/`)
- `health.yaml` - Health check related schemas
- `auth.yaml` - Authentication related schemas  
- `user.yaml` - User related schemas
- `vault.yaml` - Vault related schemas
- `audit.yaml` - Audit log related schemas
- `api-key.yaml` - API key related schemas

#### Path Files (`paths/`)
- `health.yaml` - Health check endpoints
- `auth.yaml` - Authentication endpoints
- `user.yaml` - User endpoints
- `vault.yaml` - Vault endpoints
- `audit.yaml` - Audit log endpoints
- `api-key.yaml` - API key endpoints

## Usage

### Code Generation
```bash
go generate tool.go
```

This command reads `api.yaml` and generates the Go code in `generated.go`.

### Making Changes

1. **Edit the relevant schema or path file** for organization and documentation
2. **Update the corresponding section in `api.yaml`** (the schemas must be kept inline for oapi-codegen)
3. **Run `go generate tool.go`** to regenerate the Go code

### Why This Structure?

- **Organization**: Logical separation of concerns makes the API easier to understand and maintain
- **Compatibility**: `oapi-codegen` works with the single `api.yaml` file containing inline schemas
- **Documentation**: Separate files serve as reference and make it easier to work on specific areas
- **Comments**: The main file includes comments pointing to the relevant separate files

## Scripts

- `sync-schemas.sh` - Helper script that documents the organization and workflow