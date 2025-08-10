package api

//go:generate python3 merge.py
//go:generate go tool oapi-codegen -config cfg.yaml merged_api.yaml
