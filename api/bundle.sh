#!/bin/bash
set -euo pipefail  # Exit on error, undefined variables, pipe failures

echo "🔧 Bundling OpenAPI files..."

# Check if redocly CLI is available
if ! command -v redocly &> /dev/null; then
    echo "❌ Error: redocly CLI not found."
    echo "Install with: npm install -g @redocly/cli"
    echo "Or use: npx @redocly/cli bundle openapi/api.yaml -o api.bundled.yaml"
    exit 1
fi

redocly bundle openapi/api.yaml -o api.bundled.yaml
echo "✅ OpenAPI files bundled successfully"
