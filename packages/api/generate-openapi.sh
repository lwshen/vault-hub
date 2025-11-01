#!/bin/sh
set -euo pipefail

# This script prepares the official OpenAPI Generator outputs without disturbing
# the current oapi-codegen pipeline. It is invoked by `go generate` (see
# packages/api/tool.go), bundles the spec, and then generates both server and
# client artifacts into .openapi-generator/* for inspection.
#
# Usage:
#   (cd packages/api && sh generate-openapi.sh)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
GENERATOR_DIR="${SCRIPT_DIR}/openapi-generator"
BUNDLED_SPEC="${SCRIPT_DIR}/api.bundled.yaml"

echo "[openapi] Bundling OpenAPI specification..."
sh "${SCRIPT_DIR}/bundle.sh"

if [ ! -f "${BUNDLED_SPEC}" ]; then
	echo "[openapi] bundled spec not found at ${BUNDLED_SPEC}" >&2
	exit 1
fi

CLI="npx @openapitools/openapi-generator-cli"

# TODO: pin CLI version once the target release is agreed upon.
echo "[openapi] Generating go-server artifacts..."
${CLI} generate \
	-g go-server \
	-c "${GENERATOR_DIR}/config-go-server.yaml"

echo "[openapi] Generating go client artifacts..."
${CLI} generate \
	-g go \
	-c "${GENERATOR_DIR}/config-go-client.yaml"

cat <<EOF
[openapi] Generation complete.
- Server artifacts: packages/api/.openapi-generator/server
- Client artifacts: packages/api/.openapi-generator/client

Review the outputs and update go:generate/tooling once the migration is ready.
EOF
