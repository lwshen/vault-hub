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
- Evaluate the OpenAPI Generator CLI options (Docker wrapper vs. local CLI) and decide on tooling integration (e.g., add a `tools/openapi-generator-config.yaml` or `packages/api/config.yaml`).
- Stage preview scaffolding (`packages/api/openapi-generator/*.yaml`, `packages/api/generate-openapi.sh`) so contributors can exercise the official generator without disrupting the `oapi-codegen` flow.
- Replace existing `go:generate` directives to invoke `openapi-generator-cli generate` with the Go server + client targets needed by the project.
- Configure generator options for module path, package naming, and Echo integration (e.g., customizing templates or using `--additional-properties` for Echo compatibility).
- Regenerate server stubs and shared models; update build scripts (`go generate`, CI workflows) to use the new generator and remove `oapi-codegen` dependencies.
- Review produced code for lint compliance, add wrapper layers as needed, and commit template overrides if customization is required for Echo or project-specific behavior.

### 5. Update Dependent Applications & Libraries
- Adjust CLI (`apps/cli`, `internal/cli`, `internal/encryption`) to consume the regenerated OpenAPI client or updated REST endpoints.
- Update cron jobs (`apps/cron`) and shared models (`model/`) that rely on Fiber response contracts or previously generated code.
- Refresh `packages/api` exports, ensuring downstream services or SDK consumers know about the new generator outputs and version bump.
- Update Dockerfiles, Air configuration, and release scripts to build/run the Echo server and invoke the new OpenAPI generation workflow.

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
