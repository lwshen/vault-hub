# Phase 1 Discovery — Fiber → Echo & OpenAPI Generator Migration

**Author:** Codex assistant  
**Date:** 2025-10-26  
**Scope:** Backend server (`apps/server`, `route/`, `packages/api/`), shared libraries, and API specification assets.

---

## 1. Fiber Usage Inventory

- **Legacy handlers:** `handler/auth.go`, `handler/response.go` still rely on `github.com/gofiber/fiber/v2`; they contain the vintage OIDC flow and shared `SendError` helper.
- **Route registration & middleware (archived):** `route/middleware.go.fiber_old`, `route/route.go.fiber_old`, and multiple `.fiber_old` files in `packages/api/` expose the previous Fiber-based wiring for reference.
- **Generated artifacts:** `packages/api/generated.go.fiber_old` plus endpoint-specific `.fiber_old` files show how prior generation coupled handlers to Fiber contexts.
- **Module dependencies:** No Fiber references remain under `apps/server` or `internal/` aside from the files listed above; new Echo-first modules live in `route/echo_*` and `packages/api/echo_*`.
- **Dependency state:** `github.com/gofiber/fiber/v2 v2.52.9` is still listed in `go.mod`/`go.sum`, kept temporarily for the untouched legacy helpers.

## 2. Middleware & Helper Contracts

- **Authentication context contract:**  
  - JWT routes expect `ctx.Set("user", *model.User)` with password stripped.  
  - API key routes set `ctx.Set("user_id", *uint)` and `ctx.Set("api_key", *model.APIKey)` to drive downstream authorization.
- **Error envelope:** All handlers use `SendError` helpers to respond with `{ "error": { "code": <int>, "message": <string> } }`. Echo version lives in `packages/api/echo_helpers.go`; parity must be maintained for clients depending on the shape.
- **Audit logging expectations:** Handlers call into `model.LogUserAction` / `model.LogVaultAction` with `clientIP` and `userAgent` captured via helper. Equivalent Echo helper (`getClientInfoEcho`) already mirrors Fiber logic (`handler/auth.go:getClientInfo`).
- **OIDC flow:**  
  - State management uses signed cookies (`oauth_state`) with HMAC-SHA256 (secret derived from `config.JwtSecret`).  
  - Redirect targets the frontend hash fragment (`/login#token=...&source=oidc`), so Echo migration must preserve this redirect contract.
- **Static asset serving:** Echo router mounts embedded frontend assets via `middleware.StaticWithConfig` with HTML5 mode enabled. This replaces the Fiber filesystem middleware and must continue to serve SPA routes without authentication.
- **Header expectations:** CLI endpoints enforce `Authorization: Bearer vhub_*` API keys; JWT routes must reject API keys explicitly. Middleware replicates Fiber behaviour using prefix checks.

## 3. External & Third-Party Dependencies

- **Echo stack:**  
  - Core: `github.com/labstack/echo/v4`  
  - Built-in middleware: logger + recover in `apps/server/main.go`; route layer plans to incorporate CORS, compression, and rate limiting during later phases.
- **JWT parsing:** `github.com/golang-jwt/jwt/v5` shared between middleware stacks.
- **OIDC:** `github.com/coreos/go-oidc/v3/oidc` and `golang.org/x/oauth2` remain unchanged; only request/response plumbing varies by framework.
- **Static bundling:** `internal/embed` exposes compiled Vite assets; unaffected by the framework swap but must be reachable without conflicting with API middleware.

## 4. OpenAPI Specification Workflow & Consumers

- **Source of truth:** Split specification in `packages/api/openapi/` with primary `api.yaml` referencing modular `paths/` and `schemas/`.
- **Bundling step:** `packages/api/bundle.sh` uses `npx @redocly/cli bundle` to emit `packages/api/api.bundled.yaml`.
- **Generator configuration:**  
  - `packages/api/openapi-generator-config.yaml` targets `go-echo-server`, outputs to `packages/api/generated/`, and enables interface generation.  
  - `packages/api/tool.go` runs both bundling and generator via `go generate`.
  - `packages/api/openapitools.json` pins CLI version `7.16.0`.
- **Generated output layout:**  
  - `packages/api/generated/` contains the raw OpenAPI Generator project (Go modules, handlers, models) — currently unused at runtime.  
  - `packages/api/generated_models/` and `packages/api/generated_handlers/` hold curated subsets imported by bespoke Echo handlers.
- **Downstream consumers:**  
  - **CLI:** `internal/cli/client.go` imports `github.com/lwshen/vault-hub-go-client`; this client should eventually be regenerated via the official generator (likely a separate repo/module).  
  - **Server:** Echo handlers (`packages/api/echo_*.go`) leverage `generated_models` for request/response types.  
  - **Docs/SDKs:** No other automated consumers detected in-repo, but plan assumes potential TypeScript SDK/web app usage once generator outputs are expanded.

## 5. Acceptance Criteria & Guardrails

- **Functional parity:**  
  - All existing REST endpoints (24 total) must preserve request/response schemas defined in OpenAPI.  
  - Authentication flows (password, magic link, OIDC) continue to pass manual smoke tests.  
  - CLI workflows succeed against Echo server using current `vault-hub-go-client` bindings.
- **Automated checks:**  
  - `go test ./...` and `golangci-lint run ./...` succeed without new skips.  
  - `go build -o tmp/main ./apps/server/main.go` and `go build -o vault-hub-cli ./apps/cli/main.go` remain green.  
  - `go generate packages/api/tool.go` produces deterministic artifacts committed as needed.
- **Performance guardrails (baseline to be captured pre-cutover):**  
  - p95 latency for `/api/vaults` and `/api/auth/login` within ±10% of Fiber baseline under representative load.  
  - Memory usage of server process within ±15% of current steady-state.  
  - Binary size increase limited to <5 MB compared to latest Fiber build.
- **Operational requirements:**  
  - Logging, tracing, and metrics wiring available pre-cutover; no loss of correlation IDs or request IDs.  
  - Static asset serving remains functional (SPA routes resolved).  
  - Feature flag or environment toggle available to revert to Fiber build during verification stage.
- **Rollback triggers:**  
  - Authentication regressions (failed password/OIDC flows) not resolved within 2 hours.  
  - CLI critical path (`vault list`, `vault get`) broken or returning schema-incompatible payloads.  
  - Sustained error rate >1% or latency regression >20% after deployment.  
  - Inability to rebuild or regenerate OpenAPI artifacts deterministically.

## 6. Recommended Next Steps (Pre-Phase 2)

1. Capture baseline latency/memory metrics from current Fiber deployment or staging environment.  
2. Finalise feature flag or deployment toggle strategy (e.g., separate binaries, environment variable).  
3. Validate `go generate packages/api/tool.go` on CI and document tooling prerequisites (Node.js, `npx`).  
4. Confirm ownership for each endpoint group prior to migration sprints.  
5. Socialise acceptance criteria with stakeholders (API consumers, DevOps) and secure sign-off.

