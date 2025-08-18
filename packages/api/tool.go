package api

//go:generate sh bundle.sh
//go:generate go tool oapi-codegen -config cfg.yaml api.bundled.yaml
