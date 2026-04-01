package main

import (
	"fmt"
	"os"

	"github.com/z2-cli/cmd"
)

func main() {
	cmd.FrontendFS = frontendFS()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
