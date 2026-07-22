package crawler

import (
	"context"
	"net/url"

	"github.com/SilentGr0ve/go-blog-crawler/internal/fetcher"
	"github.com/SilentGr0ve/go-blog-crawler/internal/logger"
	"github.com/SilentGr0ve/go-blog-crawler/internal/parser"
)

type Options struct {
	MaxDepth int
	MaxPages int
}

type Sequential struct {
	fetcher *fetcher.Client
	logger  *logger.Logger
}

func NewSequential(fetcher *fetcher.Client, logger *logger.Logger) *Sequential {
	return &Sequential{
		fetcher: fetcher,
		logger:  logger,
	}
}

func (s *Sequential) Run(ctx context.Context, seeds []string, options Options) ([]Result, error) {
	s.logger.Info(
		"crawler starting",
		"seeds", len(seeds),
		"max_depth", options.MaxDepth,
		"max_pages", options.MaxPages,
	)

	queue := make([]Job, 0, len(seeds))
	for _, seed := range seeds {
		queue = append(queue, Job{URL: seed, Depth: 0})
	}

	visited := make(map[string]bool)
	var results []Result

	for len(queue) > 0 {
		if err := ctx.Err(); err != nil {
			s.logger.Warn("cancelled by context")
			return results, err
		}

		job := queue[0]
		queue = queue[1:]

		if visited[job.URL] {
			continue
		}
		visited[job.URL] = true

		if options.MaxPages > 0 && len(visited) > options.MaxPages {
			s.logger.Info(
				"reached max_pages count",
				"max_pages", options.MaxPages,
			)
			break
		}

		body, statusCode, err := s.fetcher.Fetch(ctx, job.URL)
		if err != nil {
			s.logger.Warn(
				"fetch failed",
				"url", job.URL,
				"depth", job.Depth,
				"error", err,
			)
			continue
		}
		s.logger.Debug(
			"fetch completed",
			"url", job.URL,
			"status_code", statusCode,
			"body_size", len(body),
		)

		page, err := parser.Parse(job.URL, body)
		if err != nil {
			s.logger.Warn(
				"parse failed",
				"url", job.URL,
				"depth", job.Depth,
				"error", err)
			continue
		}
		s.logger.Debug(
			"parse completed",
			"url", page.URL,
			"links_count", len(page.Links),
		)

		results = append(results, Result{
			URL:   page.URL,
			Title: page.Title,
		})

		if job.Depth < options.MaxDepth {
			for _, link := range page.Links {
				if !visited[link] && sameHostOnly(job.URL, link) {
					queue = append(queue, Job{
						URL:   link,
						Depth: job.Depth + 1,
					})
				}
			}
		}
	}
	s.logger.Info(
		"crawl finished",
		"visited", len(visited),
		"results", len(results),
	)
	return results, nil
}

func sameHostOnly(a, b string) bool {
	ua, err := url.Parse(a)
	if err != nil {
		return false
	}
	ub, err := url.Parse(b)
	if err != nil {
		return false
	}
	return ua.Host == ub.Host
}
