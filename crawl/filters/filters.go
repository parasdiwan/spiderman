package filters

import (
	"net/url"
	"strings"
)

// Filter interface defines the criteria for the links which should be traversed.
// setting up as interfaces because I see this as domain specs,
// which could be modified in the future as we find out more cases
// ## FOR FUTURE:
// - we could add filters for file links
// - add filters according to robots.txt
type Filter interface {
	// Match function takes the links and returns true only if the link matches the criteria specified.
	Match(string) bool
}

func SanitizeLink(link string) string {
	link = strings.Trim(link, " ")
	link = strings.TrimSuffix(link, "#")
	link = strings.TrimPrefix(link, "https://")
	link = strings.TrimPrefix(link, "http://")

	link = strings.Trim(link, "/")
	link = strings.TrimPrefix(link, "www.")

	return link
}

// NotFragment filter checks that the link is not href starting with #
// these links are just sections of the page
type NotFragment struct{}

func (f *NotFragment) Match(link string) bool {
	link = SanitizeLink(link)
	return link != "" && !strings.HasPrefix(link, "#")
}

// internalLink filters out all possible external links
// e.g. `twitter.com/something`  OR `subdomain.domain.com`
type internalLink struct {
	// stored as sanitized baseURL
	baseUrl string
	// storing baseDomain it so don't have to recalculate it each time
	// on each filter checks
	baseDomain string
}

func NewInternalLink(baseUrl string) Filter {
	baseUrl = SanitizeLink(baseUrl)
	// split is guaranteed to return 1 element
	baseDomain := strings.Split(baseUrl, "/")[0]
	return &internalLink{baseUrl, baseDomain}
}

func (l *internalLink) Match(link string) bool {
	link = strings.Trim(link, " ")
	if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "/") {
		return true
	}
	// all external links are supposed to have ://
	// if they don't it must be internal
	if !strings.Contains(link, "://") {
		return true
	}
	u, err := url.Parse(link)
	if err != nil || u.Host == "" {
		return false
	}

	// treating www.facebook.com and facebook.com as the same website.
	linkDomain := strings.ToLower(strings.TrimPrefix(u.Host, "www."))
	return linkDomain == l.baseDomain
}

type NotEmpty struct{}

func (e *NotEmpty) Match(link string) bool {
	link = SanitizeLink(link)
	link = strings.Trim(link, "#")
	return link != ""
}

type NotMailLink struct{}

func (e *NotMailLink) Match(link string) bool {
	link = strings.Trim(link, " ")
	return !strings.HasPrefix(link, "mailto:")
}

type NotTelephone struct{}

func (e *NotTelephone) Match(link string) bool {
	link = strings.Trim(link, " ")
	return !strings.HasPrefix(link, "tel:")
}

type NotFile struct{}

var fileLinkSuffixes = []string{
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".tiff", ".ico",
	".js", ".mjs", ".cjs", ".css",
	".mp4", ".webm", ".ogv", ".avi", ".mov", ".flv", ".mkv", ".wmv",
	".mp3", ".wav", ".ogg", ".m4a", ".flac",
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp",
	".zip", ".rar", ".7z", ".tar", ".gz", ".bz2",
	".exe", ".dmg", ".apk", ".bin", ".iso",
	".csv", ".txt", ".xml", ".json", ".rss", ".rss.xml", ".woff2",
}

func (f *NotFile) Match(link string) bool {
	for _, suffix := range fileLinkSuffixes {
		if strings.HasSuffix(link, suffix) {
			return false
		}
	}
	return true
}

var _ Filter = (*NotEmpty)(nil)
var _ Filter = (*NotFragment)(nil)
var _ Filter = (*internalLink)(nil)
var _ Filter = (*NotMailLink)(nil)
var _ Filter = (*NotTelephone)(nil)
var _ Filter = (*NotFile)(nil)
