# Phase 4 – OpenAPI Generator Swap

## Goals
- Replace the `oapi-codegen`-based pipeline with the official OpenAPI Generator
  while preserving strict handler semantics required by the server and CLI. ✅
- Stage tooling changes without breaking the existing runtime. ✅ All generation
  now flows through the new CLI.
- Deliver reproducible scripts and configuration so CI and contributors can
  generate artifacts deterministically. ✅

## Current Pipeline Findings
- `packages/api/tool.go` now runs a single `go:generate` directive that invokes
  `generate-openapi.sh` (bundles the spec and calls
  `@openapitools/openapi-generator-cli`).
- Generated server artifacts live in `packages/api/openapi/server/go`; the Go
  client lives in `packages/api/openapi/client`.
- `packages/api/types.go` maintains HTTP-friendly structs so business helpers
  can work with stable field semantics without importing the generated
  packages directly.
- The CLI has been repointed to the in-repo client module, eliminating the
  dependency on `github.com/lwshen/vault-hub-go-client`.

## Official Generator Evaluation
- Tooling choice: the Node-based CLI (`npx @openapitools/openapi-generator-cli`)
  is used in both local development and CI. No Docker wrapper required.
- The `go-server` generator produces framework-agnostic handlers; thin adapters
  around Echo are implemented manually rather than via template overrides.
- The `go` generator emits a Go client consumed by the CLI and available for
  downstream consumers.

## Workflow Summary
1. Run `go generate ./packages/api/...`. The directive in `packages/api/tool.go`
   calls `generate-openapi.sh`, which bundles the spec and executes the official
   generator for both server and client outputs.
2. Generated files populate `packages/api/openapi/server/go` and
   `packages/api/openapi/client`. The script prunes auxiliary files so the
   curated outputs can live in version control.
3. Business helpers reference the new models via the structs defined in
   `packages/api/types.go`.

## Configuration Overview
- `packages/api/tool.go` wires the generator script into `go generate`.
- `packages/api/openapi-generator/config-go-server.yaml` configures module,
  package, and feature flags for the server artifacts. Outputs land in
  `packages/api/openapi/server/go`.
- `packages/api/openapi-generator/config-go-client.yaml` targets the Go client
  and writes to `packages/api/openapi/client`.
- `packages/api/generate-openapi.sh` orchestrates bundling and generation so
  contributors have a single entry point during experimentation.
- `.gitignore` excludes generated directories underneath `packages/api/openapi/`
  to prevent accidental commits.

## Follow-up Items
- Monitor generator upgrades and pin new versions in
  `packages/api/openapi-generator/*.yaml` as needed.
- Evaluate template overrides if additional strict interfaces or validation
  wrappers are required.
- Document distribution strategy for external SDKs (e.g., publishing the
  generated Go client or adding language targets).

## Validation Checklist (post-swap)
- `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...`
  passes with the new generator outputs.
- `golangci-lint run ./...` stays clean after integrating the official stubs.
- Echo routes compile and satisfy the generator interfaces without Fiber shims.
- CLI smoke tests (login, vault CRUD, API key flows) succeed against the updated
  server.
