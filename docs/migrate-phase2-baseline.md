# Phase 2 – Echo Baseline

## Overview
Phase 2 introduces an initial Echo server implementation that can run alongside the existing Fiber app by opting into the `echo` build tag. This baseline includes request logging, security middleware, static asset mounting, and graceful shutdown wiring so future phases can focus on porting routes and handlers.

## Key Additions
- Dependency: `github.com/labstack/echo/v4 v4.12.0` added in `go.mod`.
- Echo bootstrap: `internal/server/echoapp/server.go` exposes `NewServer` to build a preconfigured Echo instance (logging, recovery, security, CORS, static assets).
- Security middleware: `internal/server/echoapp/security.go` translates JWT/API key enforcement and context population to Echo.
- Error helper: `internal/server/echoapp/response.go` keeps error payloads consistent with the Fiber implementation.
- Alternate entrypoint: `apps/server/main_echo.go` (build tag `echo`) starts the Echo server and performs graceful shutdown on SIGINT/SIGTERM.

## Running the Echo Baseline
1. Ensure dependencies are installed (`go mod tidy` already run as part of this change).
2. Use the `echo` build tag when running the server:
   ```bash
   go run -tags echo ./apps/server
   ```
3. Requests will still respond with the embedded frontend; API routes currently return 404 until handlers are ported in later phases.

## Next Steps
- Wire OpenAPI-generated handlers into the Echo router once Phase 3 migration of routes begins.
- Replace Fiber-specific helpers consumed by handler packages with Echo versions, reusing the middleware state stored on the context.
- Evaluate session replacements (currently still tied to Fiber middleware) before porting OIDC flows.
