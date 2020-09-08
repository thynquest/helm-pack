package main

import (
	"os"

	"github.com/thynquest/helm-pack/cmd/helmpack"
)

func main() {
	cmd := helmpack.NewPackCmd(os.Args[1:], os.Stdout)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
