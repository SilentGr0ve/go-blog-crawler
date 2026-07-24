package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SilentGr0ve/go-blog-crawler/internal/crawler"
	"github.com/SilentGr0ve/go-blog-crawler/internal/fetcher"
	"github.com/spf13/cobra"
)

var crawlOptions CrawlOptions

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Start crawling from seed URLs",
	Long:  "Start crawling from seed URLs and follow links up to the given depth",
	RunE: func(cmd *cobra.Command, args []string) error {

		seeds, err := readLines(crawlOptions.Seeds)
		if err != nil {
			return fmt.Errorf("read seeds: %w", err)
		}

		f := fetcher.NewClient(fetcher.Options{Timeout: 10 * time.Second})
		workerPool := crawler.NewWorkerPool(f, appLogger, crawlOptions.Workers)

		results, err := workerPool.Run(
			cmd.Context(),
			seeds, crawler.Options{
				MaxDepth: crawlOptions.MaxDepth,
				MaxPages: crawlOptions.MaxPages,
			})
		if err != nil {
			return fmt.Errorf("crawler: %w", err)
		}

		for _, res := range results {
			fmt.Printf("%s %s\n", res.URL, res.Title)
		}

		appLogger.Info(
			"all pages crawled",
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().StringVar(&crawlOptions.Seeds, "seeds", "seeds.txt", "path to seeds file")
	crawlCmd.Flags().IntVar(&crawlOptions.MaxPages, "max-pages", 50, "max pages to visit (0 = unlimited)")
	crawlCmd.Flags().IntVar(&crawlOptions.MaxDepth, "depth", 1, "max crawl depth")
	crawlCmd.Flags().IntVar(&crawlOptions.Workers, "concurrency", 5, "count of crawl workers")
}

func readLines(seedsPath string) ([]string, error) {
	seedFile, err := os.Open(seedsPath)
	if err != nil {
		return nil, err
	}
	defer seedFile.Close()

	scanner := bufio.NewScanner(seedFile)
	seeds := make([]string, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		seeds = append(seeds, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return seeds, nil
}
