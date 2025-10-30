# Migration Plan: Fiber to Echo Framework with OpenAPI Generator

## ✅ MIGRATION COMPLETED

**Status**: All phases completed successfully. The project has been migrated from Fiber to Echo framework.

**Date Completed**: Migration executed per plan

---

## Executive Summary

This document outlines the complete migration plan for migrating VaultHub from the Fiber web framework to Echo framework while continuing to use OpenAPI code generation via `oapi-codegen`. The migration maintains API compatibility and improves framework alignment with Echo's standard library approach.

## Current State Analysis

### Framework Stack
- **Web Framework**: Fiber v2.52.9
- **OpenAPI Generator**: oapi-codegen v2.4.1
- **Generation Target**: `fiber-server` (current)
- **Logger Integration**: slog-fiber v1.18.1
- **Session Management**: Fiber session middleware (for OIDC)

### Key Components Using Fiber

1. **Entry Point** (`apps/server/main.go`)
   - Fiber app initialization
   - Slog-fiber middleware integration
   - Route setup

2. **Routing** (`route/route.go`)
   - `fiber.App` usage
   - OpenAPI handler registration via `RegisterHandlers`
   - Static file serving via Fiber filesystem middleware
   - Custom route groups for auth endpoints

3. **Middleware** (`route/middleware.go`)
   - `fiber.Ctx` context usage
   - JWT/API key authentication middleware
   - Route-based authentication rules
   - `c.Locals()` for context data storage

4. **Handlers** (`handler/`)
   - `fiber.Ctx` parameter in all handlers
   - `handler.SendError()` using Fiber response helpers
   - OIDC handlers using Fiber-specific methods (`c.BaseURL()`, `c.Redirect()`, etc.)

5. **API Implementations** (`packages/api/*.go`)
   - All implementations accept `*fiber.Ctx`
   - Utility functions use Fiber context methods
   - Error handling via Fiber response methods

6. **OIDC Integration** (`internal/auth/oidc.go`)
   - Fiber session store
   - `fiber.Ctx` context methods

### OpenAPI Generation Configuration

**Current** (`packages/api/cfg.yaml`):
```yaml
package: api
generate:
  models: true
  fiber-server: true
  strict-server: true
output: generated.go
```

**Generated Code Pattern**:
- `ServerInterface` with `*fiber.Ctx` parameters
- `RegisterHandlers(router fiber.Router, si ServerInterface)`
- `ServerInterfaceWrapper` uses Fiber context methods

## Migration Goals

1. ✅ **Zero API Breaking Changes**: Maintain 100% API compatibility
2. ✅ **Framework Migration**: Replace Fiber with Echo v4.x
3. ✅ **OpenAPI Regeneration**: Update codegen config to Echo-compatible generation
4. ✅ **Feature Parity**: Maintain all existing features (auth, middleware, static files)
5. ✅ **Performance**: Leverage Echo's performance optimizations
6. ✅ **Standard Library Alignment**: Better alignment with `net/http`

## Migration Strategy

### Phase 1: Preparation & Dependency Updates

#### 1.1 Update Dependencies
- **Add Echo**: `github.com/labstack/echo/v4`
- **Update oapi-codegen config**: Change from `fiber-server` to `echo-server`
- **Add Echo middlewares**: 
  - `github.com/labstack/echo/v4/middleware` (logger, recover, etc.)
  - Replace `slog-fiber` with Echo's native logging or custom middleware
- **Session Management**: Migrate from Fiber sessions to Echo-compatible session library
  - Options: `github.com/gorilla/sessions` with Echo adapter or `github.com/labstack/echo-contrib/session`

#### 1.2 Update OpenAPI Generator Config

**File**: `packages/api/cfg.yaml`
```yaml
package: api
generate:
  models: true
  echo-server: true        # Changed from fiber-server
  strict-server: true
output: generated.go
```

#### 1.3 Regenerate OpenAPI Code
- Run `go generate packages/api/tool.go`
- This will generate Echo-compatible server interfaces:
  - `ServerInterface` with `echo.Context` parameters
  - `RegisterHandlers(e *echo.Echo, si ServerInterface)`
  - Context parameter types changed from `*fiber.Ctx` to `echo.Context`

### Phase 2: Core Framework Migration

#### 2.1 Entry Point Migration (`apps/server/main.go`)

