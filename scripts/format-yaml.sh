#!/bin/sh
set -eu

# For all YAML files
find .github/workflows -maxdepth 1 \( -name "*.yml" -o -name "*.yaml" \) -type f | while read -r file; do
    echo "Formatting $file"
    yq -i -P "$file"
done
