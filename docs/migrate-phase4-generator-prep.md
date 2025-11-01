# Phase 4 – OpenAPI Generator Swap (Preparation)

## Goals
- Replace the `oapi-codegen`-based pipeline with the official OpenAPI Generator
  while preserving strict handler semantics required by the server and CLI.
- Stage tooling changes without breaking the existing Fiber-compatible
  artifacts, enabling phased verification alongside the Echo migration.
- Deliver reproducible scripts and configuration so CI and contributors can
  generate artifacts deterministically.

## Current Pipeline Findings
- `packages/api/tool.go` runs two `go:generate` directives: `bundle.sh` emits
  `api.bundled.yaml`, and `oapi-codegen` writes `generated.go` with Fiber router
  helpers plus shared models.
- The generated code bakes in Fiber-specific types (`fiber.Router`,
  `*fiber.Ctx`) that must be removed once the Echo routes cover the surface
  area.
- `packages/api/cfg.yaml` requests `fiber-server` and `strict-server` features,
  so any replacement must provide equivalent type-safe handler interfaces or a
  drop-in mapping layer.
- Downstream consumers (`vault-hub-go-client`, CLI, cron jobs) assume the
  existing model and client shape produced by `oapi-codegen`.

## Official Generator Evaluation
- Tooling options: the Node-based CLI (`npx @openapitools/openapi-generator-cli`)
  avoids bundling the Java JAR manually and fits the existing Redocly workflow.
  Docker/Standalone JAR remain fallback options if CI constraints require them.
- The `go-server` generator emits `net/http` handlers by default; custom
  templates (or a thin adapter) will be required to keep the strict handler
  pattern we rely on for Echo.
- The `go` generator can replace the published Go client once we validate
  compatibility with CLI consumers; configuration is staged but not yet wired.

## Proposed Workflow
1. Run `go generate ./packages/api/...`. The directive in `packages/api/tool.go`
   calls `generate-openapi.sh`, which bundles the spec and executes the official
   generator for both server and client outputs.
2. Inspect the generated code, iterate on template overrides, and add adapters
   as needed to reproduce strict handler behavior.
3. Once parity is proven, replace the `go:generate` directives and remove the
   `oapi-codegen` tool from `go.mod`.

## Configuration Overview
- `packages/api/tool.go` wires the new script into `go generate` so the preview
  artifacts are produced automatically.
- `packages/api/openapi-generator/config-go-server.yaml` configures module,
  package, and feature flags for the server artifacts. Outputs land in
  `packages/api/.openapi-generator/server`.
- `packages/api/openapi-generator/config-go-client.yaml` targets the Go client
  and writes to `packages/api/.openapi-generator/client`.
- `packages/api/generate-openapi.sh` orchestrates bundling and generation so
  contributors have a single entry point during experimentation.
- `.gitignore` now excludes `packages/api/.openapi-generator/` to keep the
  generated previews out of version control.

## Action Items Before Swap
- Finalize template overrides or adapters that translate the official generator
  output to the strict handler signatures used by Echo.
- Update `migrate-plan.md` phase 4 checklist with the new scripts/config and
  align CI tooling to call the official generator.
- Validate the generated Go client against the CLI to confirm compatibility
  (pagination, auth wrappers, encryption toggles).
- Remove redundant files once the swap is complete (e.g., `cfg.yaml`,
  `generated.go`, `tool.go` directives pointing at `oapi-codegen`).

## Validation Checklist (post-swap)
- `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...`
  passes with the new generator outputs.
- `golangci-lint run ./...` stays clean after integrating the official stubs.
- Echo routes compile and satisfy the generator interfaces without Fiber shims.
- CLI smoke tests (login, vault CRUD, API key flows) succeed against the updated
  server.
