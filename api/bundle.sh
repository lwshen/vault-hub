#!/bin/bash
set -euo pipefail  # Exit on error, undefined variables, pipe failures

# Bundle the split OpenAPI files into a single file
echo "Bundling OpenAPI files..."
npx @redocly/cli bundle openapi/api.yaml -o api.bundled.yaml

# Check if bundling was successful
if [ $? -eq 0 ]; then
    echo "OpenAPI files bundled successfully"
else
    echo "Failed to bundle OpenAPI files"
    exit 1
fi
