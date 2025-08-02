package crawl

import (
	"fmt"
	"spiderman/publish"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_Integration_SmallSite(t *testing.T) {
	publisher := publish.NewTestPublisher()
	crawler := NewCrawler("https://duckduckgo.com", publisher)

	err := crawler.Crawl()
	assert.NoError(t, err)
	sequentialResults := len(publisher.Published)

	// Reset publisher for parallel test
	publisher2 := publish.NewTestPublisher()
	crawler2 := NewCrawler("https://duckduckgo.com", publisher2)

	err = crawler2.CrawlParallel(2)
	assert.NoError(t, err)
	parallelResults := len(publisher2.Published)

	assert.Greater(t, sequentialResults, 0)
	assert.Greater(t, parallelResults, 0)
	assert.Equal(t, sequentialResults, parallelResults)

	fmt.Printf("Sequential found %d links, Parallel found %d links\n",
		sequentialResults, parallelResults)
}

func TestCrawler_Crawl_EmptyWebsite(t *testing.T) {
	publisher := publish.NewTestPublisher()
	crawler := NewCrawler("https://jsonplaceholder.typicode.com/users", publisher) // Simple HTML page

	start := time.Now()
	err := crawler.Crawl()
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotEmpty(t, publisher.Published)

	assert.GreaterOrEqual(t, len(publisher.Published), 1)
	assert.Equal(t, "https://jsonplaceholder.typicode.com/users", publisher.Published[0])

	// Should complete quickly
	assert.Less(t, duration, 5*time.Second)
}

func TestCrawler_Crawl_InvalidDomain(t *testing.T) {
	publisher := publish.NewTestPublisher()
	crawler := NewCrawler("https://thisdoesnotexist12345.com", publisher)

	err := crawler.Crawl()

	assert.Error(t, err)
	assert.Empty(t, publisher.Published)

	fmt.Printf("Invalid domain test: %d links found\n", len(publisher.Published))
}

func TestCrawler_CrawlParallel_InvalidDomain(t *testing.T) {
	publisher := publish.NewTestPublisher()
	crawler := NewCrawler("https://thisdoesnotexist12345.com", publisher)

	_ = crawler.CrawlParallel(3)
	assert.NotPanics(t, func() {
		_ = crawler.CrawlParallel(3)
	})

	fmt.Printf("Invalid domain parallel test: %d links found\n", len(publisher.Published))
}

func TestCrawler_IsCrawlable(t *testing.T) {
	publisher := publish.NewTestPublisher()
	crawler := NewCrawler("https://example.com", publisher)

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "valid internal link",
			link:     "https://example.com/page1",
			expected: true,
		},
		{
			name:     "empty link",
			link:     "",
			expected: false,
		},
		{
			name:     "external link",
			link:     "https://google.com",
			expected: false,
		},
		{
			name:     "link with fragment",
			link:     "#section",
			expected: false,
		},
		{
			name:     "mailto link",
			link:     "mailto:test@example.com",
			expected: false,
		},
		{
			name:     "telephone link",
			link:     "tel:+1234567890",
			expected: false,
		},
		{
			name:     "file link",
			link:     "https://example.com/file.pdf",
			expected: false,
		},
		{
			name:     "relative path",
			link:     "/about",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := crawler.isCrawlable(tt.link)
			assert.Equal(t, tt.expected, result, "Link: %s", tt.link)
		})
	}
}

func TestCrawler_BuildAbsolutePath(t *testing.T) {
	tests := []struct {
		name     string
		baseUrl  string
		link     string
		expected string
	}{
		{
			name:     "https base with relative link",
			baseUrl:  "https://example.com",
			link:     "/about",
			expected: "https://example.com/about",
		},
		{
			name:     "http base with relative link",
			baseUrl:  "http://example.com",
			link:     "/contact",
			expected: "http://example.com/contact",
		},
		{
			name:     "base without protocol with relative link",
			baseUrl:  "example.com",
			link:     "/services",
			expected: "http://example.com/services",
		},
		{
			name:     "absolute link already contains base domain",
			baseUrl:  "https://example.com",
			link:     "example.com/products",
			expected: "https://example.com/products",
		},
		{
			name:     "relative link without leading slash",
			baseUrl:  "https://example.com",
			link:     "about",
			expected: "https://example.com/about",
		},
		{
			name:     "base with trailing slash",
			baseUrl:  "https://example.com/",
			link:     "about",
			expected: "https://example.com/about",
		},
		{
			name:     "complex relative path",
			baseUrl:  "https://blog.example.com",
			link:     "posts/2023/article",
			expected: "https://blog.example.com/posts/2023/article",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := publish.NewTestPublisher()
			crawler := NewCrawler(tt.baseUrl, publisher)
			result := crawler.buildAbsolutePath(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}
