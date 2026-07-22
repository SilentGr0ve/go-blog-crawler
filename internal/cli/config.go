package cli

type CLIOptions struct {
	LogDir   string
	LogLevel string
}

type CrawlOptions struct {
	Seeds    string
	MaxPages int
	MaxDepth int
}
