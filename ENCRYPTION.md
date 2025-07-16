# Encryption in VaultHub

VaultHub now implements AES-256-GCM encryption for all configuration values stored in the database, ensuring that sensitive data like API keys and secrets are securely protected.

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

**Important Security Notes:**

1. **Use a strong encryption key**: The key should be at least 32 characters long and contain a mix of letters, numbers, and special characters.
2. **Keep the key secure**: Never commit the encryption key to version control.
3. **Back up the key**: If you lose the encryption key, you will not be able to decrypt existing data.
4. **Use different keys for different environments**: Production, staging, and development should each have their own unique encryption keys.

### Example Configuration

```bash
# Production
ENCRYPTION_KEY=prod-2xK9mNvB7qL8pRt4WyE6uI1oA5zC3dF0sG2hJ9kM

# Development
ENCRYPTION_KEY=dev-test-encryption-key-for-development

# Docker
docker run -e ENCRYPTION_KEY=your-secure-key vault-hub
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