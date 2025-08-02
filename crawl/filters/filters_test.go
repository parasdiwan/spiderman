package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal link",
			input:    "example.com/page",
			expected: "example.com/page",
		},
		{
			name:     "link with spaces",
			input:    "  example.com/page  ",
			expected: "example.com/page",
		},
		{
			name:     "link with trailing slash",
			input:    "example.com/page/",
			expected: "example.com/page",
		},
		{
			name:     "link with fragment",
			input:    "example.com/page#",
			expected: "example.com/page",
		},
		{
			name:     "https link",
			input:    "https://example.com/page",
			expected: "example.com/page",
		},
		{
			name:     "http link",
			input:    "http://example.com/page",
			expected: "example.com/page",
		},
		{
			name:     "www link",
			input:    "www.example.com/page",
			expected: "example.com/page",
		},
		{
			name:     "complex link with all elements",
			input:    "  https://www.example.com/page/#  ",
			expected: "example.com/page",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces and symbols",
			input:    "  /#  ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeLink(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoEmpty_Match(t *testing.T) {
	filter := &NotEmpty{}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "valid link",
			link:     "https://example.com",
			expected: true,
		},
		{
			name:     "valid relative link",
			link:     "/about",
			expected: true,
		},
		{
			name:     "empty string",
			link:     "",
			expected: false,
		},
		{
			name:     "only spaces",
			link:     "   ",
			expected: false,
		},
		{
			name:     "only slashes and hash",
			link:     "/#",
			expected: false,
		},
		{
			name:     "only protocol",
			link:     "https://",
			expected: false,
		},
		{
			name:     "only www",
			link:     "www.",
			expected: false,
		},
		{
			name:     "single character",
			link:     "a",
			expected: true,
		},
		{
			name:     "valid link with spaces that becomes empty after sanitization",
			link:     "  https://www./#  ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoFragment_Match(t *testing.T) {
	filter := &NotFragment{}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "normal link",
			link:     "https://example.com/page",
			expected: true,
		},
		{
			name:     "relative link",
			link:     "/about",
			expected: true,
		},
		{
			name:     "fragment link",
			link:     "#section",
			expected: false,
		},
		{
			name:     "link with fragment at end",
			link:     "https://example.com/page#section",
			expected: true, // After sanitization, fragment is removed
		},
		{
			name:     "empty string",
			link:     "",
			expected: false,
		},
		{
			name:     "only hash",
			link:     "#",
			expected: false,
		},
		{
			name:     "hash with spaces",
			link:     " #section ",
			expected: false,
		},
		{
			name:     "multiple hashes",
			link:     "##section",
			expected: false,
		},
		{
			name:     "javascript void",
			link:     "#javascript:void(0)",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoMailLink_Match(t *testing.T) {
	filter := &NotMailLink{}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "normal link",
			link:     "https://example.com",
			expected: true,
		},
		{
			name:     "mailto link",
			link:     "mailto:user@example.com",
			expected: false,
		},
		{
			name:     "mailto with subject",
			link:     "mailto:user@example.com?subject=Hello",
			expected: false,
		},
		{
			name:     "mailto with spaces",
			link:     "  mailto:user@example.com  ",
			expected: false,
		},
		{
			name:     "relative link",
			link:     "/contact",
			expected: true,
		},
		{
			name:     "empty string",
			link:     "",
			expected: true,
		},
		{
			name:     "partial mailto in middle",
			link:     "https://example.com/mailto:fake",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoTelephone_Match(t *testing.T) {
	filter := &NotTelephone{}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "normal link",
			link:     "https://example.com",
			expected: true,
		},
		{
			name:     "tel link",
			link:     "tel:+1234567890",
			expected: false,
		},
		{
			name:     "tel with country code",
			link:     "tel:+1-555-123-4567",
			expected: false,
		},
		{
			name:     "tel with spaces",
			link:     "  tel:123456789  ",
			expected: false,
		},
		{
			name:     "relative link",
			link:     "/contact",
			expected: true,
		},
		{
			name:     "empty string",
			link:     "",
			expected: true,
		},
		{
			name:     "partial tel in middle",
			link:     "https://example.com/tel:fake",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoFile_Match(t *testing.T) {
	filter := &NotFile{}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		// Valid non-file links
		{
			name:     "normal web page",
			link:     "https://example.com/page",
			expected: true,
		},
		{
			name:     "relative link",
			link:     "/about",
			expected: true,
		},
		{
			name:     "page with query params",
			link:     "https://example.com/search?q=test",
			expected: true,
		},

		// Image files
		{
			name:     "jpg image",
			link:     "https://example.com/image.jpg",
			expected: false,
		},
		{
			name:     "jpeg image",
			link:     "https://example.com/photo.jpeg",
			expected: false,
		},
		{
			name:     "png image",
			link:     "https://example.com/logo.png",
			expected: false,
		},
		{
			name:     "gif image",
			link:     "https://example.com/animation.gif",
			expected: false,
		},
		{
			name:     "svg image",
			link:     "https://example.com/icon.svg",
			expected: false,
		},
		{
			name:     "webp image",
			link:     "https://example.com/modern.webp",
			expected: false,
		},
		{
			name:     "ico favicon",
			link:     "https://example.com/favicon.ico",
			expected: false,
		},

		// JavaScript and CSS files
		{
			name:     "javascript file",
			link:     "https://example.com/script.js",
			expected: false,
		},
		{
			name:     "css file",
			link:     "https://example.com/styles.css",
			expected: false,
		},
		{
			name:     "mjs module",
			link:     "https://example.com/module.mjs",
			expected: false,
		},

		// Video files
		{
			name:     "mp4 video",
			link:     "https://example.com/video.mp4",
			expected: false,
		},
		{
			name:     "webm video",
			link:     "https://example.com/video.webm",
			expected: false,
		},
		{
			name:     "avi video",
			link:     "https://example.com/old-video.avi",
			expected: false,
		},

		// Audio files
		{
			name:     "mp3 audio",
			link:     "https://example.com/song.mp3",
			expected: false,
		},
		{
			name:     "wav audio",
			link:     "https://example.com/sound.wav",
			expected: false,
		},
		{
			name:     "ogg audio",
			link:     "https://example.com/audio.ogg",
			expected: false,
		},

		// Document files
		{
			name:     "pdf document",
			link:     "https://example.com/document.pdf",
			expected: false,
		},
		{
			name:     "word document",
			link:     "https://example.com/report.doc",
			expected: false,
		},
		{
			name:     "word docx",
			link:     "https://example.com/report.docx",
			expected: false,
		},
		{
			name:     "excel file",
			link:     "https://example.com/data.xlsx",
			expected: false,
		},
		{
			name:     "powerpoint",
			link:     "https://example.com/presentation.pptx",
			expected: false,
		},

		// Archive files
		{
			name:     "zip archive",
			link:     "https://example.com/archive.zip",
			expected: false,
		},
		{
			name:     "rar archive",
			link:     "https://example.com/files.rar",
			expected: false,
		},
		{
			name:     "tar archive",
			link:     "https://example.com/backup.tar",
			expected: false,
		},
		{
			name:     "gzip file",
			link:     "https://example.com/compressed.gz",
			expected: false,
		},

		// Executable files
		{
			name:     "exe file",
			link:     "https://example.com/installer.exe",
			expected: false,
		},
		{
			name:     "dmg file",
			link:     "https://example.com/app.dmg",
			expected: false,
		},
		{
			name:     "apk file",
			link:     "https://example.com/app.apk",
			expected: false,
		},

		// Other files
		{
			name:     "xml file",
			link:     "https://example.com/sitemap.xml",
			expected: false,
		},
		{
			name:     "json file",
			link:     "https://example.com/data.json",
			expected: false,
		},
		{
			name:     "csv file",
			link:     "https://example.com/export.csv",
			expected: false,
		},
		{
			name:     "txt file",
			link:     "https://example.com/readme.txt",
			expected: false,
		},
		{
			name:     "rss feed",
			link:     "https://example.com/feed.rss",
			expected: false,
		},
		{
			name:     "font file",
			link:     "https://example.com/font.woff2",
			expected: false,
		},

		// Edge cases
		{
			name:     "file extension in path but not filename",
			link:     "https://example.com/pdf/page",
			expected: true,
		},
		{
			name:     "file extension with query params",
			link:     "https://example.com/image.jpg?version=1",
			expected: true,
		},
		{
			name:     "uppercase extension",
			link:     "https://example.com/IMAGE.JPG",
			expected: true, // Extensions are case sensitive
		},
		{
			name:     "no extension",
			link:     "https://example.com/page-without-extension",
			expected: true,
		},
		{
			name:     "empty string",
			link:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInternalLink_Match(t *testing.T) {
	// Test with example.com as base
	filter := NewInternalLink("https://example.com")

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		// Internal links (should match)
		{
			name:     "exact same domain",
			link:     "https://example.com",
			expected: true,
		},
		{
			name:     "same domain with path",
			link:     "https://example.com/about",
			expected: true,
		},
		{
			name:     "same domain with www",
			link:     "https://www.example.com/contact",
			expected: true,
		},
		{
			name:     "relative link with slash",
			link:     "/products",
			expected: true,
		},
		{
			name:     "relative link without slash",
			link:     "services",
			expected: true,
		},
		{
			name:     "fragment link",
			link:     "#section",
			expected: true,
		},
		{
			name:     "http version of same domain",
			link:     "http://example.com/page",
			expected: true,
		},
		{
			name:     "same domain different case",
			link:     "https://EXAMPLE.COM/page",
			expected: true,
		},

		// External links (should not match)
		{
			name:     "different domain",
			link:     "https://google.com",
			expected: false,
		},
		{
			name:     "subdomain of different domain",
			link:     "https://sub.google.com",
			expected: false,
		},
		{
			name:     "similar domain name",
			link:     "https://example-fake.com",
			expected: false,
		},
		{
			name:     "domain as subdirectory",
			link:     "https://malicious.com/example.com",
			expected: false,
		},

		// Edge cases
		{
			name:     "malformed URL",
			link:     "not-a-url",
			expected: true, // No :// so treated as internal
		},
		{
			name:     "empty string",
			link:     "",
			expected: true, // Relative path
		},
		{
			name:     "just protocol",
			link:     "https://",
			expected: false, // Malformed URL
		},
		{
			name:     "invalid URL with protocol",
			link:     "https://[invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.Match(tt.link)
			assert.Equal(t, tt.expected, result, "Link: %s", tt.link)
		})
	}
}

func TestInternalLink_Match_DifferentBases(t *testing.T) {
	tests := []struct {
		name     string
		baseUrl  string
		testLink string
		expected bool
	}{
		{
			name:     "subdomain base - same subdomain",
			baseUrl:  "https://blog.example.com",
			testLink: "https://blog.example.com/post",
			expected: true,
		},
		{
			name:     "subdomain base - different subdomain",
			baseUrl:  "https://blog.example.com",
			testLink: "https://shop.example.com/item",
			expected: false,
		},
		{
			name:     "subdomain base - main domain",
			baseUrl:  "https://blog.example.com",
			testLink: "https://example.com",
			expected: false,
		},
		{
			name:     "port in base URL",
			baseUrl:  "https://localhost:8080",
			testLink: "https://localhost:8080/api",
			expected: true,
		},
		{
			name:     "port in base URL - different port",
			baseUrl:  "https://localhost:8080",
			testLink: "https://localhost:3000/api",
			expected: false,
		},
		{
			name:     "no www base with www link",
			baseUrl:  "https://github.com",
			testLink: "https://www.github.com/user",
			expected: true,
		},
		{
			name:     "www base with no www link",
			baseUrl:  "https://www.facebook.com",
			testLink: "https://facebook.com/page",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewInternalLink(tt.baseUrl)
			result := filter.Match(tt.testLink)
			assert.Equal(t, tt.expected, result, "Base: %s, Link: %s", tt.baseUrl, tt.testLink)
		})
	}
}

func TestInternalLink_Match_RealWorldExamples(t *testing.T) {
	// Test with real-world scenarios
	githubFilter := NewInternalLink("https://github.com")

	githubTests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "github main page",
			link:     "https://github.com",
			expected: true,
		},
		{
			name:     "github user profile",
			link:     "https://github.com/octocat",
			expected: true,
		},
		{
			name:     "github repository",
			link:     "https://github.com/octocat/Hello-World",
			expected: true,
		},
		{
			name:     "github with www",
			link:     "https://www.github.com/features",
			expected: true,
		},
		{
			name:     "github relative link",
			link:     "/pricing",
			expected: true,
		},
		{
			name:     "external link to stackoverflow",
			link:     "https://stackoverflow.com",
			expected: false,
		},
		{
			name:     "external link to twitter",
			link:     "https://twitter.com/github",
			expected: false,
		},
		{
			name:     "github pages subdomain",
			link:     "https://pages.github.com",
			expected: false, // Different subdomain
		},
	}

	for _, tt := range githubTests {
		t.Run(tt.name, func(t *testing.T) {
			result := githubFilter.Match(tt.link)
			assert.Equal(t, tt.expected, result, "Link: %s", tt.link)
		})
	}
}

func TestFilters_Integration(t *testing.T) {
	baseUrl := "https://example.com"
	filters := []Filter{
		&NotEmpty{},
		NewInternalLink(baseUrl),
		&NotFragment{},
		&NotMailLink{},
		&NotTelephone{},
		&NotFile{},
	}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "valid internal page",
			link:     "https://example.com/about",
			expected: true,
		},
		{
			name:     "valid relative link",
			link:     "/contact",
			expected: true,
		},
		{
			name:     "external link",
			link:     "https://google.com",
			expected: false,
		},
		{
			name:     "internal image file",
			link:     "https://example.com/logo.png",
			expected: false,
		},
		{
			name:     "internal mailto",
			link:     "mailto:contact@example.com",
			expected: false,
		},
		{
			name:     "fragment link",
			link:     "#section",
			expected: false,
		},
		{
			name:     "empty link",
			link:     "",
			expected: false,
		},
		{
			name:     "telephone link",
			link:     "tel:+1234567890",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allMatch := true
			for _, filter := range filters {
				if !filter.Match(tt.link) {
					allMatch = false
					break
				}
			}
			assert.Equal(t, tt.expected, allMatch, "Link: %s", tt.link)
		})
	}
}