**Before (Fiber)**:
```go
app := fiber.New()
app.Use(slogfiber.New(logger))
route.SetupRoutes(app)
log.Fatal(app.Listen(":" + config.AppPort))
```

**After (Echo)**:
```go
e := echo.New()
e.Use(middleware.Logger())
e.Use(middleware.Recover())
// Custom slog middleware for structured logging
e.Use(SlogMiddleware(logger))
route.SetupRoutes(e)
log.Fatal(e.Start(":" + config.AppPort))
```

**Changes Required**:
- Replace `fiber.New()` with `echo.New()`
- Replace `slogfiber.New()` with custom Echo middleware wrapper
- Replace `app.Listen()` with `e.Start()`
- Add Echo middleware (logger, recover)

#### 2.2 Route Setup Migration (`route/route.go`)

**Before (Fiber)**:
```go
func SetupRoutes(app *fiber.App) {
    app.Use(jwtMiddleware)
    server := openapi.NewServer()
    openapi.RegisterHandlers(app, server)
    // ...
}
```

**After (Echo)**:
```go
func SetupRoutes(e *echo.Echo) {
    e.Use(jwtMiddleware)
    server := openapi.NewServer()
    openapi.RegisterHandlers(e, server)
    // ...
}
```

**Changes Required**:
- Change function signature: `*fiber.App` → `*echo.Echo`
- OpenAPI `RegisterHandlers` will accept `*echo.Echo` (after regeneration)
- Static file serving: Replace Fiber filesystem middleware with Echo static handler
- Route groups: Replace `app.Group()` with `e.Group()`

**Static File Serving Migration**:
```go
// Before (Fiber)
app.Use("/", filesystem.New(filesystem.Config{
    Root:         http.FS(embedFS),
    Index:        "index.html",
    NotFoundFile: "index.html",
}))

// After (Echo)
e.StaticFS("/", embedFS)
e.GET("/*", func(c echo.Context) error {
    return c.File("index.html")
}, NotFoundMiddleware)
```

#### 2.3 Middleware Migration (`route/middleware.go`)

**Key Changes**:

1. **Context Type**: `*fiber.Ctx` → `echo.Context`

2. **Request/Response Methods**:
   - `c.Path()` → `c.Path()` (same method name, compatible)
   - `c.Get("Header")` → `c.Request().Header.Get("Header")`
   - `c.Query("param")` → `c.QueryParam("param")`
   - `c.SendStatus(code)` → `c.NoContent(code)` or `c.JSON(code, nil)`
   - `c.Status(code).JSON(data)` → `c.JSON(code, data)`
   - `c.IP()` → `c.RealIP()` or custom extraction
   - `c.Locals("key")` → `c.Set("key", value)` / `c.Get("key")`

3. **Error Handling**:
   - `return c.Next()` → `return nil` (Echo returns nil on success)
   - Fiber returns error, Echo returns error (compatible pattern)

4. **Response Helpers**:
   - `c.SendStatus()` → `c.NoContent()` or `c.JSON()`
   - `c.Redirect()` → `c.Redirect()` (same signature)

**Middleware Function Signature**:
```go
// Before (Fiber)
func jwtMiddleware(c *fiber.Ctx) error {
    // ...
    return c.Next()
}

// After (Echo)
func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // ...
        return next(c)
    }
}
```

**Middleware Registration**:
```go
// Before (Fiber)
app.Use(jwtMiddleware)

// After (Echo)
e.Use(jwtMiddleware)  // Works if middleware returns echo.MiddlewareFunc
// OR
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return jwtMiddleware
})
```

#### 2.4 Handler Response Migration (`handler/response.go`)

**Before (Fiber)**:
```go
func SendError(c *fiber.Ctx, code int, message string) error {
    return c.Status(code).JSON(fiber.Map{
        "error": ErrorResponse{
            Code:    code,
            Message: message,
        },
    })
}
```

**After (Echo)**:
```go
func SendError(c echo.Context, code int, message string) error {
    return c.JSON(code, map[string]interface{}{
        "error": ErrorResponse{
            Code:    code,
            Message: message,
        },
    })
}
```

**Changes**:
- Remove `fiber.Map` dependency (use `map[string]interface{}`)
- `c.Status(code).JSON()` → `c.JSON(code, ...)`
- Context type change

#### 2.5 OIDC Handler Migration (`handler/auth.go`)

**Key Method Translations**:

