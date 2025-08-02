package links

import (
	"golang.org/x/net/html"
)

type LinkExtractor interface {
	Extract(*html.Node) (string, bool)
}

// AHrefExtractor Extracts <a href="...">
type AHrefExtractor struct {
}

func (h *AHrefExtractor) Extract(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}
	return "", false
}

var _ LinkExtractor = (*AHrefExtractor)(nil)

// AreaHrefExtractor Extracts <area href="...">
type AreaHrefExtractor struct{}

func (e *AreaHrefExtractor) Extract(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "area" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}
	return "", false
}

var _ LinkExtractor = (*AreaHrefExtractor)(nil)

// LinkHrefExtractor Extracts <link href="...">
type LinkHrefExtractor struct{}

func (l *LinkHrefExtractor) Extract(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "link" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}
	return "", false
}

var _ LinkExtractor = (*LinkHrefExtractor)(nil)
