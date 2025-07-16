# Encryption in VaultHub

VaultHub now implements AES-256-GCM encryption for all configuration values stored in the database, ensuring that sensitive data like API keys and secrets are securely protected.

## Quick Start

**TL;DR**: Generate and set your encryption key before starting VaultHub:

```bash
# Generate a secure encryption key
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Start VaultHub
./vault-hub-server

# Or for Docker
docker run -e ENCRYPTION_KEY=$(openssl rand -base64 32) vault-hub
```

**⚠️ Important**: Save your encryption key securely! If you lose it, you cannot decrypt your existing data.

## Overview

All configuration values are automatically encrypted before being stored in the database and decrypted when retrieved. This encryption happens transparently at the model layer, so API consumers receive plaintext values while the database stores only encrypted data.

## Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: SHA-256 hash of the `ENCRYPTION_KEY` environment variable
- **Authentication**: Built-in authentication via GCM mode
- **Encoding**: Base64 encoding for database storage
- **Nonce**: Randomly generated per encryption operation (ensures different ciphertext for identical plaintext)

## Configuration

### Environment Variables

Set the following environment variable to enable encryption:

```bash
ENCRYPTION_KEY=your-secure-encryption-key-here
```

### Generating a Secure Encryption Key

**Important**: Use a cryptographically secure method to generate your encryption key. Here are several recommended approaches:

#### Method 1: Using OpenSSL (Recommended)

```bash
# Generate a 256-bit (32-byte) random key and encode it as base64
openssl rand -base64 32

# Example output: 2xK9mNvB7qL8pRt4WyE6uI1oA5zC3dF0sG2hJ9kMnPqRsT
```

#### Method 2: Using /dev/urandom (Linux/macOS)

```bash
# Generate 32 random bytes and encode as base64
head -c 32 /dev/urandom | base64

# Or generate hex string
head -c 32 /dev/urandom | xxd -p -c 32
```

#### Method 3: Using Python

```python
import secrets
import base64

# Generate 32 random bytes and encode as base64
key = base64.b64encode(secrets.token_bytes(32)).decode('utf-8')
print(f"ENCRYPTION_KEY={key}")
```

#### Method 4: Using Node.js

```javascript
const crypto = require('crypto');

// Generate 32 random bytes and encode as base64
const key = crypto.randomBytes(32).toString('base64');
console.log(`ENCRYPTION_KEY=${key}`);
```

#### Method 5: Using Go

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    key := make([]byte, 32)
    rand.Read(key)
    fmt.Printf("ENCRYPTION_KEY=%s\n", base64.StdEncoding.EncodeToString(key))
}
```

#### Method 6: Online Secure Generator

For quick testing only (not recommended for production):
- Visit: https://www.random.org/passwords/?num=1&len=32&format=html&rnd=new
- Set length to 32+ characters with mixed case, numbers, and symbols

**Important Security Notes:**

1. **Use a strong encryption key**: The key should be at least 32 characters long and contain a mix of letters, numbers, and special characters.
2. **Keep the key secure**: Never commit the encryption key to version control.
3. **Back up the key**: If you lose the encryption key, you will not be able to decrypt existing data.
4. **Use different keys for different environments**: Production, staging, and development should each have their own unique encryption keys.
5. **Generate keys securely**: Always use cryptographically secure random number generators.

### Example Configuration

```bash
# Production (use openssl rand -base64 32 to generate)
ENCRYPTION_KEY=2xK9mNvB7qL8pRt4WyE6uI1oA5zC3dF0sG2hJ9kMnPqRsT

# Development (generate a different key for dev)
ENCRYPTION_KEY=dev-mL3nK8pQ5rS7tU9vW2xY4zA6bC1dE0fG3hI5jK7mN9pQ

# Docker
docker run -e ENCRYPTION_KEY=$(openssl rand -base64 32) vault-hub

# Docker Compose
version: '3.8'
services:
  vault-hub:
    image: vault-hub
    environment:
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    # Set ENCRYPTION_KEY in your .env file
```

## What Gets Encrypted

- **Configuration Values**: The `value` field of all configurations is encrypted
- **API Responses**: Values are automatically decrypted when retrieved via the API
- **Database Storage**: Only encrypted values are stored in the database

## What Doesn't Get Encrypted

The following fields remain unencrypted for indexing and querying purposes:

- Configuration names
- Descriptions
- Categories
- Unique IDs
- User IDs
- Timestamps

## Implementation Details

### Encryption Process

1. When creating or updating a configuration, the plaintext value is encrypted using AES-256-GCM
2. A random nonce is generated for each encryption operation
3. The encrypted data is base64-encoded and stored in the database
4. The plaintext value is never stored in the database

### Decryption Process

1. When retrieving configurations, the encrypted value is fetched from the database
2. The base64-encoded data is decoded
3. The value is decrypted using AES-256-GCM
4. The plaintext value is returned to the API consumer

### Error Handling

If decryption fails (e.g., due to wrong encryption key or corrupted data):

- The operation will return an error
- No partial or corrupted data is returned
- Error messages are logged for debugging

## Security Considerations

### Key Rotation

To rotate the encryption key:

1. **WARNING**: This will make existing encrypted data unreadable
2. Export all configurations before changing the key
3. Update the `ENCRYPTION_KEY` environment variable
4. Re-import configurations (they will be encrypted with the new key)

### Performance

- Encryption/decryption adds minimal overhead (~1-2ms per operation)
- Nonce generation uses cryptographically secure random number generation
- GCM mode provides both encryption and authentication in a single operation

### Compliance

This implementation provides:

- **Confidentiality**: Data is protected with industry-standard AES-256 encryption
- **Integrity**: GCM mode ensures data hasn't been tampered with
- **Authentication**: Built-in authentication prevents unauthorized modifications

## Testing

The encryption implementation includes comprehensive tests:

- Unit tests for encryption/decryption functions
- Integration tests with database operations
- Error handling tests for invalid data
- Performance tests for large datasets

Run encryption tests:

```bash
export ENCRYPTION_KEY=test-key
go test ./internal/encryption/...
go test ./model/...
```

## Migration from Unencrypted Data

If you're upgrading from a version without encryption:

1. **Back up your database** before enabling encryption
2. Set the `ENCRYPTION_KEY` environment variable
3. Restart the application
4. Existing unencrypted data will need to be manually migrated

## Troubleshooting

### Common Issues

1. **"EncryptionKey is not set" error**
   - Solution: Set the `ENCRYPTION_KEY` environment variable

2. **"failed to decrypt" errors**
   - Check that the correct encryption key is being used
   - Verify that the data hasn't been corrupted
   - Ensure the key hasn't changed since the data was encrypted

3. **Performance issues**
   - Encryption is designed to be fast, but processing large numbers of configurations may take additional time
   - Consider implementing caching if needed

### Debug Mode

Enable debug logging to troubleshoot encryption issues:

```bash
export LOG_LEVEL=debug
```

This will provide detailed information about encryption operations without exposing sensitive data. 