| Fiber Method | Echo Equivalent |
|-------------|----------------|
| `c.BaseURL()` | `c.Scheme() + "://" + c.Request().Host` |
| `c.Redirect(url)` | `c.Redirect(code, url)` (requires status code) |
| `c.Query("param")` | `c.QueryParam("param")` |
| `c.SendStatus(code)` | `c.NoContent(code)` |
| `c.Request().URI()` | `c.Request().URL` |
| `c.Context()` | `c.Request().Context()` |

**Example Migration**:
```go
// Before (Fiber)
func LoginOidc(c *fiber.Ctx) error {
    baseUrl := c.BaseURL()
    url, err := auth.AuthCodeURL(c, baseUrl)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }
    return c.Redirect(url)
}

// After (Echo)
func LoginOidc(c echo.Context) error {
    scheme := c.Scheme()
    if c.Request().Header.Get("X-Forwarded-Proto") != "" {
        scheme = c.Request().Header.Get("X-Forwarded-Proto")
    }
    baseUrl := scheme + "://" + c.Request().Host
    url, err := auth.AuthCodeURL(c.Request().Context(), baseUrl)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }
    return c.Redirect(http.StatusFound, url)
}
```

#### 2.6 OIDC Session Migration (`internal/auth/oidc.go`)

**Current**: Uses Fiber session middleware
```go
sessionStore *session.Store  // Fiber session
```

**Migration Options**:

**Option A: Gorilla Sessions with Echo Adapter**
```go
import (
    "github.com/gorilla/sessions"
    "github.com/labstack/echo-contrib/session"
)

sessionStore := sessions.NewCookieStore([]byte(config.SessionSecret))
e.Use(session.Middleware(sessionStore))
```

**Option B: Echo Contrib Session**
```go
import "github.com/labstack/echo-contrib/session"

e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))
```

**Session Usage Changes**:
```go
// Before (Fiber)
sess, _ := sessionStore.Get(c)
sess.Set("key", value)
sess.Save()

// After (Echo)
sess, _ := session.Get("session", c)
sess.Values["key"] = value
sess.Save(c.Request(), c.Response())
```

### Phase 3: API Implementation Migration

#### 3.1 Update All API Implementations

**Files to Update**:
- `packages/api/vault.go`
- `packages/api/auth.go`
- `packages/api/user.go`
- `packages/api/api_key.go`
- `packages/api/audit_log.go`
- `packages/api/status.go`
- `packages/api/config.go`
- `packages/api/health.go`
- `packages/api/cli_vault.go`

**Change Pattern**:
```go
// Before (Fiber)
func (s *Server) GetVaults(c *fiber.Ctx) error {
    user := c.Locals("user").(*model.User)
    // ...
    return c.JSON(fiber.StatusOK, response)
}

// After (Echo)
func (s *Server) GetVaults(c echo.Context) error {
    user := c.Get("user").(*model.User)
    // ...
    return c.JSON(http.StatusOK, response)
}
```

**Common Changes**:
- `*fiber.Ctx` → `echo.Context`
- `c.Locals("key")` → `c.Get("key")`
- `c.Locals("key", value)` → `c.Set("key", value)`
- `fiber.StatusXXX` → `http.StatusXXX`
- `c.Query("param")` → `c.QueryParam("param")`
- `c.Get("Header")` → `c.Request().Header.Get("Header")`
- `c.IP()` → `c.RealIP()` or custom extraction
- Remove `fiber.Map` → use `map[string]interface{}`

#### 3.2 Utility Function Updates

**Example**: `getClientInfo` function
```go
// Before (Fiber)
func getClientInfo(c *fiber.Ctx) (string, string) {
    ip := c.Get("X-Forwarded-For")
    if ip == "" {
        ip = c.Get("X-Real-IP")
    }
    if ip == "" {
        ip = c.IP()
    }
    userAgent := c.Get("User-Agent")
    return ip, userAgent
}

// After (Echo)
func getClientInfo(c echo.Context) (string, string) {
    ip := c.Request().Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = c.Request().Header.Get("X-Real-IP")
    }
    if ip == "" {
        ip = c.RealIP()
    }
    userAgent := c.Request().Header.Get("User-Agent")
    return ip, userAgent
}
```

### Phase 4: Testing & Validation

#### 4.1 Unit Tests
- Update test mocks/stubs to use Echo context
- Update middleware tests
- Update handler tests

