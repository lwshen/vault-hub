# Migration Plan: Fiber → Echo & Official OpenAPI Generator

## Purpose
- Adopt Echo to better match the team's preferred middleware patterns and gain access to Echo's mature ecosystem of adapters, contrib packages, and first-party support.
- Standardize on the official OpenAPI Generator tooling to improve alignment with upstream OpenAPI specifications and enable wider language/client support.
- Reduce long-term maintenance by unifying server, CLI, and cron code around a single HTTP framework and OpenAPI workflow.

## Overall Process
1. Inventory the current Fiber-based server, shared middleware, and OpenAPI generation flow to identify migration touchpoints.
2. Introduce an Echo baseline (dependencies, bootstrap wiring, middleware equivalents) behind feature flags or branch protection.
3. Incrementally port routes, handlers, and middleware from Fiber abstractions to Echo while keeping parity in behavior.
4. Replace the existing `oapi-codegen` pipeline with the official OpenAPI Generator CLI, regenerate server stubs and clients, and align build tooling.
5. Update dependent binaries (CLI, cron jobs) and shared packages to consume the new Echo handlers and generated API artifacts.
6. Run comprehensive verification (tests, lint, integration smoke tests) and update documentation before merging and releasing.

## Detailed Plan

### 1. Discovery & Preparation
- Audit `apps/server`, `handler`, `route`, and `internal` packages to catalog Fiber-specific constructs (router setup, context usage, request/response helpers, middleware such as `slog-fiber`, error handling helpers).
- Document non-HTTP concerns that touch Fiber types (authentication middleware, request validation, streaming endpoints) to ensure compatible Echo patterns exist.
- Review `packages/api/openapi/api.yaml` along with existing generation commands (`go generate packages/api/tool.go`) to map how `oapi-codegen` artifacts flow into server routing, models, and clients.
- Identify environment variables, configuration loaders, and dependency injection points that will need adjustments for the Echo startup sequence.
- Decide on migration sequencing (e.g., module-by-module, functionality-based, or via parallel Echo router) and capture rollback checkpoints.

### 2. Establish Echo Baseline
- Add Echo dependencies (`github.com/labstack/echo/v4`, logging/recovery middlewares) to `go.mod`; plan for removal of Fiber-specific modules and adapters.
- Create an Echo bootstrap entry point (e.g., `apps/server/main.go` or `internal/server/echo.go`) that initializes the Echo instance, config, and shared middleware.
- Implement substitutes for core middleware: request logging (Slog integration), panic recovery, CORS, authentication, context propagation, and request validation.
- Introduce adapter helpers where immediate parity is needed (for example, bridging Fiber-specific response helpers until handlers are ported).
- Verify the Echo server starts alongside existing Fiber implementation behind a feature flag or dedicated branch for incremental testing.

### 3. Migrate Routes, Handlers, and Middleware
- Translate route registration from Fiber's chaining syntax to Echo's group/route APIs; reorganize files under `apps/server/handler` and `apps/server/route` to follow the target structure.
- Update handler signatures to use `echo.Context`, replacing Fiber-specific helpers (e.g., `ctx.Params`, `ctx.Next`) with Echo equivalents.
- Refactor middleware to Echo's `echo.MiddlewareFunc` signature, ensuring consistent behavior for auth, rate limiting, and error translation.
- Replace response helpers (status codes, JSON rendering, streaming) with Echo's APIs, covering edge cases like file downloads or SSE.
- Remove or adapt Fiber-only utilities (e.g., `slog-fiber`) and add tests for critical handler paths to confirm behavior parity.

### 4. Adopt Official OpenAPI Generator
- ✅ Integrate the official CLI via `go generate ./packages/api/...`, emitting server/client artefacts under `packages/api/openapi/(server|client)`.
- ✅ Remove `oapi-codegen` tooling (`packages/api/generated.go`, `cfg.yaml`, go.mod tool dependency) and migrate downstream usage to framework-neutral structs in `packages/api/types.go`.
- ✅ Update CLI and build scripts to rely on the new generator outputs.
- 📌 Follow-up: pin CLI versions in `packages/api/openapi-generator/*.yaml`, evaluate template customisations if stricter interfaces are required, and document how external SDKs should consume the generated packages.

### 5. Update Dependent Applications & Libraries
- ✅ CLI (`apps/cli`) now consumes `github.com/lwshen/vault-hub/packages/api/openapi/client`.
- 📌 Confirm cron jobs and any external services are vendoring the updated client or regenerated artefacts.
- 📌 Communicate the generator swap to downstream consumers and decide whether to publish a separate Go module or language-specific SDKs.
- ✅ Docker, Air, and CI workflows invoke the Echo server and official generator scripts.

### 6. Testing, Validation, and Performance
- Run `go test ./...`, `golangci-lint run ./...`, and targeted integration tests to validate Echo handlers and generated clients.
- Add new tests for areas affected by the migration (middleware behavior, error responses, request validation) to mitigate regressions.
- Perform local smoke tests (CLI commands, cron jobs) against the migrated server to verify compatibility.
- Monitor performance differences (latency, throughput) and adjust Echo server configuration (timeouts, workers) if necessary.

### 7. Documentation, Rollout, and Cleanup
- Update README, developer docs, and onboarding materials to describe the Echo server structure and new OpenAPI generation commands.
- Document any new environment variables, configuration flags, or operational considerations introduced by Echo.
- Communicate migration steps to the team, including timelines, testing requirements, and rollback plans; consider a feature branch or staged rollout.
- Remove deprecated Fiber code, tooling, and documentation once Echo is fully validated, ensuring `go.mod` and `go.sum` are tidy.
- Tag a release or create a deployment checklist once QA sign-off is complete, noting the framework change and regenerated clients for downstream consumers.
