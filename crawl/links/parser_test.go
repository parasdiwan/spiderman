package links

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_FetchLinks(t *testing.T) {
	parser := NewParser()

	result, err := parser.FetchLinks("http://monzo.com")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(result), 121)
}

func TestParser_fetchURLsFromHtml(t *testing.T) {
	htmlStr := `<links>
		<head><title>Test</title></head>
		<body>
		  <a href="http://example.com/1">One</a>
		  <a href="/2">Two</a>
		  <a>Name Only</a>
		</body>
	</links>`
	r := io.NopCloser(strings.NewReader(htmlStr))
	parser := NewParser()
	links, err := parser.fetchURLsFromHtml(r)
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://example.com/1", "/2"}, links)
}
