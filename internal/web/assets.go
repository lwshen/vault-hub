package web

import (
	"embed"
	"io/fs"
)

// Embed the web dist directory
// This will be populated when the web app is built
//go:embed dist
var distFS embed.FS

// GetDistFS returns the embedded web assets filesystem
func GetDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}

// HasAssets returns true if web assets are embedded
func HasAssets() bool {
	entries, err := distFS.ReadDir(".")
	if err != nil {
		return false
	}
	
	// Check if dist directory exists and has content
	for _, entry := range entries {
		if entry.Name() == "dist" && entry.IsDir() {
			distEntries, err := distFS.ReadDir("dist")
			return err == nil && len(distEntries) > 0
		}
	}
	return false
}