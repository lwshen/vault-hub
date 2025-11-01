package api

//go:generate sh generate-openapi.sh
//go:generate go tool oapi-codegen -config cfg.yaml api.bundled.yaml
