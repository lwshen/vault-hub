#!/bin/sh
set -eu

# For all YAML files
find . -maxdepth 5 \( -name "*.yml" -o -name "*.yaml" \) -type f ! -path "./apps/web/*" ! -path "./packages/api/api.bundled.yaml" | while read -r file; do
    echo "Formatting $file"
    yq -i -P "$file"
done
