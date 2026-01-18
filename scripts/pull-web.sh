#!/usr/bin/env bash
set -euo pipefail

# Wrapper for scripts/update-web.sh to build assets without updating to the latest frontend.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "${SCRIPT_DIR}/update-web.sh" --no-update
