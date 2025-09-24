package apps

import (
    "embed"
    "io/fs"
)

// WebDistFS contains the embedded production build of the web app.
//go:embed web/dist
var embeddedWebFS embed.FS

// WebDistFS exposes the embedded web dist directory as a standard fs.FS
var WebDistFS fs.FS

func init() {
    // Scope the embedded FS to the dist directory root
    sub, err := fs.Sub(embeddedWebFS, "web/dist")
    if err != nil {
        panic(err)
    }
    WebDistFS = sub
}

