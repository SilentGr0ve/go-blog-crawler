package crawler

import (
	"context"
	"sync"

	"github.com/SilentGr0ve/go-blog-crawler/internal/fetcher"
	"github.com/SilentGr0ve/go-blog-crawler/internal/logger"
	"github.com/SilentGr0ve/go-blog-crawler/internal/parser"
)

type WorkerPool struct {
	fetcher *fetcher.Client
	logger  *logger.Logger

	workers int
	jobsWG  sync.WaitGroup

	mu      sync.Mutex
	visited map[string]bool

	limitLogged bool
}

func NewWorkerPool(fetcher *fetcher.Client, log *logger.Logger, workers int) *WorkerPool {
	return &WorkerPool{
		fetcher: fetcher,
		logger:  log,
		workers: workers,
		visited: make(map[string]bool),
	}
}

func (w *WorkerPool) worker(ctx context.Context, worker int, jobs chan Job, results chan<- Result, wg *sync.WaitGroup, options Options) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			w.logger.Debug(
				"start process link",
				"url", job.URL,
				"worker", worker,
			)
			newJobs := w.processJob(ctx, job, results, options)
			for _, newJob := range newJobs {
				w.jobsWG.Add(1)
				select {
				case jobs <- newJob:
				case <-ctx.Done():
					w.jobsWG.Done()
					w.jobsWG.Done()
					w.logger.Warn("cancelled by context")
					return
				}
			}
			w.jobsWG.Done()
		case <-ctx.Done():
			w.logger.Warn("cancelled by context")
			return
		}

	}
}

func (w *WorkerPool) processJob(ctx context.Context, job Job, results chan<- Result, options Options) []Job {
	if ctx.Err() != nil {
		return nil
	}

	body, statusCode, err := w.fetcher.Fetch(ctx, job.URL)
	if err != nil {
		w.logger.Warn(
			"fetch failed",
			"url", job.URL,
			"depth", job.Depth,
			"error", err,
		)
		return nil
	}

	w.logger.Debug(
		"fetch completed",
		"url", job.URL,
		"status_code", statusCode,
		"body_size", len(body),
	)

	page, err := parser.Parse(job.URL, body)
	if err != nil {
		w.logger.Warn(
			"parse failed",
			"url", job.URL,
			"depth", job.Depth,
			"error", err,
		)
		return nil
	}

	w.logger.Debug(
		"parse completed",
		"url", page.URL,
		"links_count", len(page.Links),
	)

	select {
	case results <- Result{
		URL:   page.URL,
		Title: page.Title,
	}:
	case <-ctx.Done():
		return nil
	}

	var newJobs []Job
	if job.Depth < options.MaxDepth {
		for _, link := range page.Links {
			if sameHostOnly(job.URL, link) && w.canVisit(link, options.MaxPages) {
				newJobs = append(
					newJobs,
					Job{
						URL:   link,
						Depth: job.Depth + 1,
					})
			}
		}
	}
	return newJobs
}

func (w *WorkerPool) Run(ctx context.Context, seeds []string, options Options) ([]Result, error) {
	jobs := make(chan Job, 1000)
	results := make(chan Result, 1000)

	var workerWG sync.WaitGroup
	for i := 0; i < w.workers; i++ {
		workerWG.Add(1)
		go w.worker(ctx, i, jobs, results, &workerWG, options)
	}

	for _, seed := range seeds {
		if w.canVisit(seed, options.MaxPages) {
			w.jobsWG.Add(1)
			jobs <- Job{
				URL:   seed,
				Depth: 0,
			}
		}
	}

	go func() {
		w.jobsWG.Wait()
		close(jobs)
	}()

	workerWG.Wait()
	close(results)

	var all []Result
	for result := range results {
		all = append(all, result)
	}

	w.logger.Info(
		"crawl finished",
		"visited", len(w.visited),
		"results", len(all),
	)

	return all, nil
}

func (w *WorkerPool) canVisit(seed string, maxPages int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.visited[seed] {
		return false
	}

	if maxPages > 0 && len(w.visited) >= maxPages {
		if !w.limitLogged {
			w.logger.Info(
				"reached max_pages count",
				"max_pages", maxPages,
			)
			w.limitLogged = true
		}
		return false
	}

	w.visited[seed] = true
	return true
}
