#!/bin/bash

set -e

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Map darwin to macos for consistency
if [ "$OS" = "darwin" ]; then
  OS="darwin"
fi

# Construct download URL
BINARY_URL="https://github.com/lwshen/vault-hub/releases/latest/download/vault-hub-cli-${OS}-${ARCH}"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="vault-hub-cli"

echo "Downloading vault-hub-cli for ${OS}-${ARCH}..."
curl -fsSL "$BINARY_URL" -o "/tmp/${BINARY_NAME}"
chmod +x "/tmp/${BINARY_NAME}"

echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "vault-hub-cli installed successfully!"
echo "Run 'vault-hub-cli version' to verify."