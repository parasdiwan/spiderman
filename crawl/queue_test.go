package crawl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueBehavior(t *testing.T) {
	q := NewFifoQueue()
	// Single add/grab
	q.Add("x")
	assert.Equal(t, "x", q.Grab())
	// Empty queue returns empty string
	assert.Equal(t, "", q.Grab())
	// AddAll and repeated Grab
	q.Add("a")
	q.Add("b")
	q.Add("c")
	assert.Equal(t, "a", q.Grab())
	assert.Equal(t, "b", q.Grab())
	assert.Equal(t, "c", q.Grab())
	assert.Equal(t, "", q.Grab())
}
