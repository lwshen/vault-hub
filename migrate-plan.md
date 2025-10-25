# Migration Plan: Fiber to Echo Framework

**Document Version:** 1.0
**Date:** 2025-01-25
**Project:** VaultHub
**Migration Type:** Web Framework Change (Fiber → Echo)

---

## Table of Contents

1. [Purpose](#purpose)
2. [Executive Summary](#executive-summary)
3. [Overall Process](#overall-process)
4. [Detailed Migration Steps](#detailed-migration-steps)
5. [Risk Assessment](#risk-assessment)
6. [Testing Strategy](#testing-strategy)
7. [Rollback Plan](#rollback-plan)

---

## Purpose

### Why Migrate?

The current implementation uses **Fiber** web framework with **oapi-codegen** for OpenAPI code generation. A critical bug exists in oapi-codegen's Fiber server template (Issue #1806) where header parameters generate incorrect code:

```go
// Bug: c.GetReqHeaders() returns map[string][]string
// but runtime.BindStyledParameterWithOptions expects string
headers := c.GetReqHeaders()
err = runtime.BindStyledParameterWithOptions("simple", "X-Enable-Client-Encryption", value, ...)
// Error: cannot use value (variable of type []string) as string value
```

### Goals

1. ✅ **Eliminate the oapi-codegen header bug** (no workaround needed)
2. ✅ **Maintain OpenAPI-first development workflow** (keep current spec)
3. ✅ **Preserve all existing features** (encryption, auth, audit logs)
4. ✅ **Keep same tooling** (oapi-codegen, no generator switch)
5. ✅ **Gain HTTP/2 support** and better ecosystem compatibility
6. ✅ **Minimal disruption** to API consumers (no breaking changes)

### Why Echo?

| Criteria | Fiber | Echo |
|----------|-------|------|
| **OpenAPI Support** | ⚠️ Buggy (oapi-codegen #1806) | ✅ Stable |
| **Standard Library** | ❌ Uses fasthttp | ✅ Uses net/http |
| **HTTP/2** | ❌ No | ✅ Yes |
| **Performance** | 36k RPS | 34k RPS (~6% slower) |
| **Ecosystem** | Limited (fasthttp-specific) | Excellent (net/http compatible) |
| **Code Generator** | oapi-codegen (buggy) | oapi-codegen (stable) |
| **Production Ready** | ✅ Yes | ✅ Yes |

**Decision:** Migrate to Echo using oapi-codegen (NOT OpenAPI Generator, which lacks auth support).

---

## Executive Summary

### Scope

- **Files to Modify:** ~16 Go files (436 Fiber references)
- **Lines of Code:** Estimated 500-800 LOC changes
- **Database Changes:** None (purely code change)
- **API Contract Changes:** None (OpenAPI spec unchanged)
- **Client Impact:** None (TypeScript client auto-regenerates)

### Timeline

| Phase | Duration | Effort Level |
|-------|----------|--------------|
| 1. Preparation | 0.5 days | Low |
| 2. Code Generation | 0.5 days | Low |
| 3. Handler Migration | 1-2 days | High |
| 4. Middleware Migration | 1 day | Medium |
| 5. Testing | 1 day | Medium |
| 6. Documentation | 0.5 days | Low |
| **TOTAL** | **4-6 days** | **Medium** |

### Risk Level: **MEDIUM**

- ✅ No database migrations
- ✅ No external API changes
- ⚠️ Significant code changes across handlers
- ⚠️ Requires thorough testing of authentication flows

---

## Overall Process

```
Current State: Fiber
    ↓
Create Branch
    ↓
Update Dependencies
    ↓
Change oapi-codegen Config (fiber-server → echo-server)
    ↓
Regenerate Code
    ↓
Fix Compilation Errors
    ↓
Migrate Handlers (14 files)
    ↓
Migrate Middleware
    ↓
Update Tests
    ↓
Integration Testing
    ↓
Tests Pass? → No → Fix Issues
    ↓ Yes
Code Review
    ↓
Deploy to Staging
    ↓
Staging OK? → No → Rollback
    ↓ Yes
Deploy to Production
    ↓
Monitor
```

### Key Principles

1. **Incremental Changes**: Fix one component at a time
2. **Test Early, Test Often**: Run tests after each phase
3. **Backward Compatibility**: Keep OpenAPI spec unchanged
4. **Rollback Ready**: Use feature branch, easy to revert

---

## Detailed Migration Steps

### Phase 1: Preparation (0.5 days)

#### 1.1 Create Migration Branch
```bash
git checkout -b migrate-fiber-to-echo
```

#### 1.2 Update Dependencies

**File:** `go.mod`

**Add:**
```go
github.com/labstack/echo/v4 v4.13.3
```

**Remove (if not used elsewhere):**
```go
github.com/gofiber/fiber/v2 v2.52.9
```

**Run:**
```bash
go mod tidy
```

#### 1.3 Update oapi-codegen Configuration

**File:** `packages/api/cfg.yaml`

**Before:**
```yaml
package: api
generate:
  models: true
  fiber-server: true
  strict-server: true
output: generated.go
```

**After:**
```yaml
package: api
generate:
  models: true
  echo-server: true
  strict-server: true
output: generated.go
```

---

### Phase 2: Code Generation (0.5 days)

#### 2.1 Regenerate Server Code

```bash
cd packages/api
go generate tool.go
```

**Expected Output:**
- `generated.go` now uses Echo types instead of Fiber
- No compilation errors related to header parameters
- Server interface methods use `echo.Context`

#### 2.2 Verify Generated Code

**Check:**
```bash
grep -n "echo.Context" packages/api/generated.go | head -5
```

**Expected:** Handler signatures like:
```go
type ServerInterface interface {
    GetVaultByAPIKey(ctx echo.Context, uniqueId string, params GetVaultByAPIKeyParams) error
    // ...
}
```

#### 2.3 Initial Compilation Check

```bash
go build -o /dev/null ./apps/server/main.go
```

**Expected:** Compilation errors (we'll fix these next).

---

### Phase 3: Handler Migration (1-2 days)

#### 3.1 Update Handler Signatures

**Files to Update:** (14 files in `packages/api/`)
- `api_key.go`
- `audit_log.go`
- `auth.go`
- `cli_vault.go`
- `config.go`
- `health.go`
- `status.go`
- `user.go`
- `vault.go`
- Plus 5 more handler files

**Migration Pattern:**

| Fiber | Echo |
|-------|------|
| `func Handler(c *fiber.Ctx, ...) error` | `func Handler(ctx echo.Context, ...) error` |
| `c.JSON(data)` | `ctx.JSON(http.StatusOK, data)` |
| `c.Status(code).JSON(data)` | `ctx.JSON(code, data)` |
| `c.Locals("key")` | `ctx.Get("key")` |
| `c.Locals("key", value)` | `ctx.Set("key", value)` |
| `c.Params("id")` | `ctx.Param("id")` |
| `c.Query("filter")` | `ctx.QueryParam("filter")` |
| `c.Get("Header-Name")` | `ctx.Request().Header.Get("Header-Name")` |
| `c.BodyParser(&data)` | `ctx.Bind(&data)` |
| `fiber.NewError(code, msg)` | `echo.NewHTTPError(code, msg)` |

#### 3.2 Example: Vault Handler Migration

**File:** `packages/api/vault.go`

**Before (Fiber):**
```go
func (s *StrictServer) GetVaultByAPIKey(c *fiber.Ctx, uniqueId string, params GetVaultByAPIKeyParams) error {
    // Get user from context
    userID := c.Locals("user_id").(*uint)
    apiKey := c.Locals("api_key").(*model.APIKey)

    // Check encryption header
    enableEncryption := c.Get("X-Enable-Client-Encryption") == "true"

    // Business logic
    vault, err := s.vaultService.GetByUniqueID(uniqueId)
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "Vault not found")
    }

    // Return response
    return c.JSON(vault)
}
```

**After (Echo):**
```go
func (s *StrictServer) GetVaultByAPIKey(ctx echo.Context, uniqueId string, params GetVaultByAPIKeyParams) error {
    // Get user from context
    userID := ctx.Get("user_id").(*uint)
    apiKey := ctx.Get("api_key").(*model.APIKey)

    // Check encryption header
    enableEncryption := ctx.Request().Header.Get("X-Enable-Client-Encryption") == "true"

    // Business logic
    vault, err := s.vaultService.GetByUniqueID(uniqueId)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Vault not found")
    }

    // Return response
    return ctx.JSON(http.StatusOK, vault)
}
```

#### 3.3 Update Response Helpers

**File:** `handler/response.go`

**Before:**
```go
func SendError(c *fiber.Ctx, code int, message string) error {
    return c.Status(code).JSON(fiber.Map{"error": message})
}
```

**After:**
```go
func SendError(ctx echo.Context, code int, message string) error {
    return ctx.JSON(code, map[string]string{"error": message})
}
```

---

### Phase 4: Middleware Migration (1 day)

#### 4.1 Update JWT Middleware

**File:** `route/middleware.go`

**Before (Fiber):**
```go
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization")
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := auth.ValidateJWT(token)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        c.Locals("user", claims.User)
        return c.Next()
    }
}
```

**After (Echo):**
```go
func JWTMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            authHeader := ctx.Request().Header.Get("Authorization")
            if authHeader == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization")
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := auth.ValidateJWT(token)
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
            }

            ctx.Set("user", claims.User)
            return next(ctx)
        }
    }
}
```

#### 4.2 Update API Key Middleware

**Similar pattern to JWT middleware** - update signatures and context methods.

#### 4.3 Update CORS Middleware

**Option 1: Use Echo's Built-in CORS**
```go
import "github.com/labstack/echo/v4/middleware"

e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
}))
```

---

### Phase 5: Main Server Setup (0.5 days)

#### 5.1 Update Server Initialization

**File:** `apps/server/main.go`

**Before (Fiber):**
```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/lwshen/vault-hub/route"
)

func main() {
    app := fiber.New()

    // Setup routes
    route.Setup(app)

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

**After (Echo):**
```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/lwshen/vault-hub/route"
)

func main() {
    e := echo.New()

    // Built-in middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Setup routes
    route.Setup(e)

    // Start server
    e.Logger.Fatal(e.Start(":3000"))
}
```

#### 5.2 Update Route Registration

**File:** `route/route.go`

**Before (Fiber):**
```go
func Setup(app *fiber.App) {
    api.RegisterHandlers(app, handler)
}
```

**After (Echo):**
```go
func Setup(e *echo.Echo) {
    api.RegisterHandlers(e, handler)
}
```

---

### Phase 6: Static File Serving (0.5 days)

#### 6.1 Update Static Assets Handler

**Before (Fiber):**
```go
app.Static("/", "./internal/embed/dist")
```

**After (Echo):**
```go
e.Static("/", "internal/embed/dist")
```

---

### Phase 7: Testing (1 day)

#### 7.1 Unit Tests

**Update test files** that use Fiber test helpers:

**Before:**
```go
import "github.com/gofiber/fiber/v2"

app := fiber.New()
req := httptest.NewRequest("GET", "/api/vaults", nil)
resp, _ := app.Test(req)
```

**After:**
```go
import (
    "github.com/labstack/echo/v4"
    "net/http/httptest"
)

e := echo.New()
req := httptest.NewRequest(http.MethodGet, "/api/vaults", nil)
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)
```

#### 7.2 Run Test Suite

```bash
JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...
```

**Expected:** All tests pass.

#### 7.3 Integration Testing Checklist

- [ ] User registration and login
- [ ] JWT token generation and validation
- [ ] API key authentication
- [ ] Vault CRUD operations
- [ ] Client-side encryption (`X-Enable-Client-Encryption` header)
- [ ] Audit log creation
- [ ] OIDC authentication (if enabled)
- [ ] Static file serving (frontend)
- [ ] CORS functionality
- [ ] Error responses (4xx, 5xx)

---

### Phase 8: Documentation Updates (0.5 days)

#### 8.1 Update CLAUDE.md

**Update framework references:**
```markdown
# Old
- **Framework**: Fiber v2.52.9

# New
- **Framework**: Echo v4.13.3
```

**Update development commands:**
```markdown
# Old
Uses Fiber with fasthttp for high performance

# New
Uses Echo with standard net/http and HTTP/2 support
```

#### 8.2 Update README.md

Update any Fiber-specific documentation.

#### 8.3 Update API Documentation

If you have developer documentation referencing Fiber, update it.

---

## Risk Assessment

### High-Risk Areas

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Authentication Breaks** | Critical | Thorough testing of JWT/API key flows |
| **Header Parsing Issues** | High | Test `X-Enable-Client-Encryption` explicitly |
| **Context Data Loss** | High | Verify all `Locals`→`Get/Set` conversions |
| **Error Response Changes** | Medium | Test error handling thoroughly |

### Medium-Risk Areas

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Performance Regression** | Medium | Benchmark before/after |
| **Frontend Integration** | Medium | Test web UI thoroughly |
| **CORS Changes** | Medium | Test from different origins |

### Low-Risk Areas

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Static File Serving** | Low | Simple path change |
| **Logging Format** | Low | Verify log output |

---

## Testing Strategy

### 1. Automated Testing

```bash
# Run all Go tests
JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linting
golangci-lint run ./...
gofmt -w .
```

### 2. Manual Testing Checklist

#### Authentication
- [ ] POST `/api/auth/register` - Create new user
- [ ] POST `/api/auth/login` - Login with credentials
- [ ] GET `/api/users/me` - Get current user (JWT)
- [ ] POST `/api/auth/login/oidc` - OIDC login (if configured)

#### API Keys
- [ ] POST `/api/api-keys` - Create API key
- [ ] GET `/api/api-keys` - List API keys
- [ ] DELETE `/api/api-keys/{id}` - Delete API key

#### Vaults (Web API)
- [ ] POST `/api/vaults` - Create vault
- [ ] GET `/api/vaults` - List vaults
- [ ] GET `/api/vaults/{id}` - Get vault
- [ ] PUT `/api/vaults/{id}` - Update vault
- [ ] DELETE `/api/vaults/{id}` - Delete vault

#### Vaults (CLI API)
- [ ] GET `/api/cli/vaults` - List vaults (API key auth)
- [ ] GET `/api/cli/vault/{uniqueId}` - Get vault (API key auth)
- [ ] GET `/api/cli/vault/name/{name}` - Get vault by name
- [ ] GET `/api/cli/vault/{uniqueId}` with `X-Enable-Client-Encryption: true`

#### System
- [ ] GET `/api/status` - Health check
- [ ] GET `/` - Frontend loads
- [ ] CORS from different origin

### 3. Performance Testing

```bash
# Benchmark before migration (Fiber)
ab -n 10000 -c 100 http://localhost:3000/api/status

# Benchmark after migration (Echo)
ab -n 10000 -c 100 http://localhost:3000/api/status

# Compare results (expect ~6% slower, acceptable)
```

### 4. Client-Side Encryption Testing

**Critical:** Verify the header parsing fix works!

```bash
# Test encrypted vault retrieval
curl -H "Authorization: Bearer vhub_xxx" \
     -H "X-Enable-Client-Encryption: true" \
     http://localhost:3000/api/cli/vault/abc123

# Expected: Encrypted value returned
# Verify: No compilation errors about []string vs string
```

---

## Rollback Plan

### If Migration Fails

1. **Immediately:**
   ```bash
   git checkout main
   git branch -D migrate-fiber-to-echo
   ```

2. **For deployed environments:**
   ```bash
   # Deploy previous version
   ./scripts/deploy.sh --version=v1.4.5
   ```

3. **No database rollback needed** (no schema changes)

### Rollback Decision Criteria

Rollback if:
- ❌ Critical authentication bugs discovered in production
- ❌ Performance degradation >20%
- ❌ Client-side encryption fails
- ❌ Unable to fix within 2 hours

Do NOT rollback if:
- ✅ Minor UI glitches (can fix forward)
- ✅ Performance degradation <10% (acceptable)
- ✅ Non-critical features affected

---

## Post-Migration Validation

### Day 1: Immediate Checks
- [ ] All endpoints responding
- [ ] No 5xx errors in logs
- [ ] Authentication working
- [ ] Frontend functional

### Week 1: Monitoring
- [ ] Error rate <1%
- [ ] Response times within 10% of baseline
- [ ] No memory leaks
- [ ] All background jobs running

### Month 1: Long-term Validation
- [ ] Performance stable
- [ ] No unexpected errors
- [ ] Client feedback positive
- [ ] Consider removing Fiber dependency entirely

---

## Success Criteria

Migration is considered successful when:

1. ✅ All automated tests pass
2. ✅ All manual test scenarios pass
3. ✅ No compilation errors or warnings
4. ✅ Header parameter bug eliminated
5. ✅ Performance within 10% of baseline
6. ✅ Frontend fully functional
7. ✅ CLI tool works with API
8. ✅ Client-side encryption working
9. ✅ Zero production incidents for 1 week
10. ✅ Code review approved

---

## Appendix: Quick Reference

### Fiber → Echo Cheat Sheet

```go
// Context
fiber.Ctx → echo.Context

// Response
c.JSON(data) → ctx.JSON(200, data)
c.Status(code).JSON(data) → ctx.JSON(code, data)
c.SendString(str) → ctx.String(200, str)

// Request
c.Params("id") → ctx.Param("id")
c.Query("q") → ctx.QueryParam("q")
c.Get("Header") → ctx.Request().Header.Get("Header")
c.BodyParser(&v) → ctx.Bind(&v)

// Context Storage
c.Locals("key") → ctx.Get("key")
c.Locals("key", val) → ctx.Set("key", val)

// Errors
fiber.NewError(code, msg) → echo.NewHTTPError(code, msg)
fiber.ErrUnauthorized → echo.ErrUnauthorized

// Middleware
fiber.Handler → echo.MiddlewareFunc
c.Next() → next(ctx)

// App
fiber.New() → echo.New()
app.Listen(":3000") → e.Start(":3000")
```

---

**End of Migration Plan**
