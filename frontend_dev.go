//go:build !production

package main

import "io/fs"

func frontendFS() fs.FS {
	return nil
}
