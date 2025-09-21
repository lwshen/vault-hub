# Security

VaultHub implements multiple layers of security to protect your sensitive data.

## Encryption

### AES-256-GCM Encryption

All vault values are encrypted using AES-256-GCM before being stored in the database. This provides both confidentiality and authenticity.

**Key Features:**
- 256-bit encryption key
- Galois/Counter Mode (GCM)
- Authenticated encryption
- Unique IV per encryption

**Security Benefits:**
- Data confidentiality
- Data integrity
- Authentication
- Resistance to tampering

> **⚠️ Important**  
> Keep your `ENCRYPTION_KEY` secure and backed up.  
> Without it, encrypted vault data cannot be recovered.

### Encryption Process

1. **Key Generation**: Use a cryptographically secure random key
   ```bash
   # Generate a secure encryption key
   openssl rand -base64 32
   ```

2. **Data Encryption**: Each vault value is encrypted individually
   - Unique initialization vector (IV) per encryption
   - AEAD (Authenticated Encryption with Associated Data)
   - Integrity verification through authentication tags

3. **Storage**: Only encrypted data is stored in the database
   - Original plaintext is never persisted
   - Encryption keys are never stored with data

## Access Control

### Authentication Methods

VaultHub supports multiple authentication mechanisms:

- **🔐 JWT tokens** for web access
- **🔑 API keys** for programmatic access  
- **🌐 OIDC integration** (optional)

### Route Protection

**Public Routes:**
- Login, registration, static assets

**JWT Protected:**
- Web interface, user management

**API Key Protected:**  
- CLI endpoints, vault access

### Permission Model

- **User-based access**: Each user has access to their own vaults
- **API key scoping**: API keys can be scoped to specific vaults
- **Fine-grained permissions**: Read/write access control
- **Session management**: Secure session handling with expiration

## Security Best Practices

### Key Management

```bash
# Generate secure keys
JWT_SECRET=$(openssl rand -base64 64)
ENCRYPTION_KEY=$(openssl rand -base64 32)

# Store securely (example with AWS Secrets Manager)
aws secretsmanager create-secret \
  --name "vaulthub/encryption-key" \
  --secret-string "$ENCRYPTION_KEY"
```

### Environment Variables

Never expose sensitive configuration:

```bash
# ❌ Don't do this
export ENCRYPTION_KEY=my-simple-key

# ✅ Use secure random keys
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# ✅ Use external secret management
export ENCRYPTION_KEY=$(aws secretsmanager get-secret-value \
  --secret-id "vaulthub/encryption-key" \
  --query SecretString --output text)
```

### Database Security

- **Encrypted connections**: Use TLS for database connections
- **Access control**: Limit database access to VaultHub service
- **Regular backups**: Encrypted backup strategy
- **Monitoring**: Database access logging

### Network Security

```nginx
# Use HTTPS in production
server {
    listen 443 ssl http2;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}
```

## Audit Logs

VaultHub maintains comprehensive audit logs of all operations for security monitoring and compliance.

### Logged Events

**Vault Operations:**
- Vault creation and deletion
- Vault value updates  
- Vault access (read operations)
- Permission changes

**Authentication Events:**
- User login and logout
- API key creation and usage
- Failed authentication attempts
- Session management

**System Events:**
- Configuration changes
- System startup/shutdown
- Error conditions

### Audit Log Format

```json
{
  "timestamp": "2025-01-01T12:00:00Z",
  "userId": "user-uuid",
  "action": "vault.read",
  "resource": "vault-uuid",
  "ipAddress": "192.168.1.100",
  "userAgent": "vault-hub-cli/1.2.9",
  "result": "success",
  "metadata": {
    "vaultName": "production-secrets",
    "accessMethod": "api-key"
  }
}
```

### Audit Log Query

Access audit logs through the web interface or API:

```http
GET /api/audit-logs
Authorization: Bearer jwt_token_here

Query Parameters:
- action: Filter by action type
- userId: Filter by user
- startDate: Start date filter
- endDate: End date filter
- limit: Result limit
- offset: Pagination offset
```

## Compliance

### Standards Compliance

VaultHub is designed to help organizations meet security and compliance requirements:

- **SOC 2 Type II**: Security, availability, and confidentiality
- **GDPR**: Data protection and privacy
- **HIPAA**: Healthcare data protection (when properly configured)
- **ISO 27001**: Information security management

### Data Protection

- **Data minimization**: Only necessary data is collected
- **Right to deletion**: Complete data removal capability
- **Data portability**: Export functionality
- **Breach notification**: Audit trail for incident response

### Security Controls

| Control | Implementation |
|---------|----------------|
| Encryption at Rest | AES-256-GCM for all vault data |
| Encryption in Transit | TLS 1.3 for all connections |
| Authentication | Multi-factor support via OIDC |
| Authorization | Role-based access control |
| Audit Logging | Complete operation trail |
| Backup & Recovery | Encrypted backup procedures |
| Incident Response | Monitoring and alerting |

## Security Reporting

### Vulnerability Disclosure

If you discover a security vulnerability in VaultHub:

1. **Do not** disclose the vulnerability publicly
2. Email security details to the maintainers
3. Include steps to reproduce the issue
4. Allow reasonable time for response and fix

### Security Updates

- Subscribe to release notifications for security updates
- Monitor the [GitHub repository](https://github.com/lwshen/vault-hub) for security advisories
- Test security updates in staging environments before production deployment

### Security Monitoring

Implement monitoring for:

```bash
# Failed authentication attempts
grep "authentication failed" /var/log/vaulthub.log

# Unusual access patterns
grep "vault.read" /var/log/vaulthub.log | awk '{print $4}' | sort | uniq -c

# API key usage
grep "api-key" /var/log/vaulthub.log | grep -v "success"
```