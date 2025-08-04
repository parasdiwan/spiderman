package publish

import (
	"fmt"
	"time"
)

type Publisher interface {
	Publish(title string, lines []string) error
	PublishStats() error
	RecordError(url string, failedFor ErrType, err error) error
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
	totalErrors := 0
	for _, s := range c.erroredPages {
		totalErrors += len(s)
	}
	fmt.Println("---------------------Crawler stats ------------------")
	fmt.Println("-")
	fmt.Printf("Total time spent: %v seconds\n", totalTimeSpent)
	fmt.Println("Total pages crawled: ", c.totalPages)
	fmt.Println("Total links found: ", c.totalLinks)
	fmt.Println("Total Errors: ", len(c.erroredPages))
	if c.totalErrors > 0 {
		fmt.Println("---------------------Error stats --------------------")
		for k, s := range c.erroredPages {
			fmt.Println("[Error]: ", k)
			for _, s := range s {
				fmt.Println("-- ", s)
			}
		}
		fmt.Print("- ")
	}
	fmt.Println("-----------------------------------------------------")
	return nil
}

func (c *consoleLinkPublisher) RecordError(url string, cause ErrType, error error) error {
	c.totalErrors++
	pages, exists := c.erroredPages[cause]
	if !exists {
		pages = make([]string, 0)
	}
	pages = append(pages, url)
	c.erroredPages[cause] = pages
	return nil
}

var _ Publisher = (*consoleLinkPublisher)(nil)

func NewConsolePublisher() Publisher {
	return &consoleLinkPublisher{
		createdAt:    time.Now(),
		totalPages:   0,
		totalLinks:   0,
		erroredPages: make(map[ErrType][]string, 0),
	}
}