#### 4.2 Integration Tests
- Test all API endpoints
- Verify authentication flows (JWT, API key, OIDC)
- Test static file serving
- Test error handling

#### 4.3 Manual Testing Checklist
- [ ] Server starts successfully
- [ ] All API endpoints respond correctly
- [ ] JWT authentication works
- [ ] API key authentication works
- [ ] OIDC login flow works
- [ ] Static frontend assets load correctly
- [ ] Error responses are properly formatted
- [ ] Audit logging captures correct IP/user agent
- [ ] Middleware executes in correct order
- [ ] CORS headers (if any) are preserved

### Phase 5: Documentation & Cleanup

#### 5.1 Update Documentation
- Update `CLAUDE.md` with Echo-specific commands
- Update API documentation if needed
- Update deployment docs
- Update development setup instructions

#### 5.2 Code Cleanup
- Remove unused Fiber imports
- Remove `slog-fiber` dependency
- Clean up commented code
- Verify no Fiber references remain

## Detailed File-by-File Migration Checklist

### Core Framework Files

- [ ] `apps/server/main.go`
  - Replace `fiber.New()` with `echo.New()`
  - Replace `slogfiber.New()` with custom Echo middleware
  - Replace `app.Listen()` with `e.Start()`
  - Add Echo middleware (logger, recover)

- [ ] `route/route.go`
  - Change function signature to accept `*echo.Echo`
  - Update `RegisterHandlers` call (after regeneration)
  - Migrate static file serving
  - Update route group creation

- [ ] `route/middleware.go`
  - Change all middleware signatures to `echo.MiddlewareFunc`
  - Update context access methods
  - Update `c.Locals()` → `c.Set()` / `c.Get()`
  - Update error responses
  - Update `isPublicRoute()` helper (if needed)

### Handler Files

- [ ] `handler/response.go`
  - Change `SendError` signature to use `echo.Context`
  - Remove `fiber.Map` dependency

- [ ] `handler/auth.go`
  - Update `LoginOidc` function
  - Update `LoginOidcCallback` function
  - Update `getClientInfo` helper
  - Change redirect methods

### API Implementation Files

- [ ] `packages/api/impl.go`
  - Update `Server` struct if needed
  - Verify `ServerInterface` conformance (after regeneration)

- [ ] `packages/api/vault.go`
  - Update all handler methods to use `echo.Context`
  - Update context access
  - Update error handling

- [ ] `packages/api/auth.go`
  - Update authentication handlers
  - Update context access

- [ ] `packages/api/user.go`
  - Update user management handlers

- [ ] `packages/api/api_key.go`
  - Update API key handlers

- [ ] `packages/api/audit_log.go`
  - Update audit log handlers

- [ ] `packages/api/status.go`
  - Update status endpoint

- [ ] `packages/api/config.go`
  - Update config endpoint

- [ ] `packages/api/health.go`
  - Update health endpoint

- [ ] `packages/api/cli_vault.go`
  - Update CLI vault handlers

### Authentication & OIDC

- [ ] `internal/auth/oidc.go`
  - Replace Fiber session store with Echo-compatible sessions
  - Update session methods
  - Update context usage

### Configuration Files

- [ ] `packages/api/cfg.yaml`
  - Change `fiber-server` to `echo-server`

- [ ] `go.mod`
  - Remove `github.com/gofiber/fiber/v2`
  - Remove `github.com/samber/slog-fiber`
  - Add `github.com/labstack/echo/v4`
  - Add `github.com/labstack/echo-contrib/session` (or gorilla sessions)
  - Run `go mod tidy`

- [ ] `go.sum`
  - Will be updated automatically with `go mod tidy`

## Dependencies Update Summary

### Remove
```go
github.com/gofiber/fiber/v2 v2.52.9
github.com/samber/slog-fiber v1.18.1
github.com/gofiber/fiber/v2/middleware/filesystem
github.com/gofiber/fiber/v2/middleware/session
```

### Add
```go
github.com/labstack/echo/v4 v4.12.0
github.com/labstack/echo-contrib/session v0.0.0-...
// OR
github.com/gorilla/sessions v1.2.2
```

### Keep (Already Compatible)
```go
github.com/oapi-codegen/oapi-codegen/v2 v2.4.1
github.com/oapi-codegen/runtime v1.1.2
github.com/golang-jwt/jwt/v5 v5.3.0
// ... all other dependencies remain the same
```

