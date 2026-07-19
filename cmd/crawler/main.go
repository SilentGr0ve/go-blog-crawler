package main

import (
	"os"

	"github.com/SilentGr0ve/go-blog-crawler/internal/logger"
)

func main() {
	log, err := logger.NewLogger(logger.Config{
		Level:  "debug",
		Folder: "./out/logs",
	})

	if err != nil {
		os.Stderr.WriteString("failed to init logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer log.Close()

	log.Info("crawler starting")
	log.Debug("debug message", "key", "value")
}
