package cli

import (
	"fmt"

	"github.com/SilentGr0ve/go-blog-crawler/internal/logger"
	"github.com/spf13/cobra"
)

var (
	opts      CLIOptions
	appLogger *logger.Logger
)

var rootCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Concurrent Web crawler for Go blogs",
	Long:  "go-blog-crawler is a concurrent web crawler that discovers and fetches Go-related blog pages, respects robots.txt, and exports results to JSON",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loggerConfig := logger.Config{
			Level:  opts.LogLevel,
			Folder: opts.LogDir,
		}

		l, err := logger.NewLogger(loggerConfig)
		if err != nil {
			return fmt.Errorf("failed to init logger: %w", err)
		}

		appLogger = l
		if appLogger == nil {
			return cmd.Help()
		}

		appLogger.Info("crawler starting...")
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&opts.LogLevel,
		"log-level",
		"info",
		"log level (debug|info|warn|error)",
	)

	rootCmd.PersistentFlags().StringVar(
		&opts.LogDir,
		"log-dir",
		"",
		"directory for log files (if empty = stderr only)",
	)
}
