# VaultHub

A lightweight and secure solution for managing environment variables and API keys for your applications.

## Overview

VaultHub provides a simple and secure way to store, manage, and access sensitive configuration data such as API keys and environment variables. All configuration values are automatically encrypted using AES-256-GCM before being stored in the database, ensuring your sensitive data remains protected. It helps developers maintain security best practices while making configuration management easier.

## Features

- 🔐 **AES-256-GCM encryption** for all configuration values stored in the database
- 🔄 Simple API for storing and retrieving environment variables
- 🌍 Support for multiple environments (development, testing, production)
- 🧩 Easy integration with existing applications
- 🖥️ Command-line interface for convenient management
- 🔒 Security-focused design with best practices built in
- 🛡️ **Transparent encryption/decryption** at the model layer
- 🔑 **Secure key management** with environment variable configuration

## Security

VaultHub implements industry-standard AES-256-GCM encryption for all configuration values. See [ENCRYPTION.md](ENCRYPTION.md) for detailed information about the encryption implementation, key management, and security considerations.

**Important**: Set the `ENCRYPTION_KEY` environment variable before starting the application:

```bash
export ENCRYPTION_KEY=your-secure-encryption-key-here
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
