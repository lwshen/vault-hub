package staticfs

import "embed"

// WebDist embeds the built web assets from apps/web/dist
//go:embed all:dist
var WebDist embed.FS

