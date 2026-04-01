//go:build production

package main

import (
	"embed"
	"io/fs"
)

//go:embed all:web/build
var frontendBuildDir embed.FS

func frontendFS() fs.FS {
	sub, _ := fs.Sub(frontendBuildDir, "web/build")
	return sub
}
