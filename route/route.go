package route

import (
    "io/fs"
    "path/filepath"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/lwshen/vault-hub/handler"
    "github.com/lwshen/vault-hub/internal/staticfs"
    openapi "github.com/lwshen/vault-hub/packages/api"
)

func SetupRoutes(app *fiber.App) {
	app.Use(jwtMiddleware)

	server := openapi.NewServer()
	openapi.RegisterHandlers(app, server)

	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Get("/login/oidc", handler.LoginOidc)
	auth.Get("/callback/oidc", handler.LoginOidcCallback)

    // Web (embedded)
    sub, err := fs.Sub(staticfs.WebDist, "dist")
    if err == nil {
        app.Get("/*", func(c *fiber.Ctx) error {
            reqPath := c.Params("*")
            if reqPath == "" || strings.HasSuffix(reqPath, "/") {
                reqPath = reqPath + "index.html"
            }

            // Try to serve the requested asset
            if f, openErr := sub.Open(reqPath); openErr == nil {
                defer f.Close()
                c.Type(contentTypeForPath(reqPath))
                return c.SendStream(f)
            }

            // Fallback to SPA index.html
            f, openErr := sub.Open("index.html")
            if openErr != nil {
                return c.SendStatus(fiber.StatusInternalServerError)
            }
            defer f.Close()
            c.Type("html")
            return c.SendStream(f)
        })
    }
}

func contentTypeForPath(p string) string {
    ext := strings.ToLower(filepath.Ext(p))
    switch ext {
    case ".html":
        return "html"
    case ".css":
        return "css"
    case ".js":
        return "javascript"
    case ".json":
        return "json"
    case ".svg":
        return "svg"
    case ".png":
        return "png"
    case ".jpg", ".jpeg":
        return "jpeg"
    case ".ico":
        return "x-icon"
    case ".map":
        return "json"
    case ".woff":
        return "font"
    case ".woff2":
        return "font"
    case ".ttf":
        return "font"
    default:
        return "octet-stream"
    }
}
