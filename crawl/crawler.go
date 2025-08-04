package crawl

import (
	"fmt"
	"log"
	"spiderman/crawl/filters"
	"strings"
	"sync"

	"spiderman/crawl/links"
	"spiderman/publish"
)

type Crawler struct {
	parser    *links.Parser
	publisher publish.Publisher
	filters   []filters.Filter
	baseUrl   string
}

func NewCrawler(baseUrl string, publisher publish.Publisher) *Crawler {
	return &Crawler{
		parser:    links.NewParser(),
		publisher: publisher,
		filters: []filters.Filter{
			&filters.NotEmpty{},
			filters.NewInternalLink(baseUrl),
			&filters.NotFragment{},
			&filters.NotMailLink{},
			&filters.NotTelephone{},
			&filters.NotFile{},
		},
		baseUrl: baseUrl,
	}
}

func (m *Crawler) isCrawlable(link string) bool {
	for _, filter := range m.filters {
		if !filter.Match(link) {
			return false
		}
	}
	return true
}

func (m *Crawler) buildAbsolutePath(link string) string {
	internalLink := filters.SanitizeLink(link)
	baseDomain := filters.SanitizeLink(m.baseUrl)
	// taking default as http
	protocol := "http://"
	if strings.Contains(m.baseUrl, "://") && strings.HasPrefix(m.baseUrl, "https://") {
		protocol = "https://"
	}
	if strings.HasPrefix(internalLink, baseDomain) {
		return protocol + internalLink
	}
	return protocol + baseDomain + "/" + internalLink
}

func (m *Crawler) Crawl() error {
	_, err := m.parser.FetchLinks(m.baseUrl)
	if err != nil {
		return fmt.Errorf("failed to access initial URL %s: %w", m.baseUrl, err)
	}

	queue := NewFifoQueue()

	queue.Add(m.baseUrl)
	nextUrl := queue.Grab()
	for ; nextUrl != ""; nextUrl = queue.Grab() {
		m.crawlAndPublishLinks(nextUrl, queue)
	}
	_ = m.publisher.PublishStats()
	return nil
}

func (m *Crawler) crawlAndPublishLinks(nextUrl string, queue *FifoQueue) {
	linksForPage, err := m.parser.FetchLinks(nextUrl)
	if err != nil {
		_ = m.publisher.RecordError(nextUrl, resolveErrType(err), err)
		log.Printf("[Error] failed to crawl page: %s\n", err)
		return
	}
	err = m.publisher.Publish(nextUrl, linksForPage)
	if err != nil {
		log.Printf("[Error] failed to publish page: %s\n", err)
	}
	for _, link := range linksForPage {
		if m.isCrawlable(link) {
			queue.Add(m.buildAbsolutePath(link))
		}
	}
}

// CrawlParallel starts the crawl with the specified number of workers.
func (c *Crawler) CrawlParallel(maxWorkers int) error {
	bufferSize := maxWorkers * 500
	queue := NewTaskQueue(c.baseUrl, bufferSize)
	var wg sync.WaitGroup
	wg.Add(1)

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		go func(id int) {
			for url := range queue.QueuedTasks() {
				c.processURL(id, url, queue, &wg)
			}
		}(i)
	}

	// Wait for all URLs to be processed
	wg.Wait()
	queue.Close()

	return c.publisher.PublishStats()
}

// processURL processes a single URL: publishing, fetching links, and enqueuing.
func (c *Crawler) processURL(workerID int, url string, queue *TaskQueue, wg *sync.WaitGroup) {
	defer wg.Done()

	// If already visited, skip
	if !queue.MarkVisited(url) {
		return
	}

	linksForPage, err := c.parser.FetchLinks(url)
	if err != nil {
		log.Printf("[Worker %d] Error fetching %s: %v", workerID, url, err)
		_ = c.publisher.RecordError(url, resolveErrType(err), err)
		return
	}
	_ = c.publisher.Publish(url, linksForPage)

	for _, link := range linksForPage {
		if !c.isCrawlable(link) {
			continue
		}
		absoluteLink := c.buildAbsolutePath(link)
		if queue.Add(absoluteLink) {
			wg.Add(1)
		}
	}
}

func resolveErrType(err error) publish.ErrType {
	switch err.Error() {
	case "failed with status 500":
		return publish.ErrTypeInternal
	case "failed with status 404":
		return publish.ErrTypeNotFound
	case "failed with status 403":
		return publish.ErrTypeNoAccess
	default:
		fmt.Println("-------------------ERROR FOUND------------------------------")
		fmt.Println(err.Error())
	}
	return publish.ErrTypeUnknown
}
