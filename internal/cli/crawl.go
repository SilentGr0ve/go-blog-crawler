package cli

import "github.com/spf13/cobra"

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Start crawling from seed URLs",
	Long:  "Start crawling from seed URLs and follow links up to the given depth",
	RunE: func(cmd *cobra.Command, args []string) error {
		appLogger.Info("crawl command started...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)
}
