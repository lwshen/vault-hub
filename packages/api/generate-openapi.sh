#!/bin/sh
set -euo pipefail

# This script prepares the official OpenAPI Generator outputs consumed by the
# Go server and CLI. It is invoked by `go generate` (see packages/api/tool.go),
# bundles the spec, and then generates both server and client artifacts under
# packages/api/openapi/*.
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
SERVER_DIR="${SCRIPT_DIR}/openapi/server"
CLIENT_DIR="${SCRIPT_DIR}/openapi/client"

rm -rf "${SERVER_DIR}" "${CLIENT_DIR}"

echo "[openapi] Generating go-server artifacts..."
${CLI} generate \
	-g go-server \
	-c "${GENERATOR_DIR}/config-go-server.yaml"

echo "[openapi] Generating go client artifacts..."
${CLI} generate \
	-g go \
	-c "${GENERATOR_DIR}/config-go-client.yaml"

find "${SERVER_DIR}" -maxdepth 1 -type f \( -name 'go.mod' -o -name 'go.sum' -o -name 'main.go' -o -name 'Dockerfile' -o -name 'README.md' \) -delete
find "${CLIENT_DIR}" -maxdepth 1 -type f \( -name 'go.mod' -o -name 'go.sum' -o -name 'README.md' -o -name 'git_push.sh' -o -name '.travis.yml' \) -delete
rm -rf "${SERVER_DIR}/api" "${CLIENT_DIR}/api" "${CLIENT_DIR}/docs" "${SERVER_DIR}/.openapi-generator" "${CLIENT_DIR}/.openapi-generator"
rm -f "${CLIENT_DIR}/.gitignore" "${CLIENT_DIR}/.openapi-generator-ignore"
rm -rf "${SCRIPT_DIR}/.openapi-generator"
rm -f "${SCRIPT_DIR}/openapitools.json"

cat <<EOF
[openapi] Generation complete.
- Server artifacts: packages/api/openapi/server
- Client artifacts: packages/api/openapi/client

Review the outputs and update go:generate/tooling once the migration is ready.
EOF
