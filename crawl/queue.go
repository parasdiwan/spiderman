package crawl

import "sync"

// FifoQueue guarantees ordering for messages with a First-in first-out approach
// it is not safe for concurrent usage.
type FifoQueue struct {
	elements   []string
	elementSet map[string]bool
}

func NewFifoQueue() *FifoQueue {
	return &FifoQueue{
		elementSet: make(map[string]bool),
		elements:   make([]string, 0),
	}
}

func (q *FifoQueue) Add(element string) {
	if q.elementSet[element] {
		return
	}
	q.elements = append(q.elements, element)
	q.elementSet[element] = true
}

// Grab returns the first element from the queue
// if queue is empty, returns an empty string
func (q *FifoQueue) Grab() string {
	if len(q.elements) == 0 {
		return ""
	}
	result := q.elements[0]
	q.elements = q.elements[1:]
	return result
}

// TaskQueue doesn't guarantee ordering but is safe for concurrent usage
// it guarantees non-duplicate elements
type TaskQueue struct {
	queue   chan string
	visited sync.Map
	queued  sync.Map
	closed  chan struct{}
}

func NewTaskQueue(seedUrl string, size int) *TaskQueue {
	q := &TaskQueue{
		queue:  make(chan string, size), // buffer size can be configurable
		closed: make(chan struct{}),
	}

	q.queued.Store(seedUrl, true)
	q.queue <- seedUrl

	return q
}

// Add adds a URL to the queue if not already queued.
// Returns true if URL was added, false if duplicate or closed.
func (q *TaskQueue) Add(url string) bool {
	if _, loaded := q.queued.LoadOrStore(url, true); loaded {
		return false
	}

	select {
	case q.queue <- url:
		return true
	case <-q.closed:
		return false
	}
}

// MarkVisited returns true if the url was NOT visited before, marking it now.
// Returns false if already visited.
func (q *TaskQueue) MarkVisited(url string) bool {
	_, loaded := q.visited.LoadOrStore(url, true)
	return !loaded
}

func (q *TaskQueue) QueuedTasks() <-chan string {
	return q.queue
}

func (q *TaskQueue) Close() {
	close(q.closed)
	close(q.queue)
}
