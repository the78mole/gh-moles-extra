package main

import (
	"os"

	"github.com/the78mole/gh-moles-extra/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
