# Migration Plan: Fiber to Echo & Official OpenAPI Generator

## Purpose
- Standardize the HTTP layer on Echo to align with organizational Go service conventions and simplify middleware reuse.
- Replace `oapi-codegen` with the official OpenAPI Generator tooling to unlock broader language targets, consistent code style, and automated server/client parity.
- Reduce maintenance overhead by converging on supported frameworks, improving onboarding, and tightening contract-first development around the OpenAPI spec.

## Overall Process
1. **Discovery & Design** – Inventory Fiber-specific patterns, catalog generated API bindings, and design the target Echo-based architecture plus generation workflow.
2. **OpenAPI Pipeline Migration** – Introduce the official generator, reproduce existing artifacts, and adapt the codebase to the new output shape while keeping the spec untouched initially.
3. **Server Framework Migration** – Incrementally port middleware, routing, handlers, and integration points from Fiber to Echo, ensuring feature parity and updated tests.
4. **Stabilization & Cleanup** – Run full QA (tests, lint, manual smoke), update docs/tooling, and remove dead Fiber/oapi-codegen assets.

## Detailed Steps

### Phase 1 – Discovery & Design
- Audit current HTTP stack: middleware (`route/middleware.go`), routing (`route/route.go`), and handlers (`handler/*.go`) to understand Fiber-specific usage (context helpers, groups, responses).
- Document cross-cutting concerns (logging via `slog-fiber`, JWT middleware, static asset serving, request context usage) and the equivalent Echo strategy.
- Review `packages/api` generation flow (`tool.go`, `cfg.yaml`, `generated.go`, `impl.go`) to capture the current contract-to-implementation handoff and any custom patches on the generated code.
- Decide on the official generator targets:
  - Confirm the `openapi-generator-cli` version and install path (likely vendored via `go:generate` or `tools.go`).
  - Select generators (`go-echo-server` or `go-server` with Echo library, `typescript-fetch` if web consumers need regeneration).
- Produce a migration design doc snippet (module layout, package boundaries) to validate with stakeholders before coding.

### Phase 2 – OpenAPI Pipeline Migration
- Add tooling dependencies: vendor/download `openapi-generator-cli` (Docker or jar) and wire `go:generate` or `Makefile` helpers paralleling the existing `bundle.sh`; do not install it globally (avoid `npm install -g`), prefer `npx` or a checked-in wrapper script.
- Translate `cfg.yaml` settings to the official generator config (`.json` or `.yaml`) covering package names, enum handling, nullable types, and interface generation.
- Generate server stubs and models into a staging package (e.g., `packages/api/gen`) without deleting existing `oapi-codegen` outputs; compare shape and identify required adapter layers.
- Update manual implementation files (`packages/api/impl.go`, audit log wiring, etc.) to satisfy the new generated interface signatures.
- Regenerate client/SDK consumers (CLI, cron jobs) if they rely on the previous generated Go types; refactor imports accordingly.
- Adjust unit tests to account for any struct/tag differences; ensure backward compatibility for serialized JSON where possible.

### Phase 3 – Server Framework Migration
- Introduce Echo dependencies in `go.mod` and add initial bootstrap in `apps/server/main.go` (Echo instance, logger integration, graceful shutdown).
- Replace Fiber middleware with Echo equivalents:
  - Logging: swap `slog-fiber` for Echo-compatible middleware (custom or community) and wrap `slog` manually if needed.
  - Authentication/JWT: port `jwtMiddleware` logic to Echo middleware chain, ensuring request context semantics are preserved.
  - Static assets: reproduce SPA hosting via Echo’s `StaticFS` with the embedded dist filesystem.
- Refactor routing:
  - Update `route/` to construct Echo groups (`/api`, `/api/auth`, etc.) and register generated OpenAPI handlers using the new generator’s router glue.
  - Migrate inline closures (e.g., `/magic-link/token`) to Echo handler signatures (`echo.Context`), converting response helpers (`c.SendStatus`, `c.JSON`, etc.) to Echo equivalents.
- Update `handler/` layer to consume `echo.Context` (body/query/headers access, context propagation) and revise helper utilities (`getClientInfo`) for the new API.
- Ensure middlewares and other packages using Fiber types (e.g., request context in `internal/auth`) are updated to work with Echo’s context/Request objects.
- Run `go test ./...` and resolve compilation/test failures until parity is achieved.

### Phase 4 – Stabilization & Cleanup
- Remove Fiber-specific dependencies (`github.com/gofiber/*`, `slog-fiber`) and old generated artifacts; tidy modules with `go mod tidy`.
- Update documentation (`README.md`, internal runbooks) to reflect Echo commands, middleware guidance, and the new OpenAPI generation workflow.
- Refresh CI scripts and lint configs if they reference Fiber or the old generator.
- Perform manual smoke tests: OIDC login redirect, magic link consumption, static frontend serving, CLI interactions against the new server.
- Capture migration notes, highlight breaking changes (if any), and align rollout steps (deploy order, toggles, rollback plan).
