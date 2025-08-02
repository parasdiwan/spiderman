package links

import (
	"errors"
	"io"
	"log"
	"spiderman/crawl/filters"
	"spiderman/crawl/http"

	"golang.org/x/net/html"
)

type Parser struct {
	extractors []LinkExtractor
	fetcher    *http.Fetcher
	filters    []filters.Filter
}

func NewParser() *Parser {
	return &Parser{
		extractors: []LinkExtractor{
			&AHrefExtractor{},
			&LinkHrefExtractor{},
			&AreaHrefExtractor{},
		},
		fetcher: http.NewFetcher(),
		filters: []filters.Filter{
			&filters.NotFile{},
			&filters.NotEmpty{},
			&filters.NotMailLink{},
			&filters.NotTelephone{},
		},
	}
}

func (p *Parser) FetchLinks(baseUrl string) ([]string, error) {
	result := p.fetcher.Fetch(baseUrl)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Location != "" {
		// It’s a redirect: return as a single link
		return []string{result.Location}, nil
	}

	if result.Body == nil {
		return nil, errors.New("there's no body here")
	}
	defer result.Body.Close()
	return p.fetchURLsFromHtml(result.Body)
}

// input links is assumed to be utf-8 encoded
func (p *Parser) fetchURLsFromHtml(reader io.ReadCloser) ([]string, error) {
	baseNode, err := html.Parse(reader)
	if err != nil {
		log.Printf("[Error] failed to parse links: %s\n", err)
		return nil, err
	}

	links := make([]string, 0)
	p.extractLinks(baseNode, &links)
	return links, nil
}

func (p *Parser) extractLinks(node *html.Node, links *[]string) {
	for _, ex := range p.extractors {
		link, exists := ex.Extract(node)
		if !exists || !p.isValidLink(link) {
			continue
		}
		*links = append(*links, link)
	}
	node = node.FirstChild
	for ; node != nil; node = node.NextSibling {
		p.extractLinks(node, links)
	}
}

func (p *Parser) isValidLink(link string) bool {
	for _, filter := range p.filters {
		if !filter.Match(link) {
			return false
		}
	}
	return true
}
