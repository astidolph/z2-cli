package cmd

import "io/fs"

// FrontendFS holds the embedded frontend filesystem, set by main before Execute.
var FrontendFS fs.FS
