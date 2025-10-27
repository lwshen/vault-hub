#!/bin/bash

# OpenAPI Generator CLI wrapper script
# This script downloads and runs the OpenAPI Generator CLI

set -euo pipefail

# Directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# OpenAPI Generator version and download URL
OPENAPI_GENERATOR_VERSION="7.8.0"
OPENAPI_GENERATOR_JAR="$SCRIPT_DIR/openapi-generator-cli.jar"
OPENAPI_GENERATOR_URL="https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/${OPENAPI_GENERATOR_VERSION}/openapi-generator-cli-${OPENAPI_GENERATOR_VERSION}.jar"

# Download JAR if not exists
if [ ! -f "$OPENAPI_GENERATOR_JAR" ]; then
    echo "Downloading OpenAPI Generator CLI v${OPENAPI_GENERATOR_VERSION}..."
    curl -L -o "$OPENAPI_GENERATOR_JAR" "$OPENAPI_GENERATOR_URL"
    echo "Downloaded OpenAPI Generator CLI to $OPENAPI_GENERATOR_JAR"
fi

# Change to project root for relative paths
cd "$PROJECT_ROOT"

# Run OpenAPI Generator CLI with all passed arguments
java -jar "$OPENAPI_GENERATOR_JAR" "$@"