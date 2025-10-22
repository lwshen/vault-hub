#!/usr/bin/env bash
set -euo pipefail

command -v pnpm >/dev/null 2>&1 || {
  echo "pnpm is required but not found in PATH" >&2
  exit 1
}

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Updating apps/web submodule..."
git -C "${REPO_ROOT}" submodule update --init --remote apps/web

echo "Installing frontend dependencies..."
pnpm --dir "${REPO_ROOT}/apps/web" install

echo "Building frontend assets..."
EMBED_DIST="${REPO_ROOT}/internal/embed/dist"
pnpm --dir "${REPO_ROOT}/apps/web" run build
mkdir -p "${EMBED_DIST}"
cp -a "${REPO_ROOT}/apps/web/dist/." "${EMBED_DIST}/"

echo "apps/web update complete."
