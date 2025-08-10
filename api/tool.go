package api

//go:generate sh -c "redocly bundle openapi/main.yaml -o api.yaml --force && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml"
