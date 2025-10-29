# Phase 3 – Route & Handler Migration (In Progress)

## Summary
- Connected the new Echo server to public API endpoints:
  - `GET /api/health` now served via Echo using shared `api.HealthCheck()` logic.
  - `GET /api/config` uses `api.PublicConfig()` to return feature flags.
  - `GET /api/status` reuses `api.BuildStatusResponse()` for operational insights.
  - `GET /api/user` now leverages `api.BuildCurrentUserResponse()` and the Echo security middleware to return the authenticated profile.
  - `GET/POST/PUT/DELETE /api/vaults` flow through shared `api.GetVaultsForUser` et al., returning parity responses and audit logging via reusable helpers.
  - `GET/POST/PATCH/DELETE /api/api-keys` use new framework-agnostic helpers (`api.CreateAPIKeyForUser`, etc.) to manage API keys and audit trails.
- Extracted response builders in `packages/api` so both Fiber and Echo handlers can consume the same business logic without duplicating state checks.
- Echo server now registers routes through `internal/server/echoapp/routes.go`, while static assets mount after route registration to avoid conflicts.

## Outstanding Work
- Port remaining endpoints (auth flows, CLI vault endpoints, audit logs, magic link/OIDC callbacks) to Echo, ensuring middleware injects request context data similar to Fiber’s `Locals`.
- Replace Fiber session middleware in `internal/auth/oidc.go` with an Echo-compatible implementation before exposing OIDC endpoints.
- Update generated strict handlers (`packages/api/generated.go`) and request parsing once the official OpenAPI Generator is adopted in Phase 4.
- Expand automated tests to cover both health/status/config endpoints via the new Echo server.
