# Phase 3 – Route & Handler Migration (Complete)

## Summary
- Completed parity coverage for the Echo server:
  - Public health/config/status/user endpoints share the same helpers with Fiber and are secured by the new middleware gate.
  - Vault CRUD and API-key management reuse framework-neutral logic (`packages/api/vault.go`, `packages/api/api_key.go`) while preserving audit logging.
  - Auth flows (`/api/auth/login|signup|logout`), password reset, and magic-link endpoints now call shared helper functions that consolidate validation, rate limiting, and audit logging.
  - CLI vault endpoints (`/api/cli/*`) delegate to the new helper layer and respect optional client-side encryption across frameworks.
  - Audit log listing/metrics endpoints run through unified parameter parsing and data builders for consistent pagination semantics.
  - OIDC login + callback are wired to Echo, backed by a framework-agnostic in-memory state cache and the existing user bootstrap logic.
- Introduced reusable helpers in `packages/api/auth.go`, `packages/api/audit_log.go`, and `packages/api/cli_vault.go` so Fiber and Echo share identical business rules.
- Replaced the Fiber session store in `internal/auth/oidc.go` with a lightweight expiring state cache, enabling both routers to participate in the OIDC flow safely.
- Expanded `internal/server/echoapp/routes.go` to mount the full authenticated surface while keeping static asset mounting isolated.

## Outstanding Work
- Adopt the official OpenAPI generator pipeline in Phase 4 and update the strict handlers accordingly.
- Add integration coverage for the shared helper layer and Echo routes (JWT/API-key/OIDC smoke tests) to guard behavior parity.
