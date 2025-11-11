#!/usr/bin/env bash

set -euo pipefail

export COREPACK_ENABLE_DOWNLOAD_PROMPT=0

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${REPO_ROOT}"

echo "[post-start] Updating apps/web submodule"
git submodule update --init --remote apps/web

echo "[post-start] Enabling corepack and installing frontend dependencies"
corepack enable
pnpm --dir apps/web install --frozen-lockfile

echo "[post-start] Building frontend"
pnpm --dir apps/web run build

echo "[post-start] Building backend"
mkdir -p tmp
go build -o tmp/main ./apps/server/main.go

echo "[post-start] Installing Air"
go install github.com/air-verse/air@latest
