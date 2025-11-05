package api

//go:generate sh bundle.sh
//go:generate npx @openapitools/openapi-generator-cli generate -c openapi-generator-config.yaml
//go:generate rm -f generated/go.mod
