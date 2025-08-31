#!/bin/sh
set -eu

find . -maxdepth 5 \
  \( -name "*.yml" -o -name "*.yaml" \) \
  -type f \
  ! -path "./apps/web/*" \
  ! -path "./packages/api/api.bundled.yaml" \
  -print0 | \
  xargs -0 npx prettier --write
