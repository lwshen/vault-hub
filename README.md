# VaultHub

A lightweight and secure solution for managing environment variables and API keys for your applications.

## Overview

VaultHub provides a simple and secure way to store, manage, and access sensitive vault data such as API keys and environment variables. All vault values are automatically encrypted using AES-256-GCM before being stored in the database, ensuring your sensitive data remains protected. It helps developers maintain security best practices while making vault management easier.

## Features

- 🔐 **AES-256-GCM encryption** for all vault values stored in the database
- 🔄 Simple API for storing and retrieving environment variables
- 🌍 Support for multiple environments (development, testing, production)
- 🧩 Easy integration with existing applications
- 🖥️ Command-line interface for convenient management
- 🔒 Security-focused design with best practices built in
- 🛡️ **Transparent encryption/decryption** at the model layer
- 🔑 **Secure key management** with environment variable setup

## Security

VaultHub implements industry-standard AES-256-GCM encryption for all vault values. See [ENCRYPTION.md](ENCRYPTION.md) for detailed information about the encryption implementation, key management, and security considerations.

**Important**: Set the `ENCRYPTION_KEY` environment variable before starting the application:

```bash
# Generate a secure encryption key
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Or set your own secure key
export ENCRYPTION_KEY=your-secure-encryption-key-here
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
