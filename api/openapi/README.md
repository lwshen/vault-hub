# OpenAPI Specification Structure

This directory contains the OpenAPI specification split into multiple files for better organization and maintainability.

## Directory Structure

```
openapi/
├── api.yaml         # Main OpenAPI file with references to split files
├── paths/           # API endpoint definitions
│   ├── health.yaml  # Health check endpoint
│   ├── auth.yaml    # Authentication endpoints (login, signup, logout)
│   ├── user.yaml    # User endpoints
│   ├── vault.yaml   # Vault management endpoints
│   ├── audit.yaml   # Audit log endpoints
│   └── apikey.yaml  # API key management endpoints
└── schemas/         # Data model definitions
    ├── health.yaml  # Health check schemas
    ├── auth.yaml    # Authentication schemas
    ├── user.yaml    # User schemas
    ├── vault.yaml   # Vault schemas
    ├── audit.yaml   # Audit log schemas
    └── apikey.yaml  # API key schemas
```

## Path Reference Encoding

When referencing paths in external files, OpenAPI uses JSON Pointer notation. You'll see references like:

```yaml
$ref: './paths/auth.yaml#/~1api~1auth~1login'
```

The `~1` is the encoded form of `/` in JSON Pointer. Here's how to decode these references:

- `~1` → `/`
- `~0` → `~`

So the above reference points to the `/api/auth/login` path definition in the `auth.yaml` file.

### Examples:
- `#/~1api~1health` → references `/api/health`
- `#/~1api~1auth~1login` → references `/api/auth/login`
- `#/~1api~1vaults~1{uniqueId}` → references `/api/vaults/{uniqueId}`

This encoding is required by the OpenAPI specification when referencing paths in external files.

## How it Works

1. The main `api.yaml` file references the split files using `$ref` directives
2. When running `go generate`, the `bundle.sh` script is executed first
3. The script uses Redocly CLI to bundle all referenced files into `api.bundled.yaml`
4. The bundled file is then used by `oapi-codegen` to generate the Go code

## Making Changes

1. Edit the appropriate file in the `paths/` or `schemas/` directory
2. Run `go generate ./...` from the api directory
3. The bundled file and generated Go code will be automatically updated

## Note

- The `api.bundled.yaml` file is generated and should not be edited directly
- It is excluded from version control via `.gitignore`