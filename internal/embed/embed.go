package embed

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

// GetDistFS returns the embedded filesystem containing the web dist files
func GetDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}

// GetHTTPFS returns an http.FileSystem for serving the embedded files
func GetHTTPFS() (http.FileSystem, error) {
	sub, err := GetDistFS()
	if err != nil {
		return nil, err
	}
	return http.FS(sub), nil
}