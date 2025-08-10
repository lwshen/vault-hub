package api

//go:generate python3 merge-spec.py
//go:generate go tool oapi-codegen -config cfg.yaml api.yaml
