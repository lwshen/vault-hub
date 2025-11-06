# Phase 1 Discovery — Echo Adoption & Official OpenAPI Generator

**Author:** Codex assistant  
**Updated:** 2026-03-17  
**Scope:** Backend server (`apps/server`, `internal/server/echoapp/`, `packages/api/`), shared libraries, and API specification assets.

> This document captures the original findings from the Fiber → Echo discovery work and annotates the current Echo-first architecture. Treat all references to Fiber as historical context—the codebase now runs entirely on Echo with the official OpenAPI Generator pipeline.

---

## 1. Current Server Layout

- **Entry point:** `apps/server/main.go` boots the Echo server directly (no build tags). It configures logging, recovers, security middleware, and graceful shutdown.
- **Routing & middleware:** `internal/server/echoapp/` houses reusable middleware (`security.go`, `response.go`) and the consolidated route registration in `routes.go`. The directory replaces the old `handler/` and `route/` Fiber packages.
- **Static assets:** `echoapp.MountStatic` serves the embedded Vite build exported via `internal/embed`. HTML5 history mode is preserved for the SPA.
- **Security context contract:**  
  - Authenticated web routes set `context.Set("user", *model.User)` to expose the current user.  
  - CLI routes set `context.Set("api_key", *model.APIKey)` alongside `context.Set("user_id", *uint)` for downstream authorization checks.  
  - `SecurityMiddleware` mirrors the historic path gating (public, JWT-only, API-key-only) while returning the same JSON error envelope.

## 2. Framework-Agnostic Business Helpers

- Business logic in `packages/api/*.go` is now independent of web frameworks.
- Common helpers (`ExtractClientInfo`, `GetVaultsForUser`, `GetAuditLogsForUser`, etc.) are invoked by Echo handlers with explicit context-derived parameters (user, api key, client metadata).
- Fiber-specific helpers (`handler/response.go`, `handler/auth.go`, `route/*.go`) were removed; the Echo layer owns request parsing and response encoding.

## 3. OpenAPI Specification & Generation

- **Spec structure:** `packages/api/openapi/api.yaml` remains the source of truth, referencing modular `paths/` and `schemas/` files.
- **Bundling:** `packages/api/bundle.sh` uses Redocly CLI to emit `packages/api/api.bundled.yaml`.
- **Generation workflow:**  
  - `go generate ./packages/api/...` executes `packages/api/generate-openapi.sh`.  
  - The script bundles the spec and runs `npx @openapitools/openapi-generator-cli generate` twice—once for the Go server (`packages/api/openapi/server/go`) and once for the Go client (`packages/api/openapi/client`).  
  - The generated code is committed to the repo so downstream binaries can consume it without an additional generation step.
- **Configuration:** Generator options live under `packages/api/openapi-generator/*.yaml`, defining module paths `github.com/lwshen/vault-hub/packages/api/openapi/(server|client)`.
- **Downstream consumers:**  
  - The CLI now imports `github.com/lwshen/vault-hub/packages/api/openapi/client`.  
  - Echo handlers map to generator models through hand-maintained structs in `packages/api/types.go`, keeping response shapes stable while avoiding direct dependencies on generator packages.

## 4. Authentication & OIDC Notes

- **State management:** `internal/auth/oidc.go` implements an expiring in-memory state cache shared by Fiber and Echo during migration; with Fiber removed, Echo remains the sole consumer.
- **OIDC callback:** Echo routes (`loginOIDCHandler`, `loginOIDCCallbackHandler`) reuse the same helper logic, ensuring JWT issuance and audit logging behaviour remains consistent with the legacy implementation.
- **Magic links & password resets:** Echo handlers in `routes.go` call into `api.RequestPasswordResetEmail`, `api.RequestMagicLinkEmail`, and `api.ConsumeMagicLinkToken`, returning the same JSON envelopes as before (including `Retry-After` headers when rate-limited).

## 5. Integration & Testing Considerations

- **Parity checks:** Manual smoke tests confirmed login/signup/logout, password reset, magic link, vault CRUD, and CLI flows operate identically under Echo.
- **Pending automation:** Integration tests should target:  
  - Auth endpoints (password + OIDC) covering success and failure cases.  
  - CLI endpoints verifying client-side encryption toggles.  
  - Audit log listing/metrics to confirm query semantics (pagination, filters).  
  - Rate-limiting responses to ensure `Retry-After` headers surface correctly.
- **Performance baselines:** Echo performance matches the original Fiber build within previously agreed tolerances; continue monitoring `/api/vaults` and `/api/auth/login` latency in staging/production.

## 6. Tooling & Operational Checklist

- `go run ./apps/server/main.go` (Echo) and `go build -o vault-hub-cli ./apps/cli/main.go` remain the recommended local commands.
- `go generate ./packages/api/...` must succeed locally and in CI; ensure Node.js/npm are available for the generator CLI.
- `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` validates business helpers and database flows.
- `golangci-lint run ./...` should pass with no new exceptions.
- Docker builds continue embedding the frontend bundle via `pnpm --dir apps/web run build` followed by Go compilation; no changes required beyond the updated server entrypoint.

---

**Summary:** Phase 1 established a clear path to Echo and the official OpenAPI Generator. The migration is now complete: Fiber artifacts and dependencies have been removed, Echo serves as the only HTTP framework, and the generator outputs live under `packages/api/openapi/`. Future work should focus on hardened integration tests, template customisation, and client regeneration strategies.
