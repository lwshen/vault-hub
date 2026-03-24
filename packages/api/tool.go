package api

//go:generate sh bundle.sh
//go:generate go tool oapi-codegen -config cfg.yaml api.bundled.yaml
//go:generate go run ./cmd/fiber-v3-fix generated.go
