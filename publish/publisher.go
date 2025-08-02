package publish

import (
	"fmt"
	"time"
)

type Publisher interface {
	Publish(title string, lines []string) error
	PublishStats() error
	RecordError(failedFor string, err ErrType) error
}

type consoleLinkPublisher struct {
	createdAt    time.Time
	totalPages   int
	totalLinks   int
	totalErrors  int
	erroredPages map[ErrType][]string
}

func (c *consoleLinkPublisher) Publish(title string, lines []string) error {
	c.totalPages++
	c.totalLinks += len(lines)
	fmt.Println("Links found on: ", title)
	for _, s := range lines {
		fmt.Println(" - " + s)
	}
	return nil
}

func (c *consoleLinkPublisher) PublishStats() error {
	totalTimeSpent := time.Since(c.createdAt).Seconds()
	fmt.Println("---------------------Crawler stats ------------------")
	fmt.Println("-")
	fmt.Printf("Total time spent: %v seconds\n", totalTimeSpent)
	fmt.Println("Total pages crawled: ", c.totalPages)
	fmt.Println("Total links found: ", c.totalLinks)
	fmt.Println("-----------------------------------------------------")
	return nil
}

func (c *consoleLinkPublisher) RecordError(failedFor string, errType ErrType) error {
	c.totalErrors++
	pages, exists := c.erroredPages[errType]
	if !exists {
		pages = make([]string, 0)
	}
	pages = append(pages, failedFor)
	c.erroredPages[errType] = pages
	return nil
}

var _ Publisher = (*consoleLinkPublisher)(nil)

func NewConsolePublisher() Publisher {
	return &consoleLinkPublisher{
		createdAt:    time.Now(),
		totalPages:   0,
		totalLinks:   0,
		erroredPages: make(map[ErrType][]string),
	}
}

type ErrType string

const (
	ErrType404 = ErrType("Error 404")
)
