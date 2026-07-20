package main

import (
	"os"

	"github.com/SilentGr0ve/go-blog-crawler/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
