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
cd "${REPO_ROOT}/apps/web"
pnpm install --frozen-lockfile

echo "Building frontend assets..."
EMBED_DIST="${REPO_ROOT}/internal/embed/dist"
pnpm run build
mkdir -p "${EMBED_DIST}"
cp -a "${REPO_ROOT}/apps/web/dist/." "${EMBED_DIST}/"

echo "apps/web update complete."