## Testing Strategy

### 1. Unit Tests
```bash
# Update test files to use Echo context mocks
JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...
```

### 2. Integration Tests
- Test all endpoints via HTTP client
- Verify OpenAPI spec compliance
- Test authentication flows

### 3. Manual Verification
```bash
# Start server
go run ./apps/server/main.go

# Test endpoints
curl http://localhost:3000/api/health
curl http://localhost:3000/api/status
```

## Risk Assessment & Mitigation

### High Risk Areas

1. **Middleware Chain Execution**
   - **Risk**: Echo middleware execution order differs from Fiber
   - **Mitigation**: Test middleware chain thoroughly, verify order

2. **Context Data Access**
   - **Risk**: `c.Locals()` vs `c.Set()`/`c.Get()` might have different behavior
   - **Mitigation**: Search and replace systematically, test all context usage

3. **Error Handling**
   - **Risk**: Error propagation might differ
   - **Mitigation**: Echo uses standard error return pattern (compatible)

4. **Static File Serving**
   - **Risk**: SPA routing might break
   - **Mitigation**: Test frontend navigation thoroughly

5. **OIDC Session Management**
   - **Risk**: Session storage format might differ
   - **Mitigation**: Test OIDC flow end-to-end

### Medium Risk Areas

1. **IP Address Extraction**
   - Echo's `RealIP()` might handle proxies differently
   - Test with various proxy configurations

2. **Request Context**
   - Context passing for OIDC/cancellation
   - Verify `c.Request().Context()` usage

## Rollback Plan

1. **Git Branch Strategy**
   - Create feature branch: `migrate/echo-framework`
   - Keep `main` branch untouched until migration is validated
   - Can revert branch if issues arise

2. **Incremental Migration**
   - Migrate one module at a time
   - Test after each module migration
   - Use feature flags if needed

3. **Dependency Rollback**
   - Keep `go.mod` history
   - Can revert to Fiber if critical issues found

## Timeline Estimate

### Phase 1: Preparation (1-2 hours)
- Update dependencies
- Regenerate OpenAPI code
- Initial testing

### Phase 2: Core Migration (4-6 hours)
- Entry point migration
- Route setup migration
- Middleware migration
- Handler migration

### Phase 3: API Implementation (6-8 hours)
- Update all API implementation files
- Update utility functions
- OIDC session migration

### Phase 4: Testing (4-6 hours)
- Unit test updates
- Integration testing
- Manual verification
- Bug fixes

### Phase 5: Documentation (1-2 hours)
- Update documentation
- Code cleanup

**Total Estimated Time**: 16-24 hours

## Success Criteria

1. ✅ All tests pass
2. ✅ Server starts without errors
3. ✅ All API endpoints respond correctly
4. ✅ Authentication flows work (JWT, API key, OIDC)
5. ✅ Frontend loads and functions correctly
6. ✅ No Fiber dependencies remain
7. ✅ Code passes linter checks
8. ✅ Performance is maintained or improved
9. ✅ Documentation is updated

## Post-Migration Improvements

After successful migration, consider:

1. **Echo-Specific Features**
   - Leverage Echo's middleware ecosystem
   - Use Echo's validation features
   - Consider Echo's CORS middleware

2. **Performance Optimization**
   - Benchmark and compare performance
   - Optimize middleware chain
   - Consider Echo's HTTP/2 support

3. **Code Quality**
   - Standardize error handling with Echo's error handler
   - Use Echo's binding features for request validation
   - Leverage Echo's context methods for cleaner code

## Notes

- **OpenAPI Compatibility**: The OpenAPI spec itself doesn't change, only the generated server code
- **API Compatibility**: All API endpoints remain identical, only the framework implementation changes
- **Breaking Changes**: None to external API consumers
- **Client Libraries**: No changes needed to TypeScript/Go client libraries

## Reference Links

- [Echo Framework Documentation](https://echo.labstack.com/)
- [oapi-codegen Echo Server Generation](https://github.com/oapi-codegen/oapi-codegen/blob/master/pkg/codegen/templates/echo-server.tmpl)
- [Echo Middleware Guide](https://echo.labstack.com/middleware/)
- [Echo Context API](https://pkg.go.dev/github.com/labstack/echo/v4#Context)
