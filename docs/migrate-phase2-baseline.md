# Phase 2 – Echo Baseline

## Overview
Phase 2 established the Echo server baseline that ultimately replaced the Fiber implementation. The build tag toggle has since been removed—Echo is now the default runtime—but the key architectural decisions remain relevant.

## Key Additions
- Dependency: `github.com/labstack/echo/v4 v4.12.0` is now a first-class requirement in `go.mod`.
- Echo bootstrap: `internal/server/echoapp/server.go` exposes `NewServer` to build a preconfigured Echo instance (logging, recovery, security, CORS, static assets).
- Security middleware: `internal/server/echoapp/security.go` handles JWT/API key enforcement and context population.
- Error helper: `internal/server/echoapp/response.go` keeps error payloads consistent with the legacy JSON envelope.
- Server entrypoint: `apps/server/main.go` calls the Echo bootstrap directly and manages graceful shutdown on SIGINT/SIGTERM.

## Running the Echo Baseline
1. Ensure dependencies are installed (`go mod tidy` is part of the standard workflow).
2. Run the server without build tags:
   ```bash
   go run ./apps/server
   ```
3. Embedded frontend assets are served automatically; all API routes are backed by Echo handlers.

## Next Steps
- (Completed in Phase 3) Wire OpenAPI-driven helpers into Echo routes and retire Fiber shims.
- (Completed in Phase 3) Replace Fiber-specific helpers with framework-agnostic equivalents in `packages/api`.
- (Completed in Phase 3) Swap Fiber session middleware for an Echo-friendly OIDC state cache.
