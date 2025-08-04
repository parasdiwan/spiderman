package main

import (
	"fmt"
	"net/url"
	"os"
	"spiderman/crawl"
	"spiderman/publish"
	"strconv"
	"strings"
)

func main() {
	input, numWorkers := sanitizeInputs()

	crawler := crawl.NewCrawler(input, publish.NewConsolePublisher())
	var err error
	if numWorkers == 1 {
		err = crawler.Crawl()
	} else {
		err = crawler.CrawlParallel(numWorkers) // use the optional argument here
	}
	if err != nil {
		fmt.Println("Spider had issues spidering")
		fmt.Printf("[Error]: %v\n", err)
		return
	}
}

func sanitizeInputs() (string, int) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: spider <base_website_link> [num_workers::optional]")
		os.Exit(1)
	}

	input := os.Args[1]
	// validate input
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Usage: spider <base_website_link> [num_workers]")
		os.Exit(1)
	}

	valid := IsValidHTTPLink(input)
	if !valid {
		fmt.Println("Not a valid link \nUsage: spider <base_website_link> [num_workers]")
		os.Exit(1)
	}

	// default number of workers
	numWorkers := 20
	if len(os.Args) >= 3 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n <= 0 {
			fmt.Println("num_workers must be a positive integer")
			os.Exit(1)
		}
		numWorkers = n
	}
	return input, numWorkers
}

func IsValidHTTPLink(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(u.Scheme)
	if (scheme == "http" || scheme == "https") && u.Host != "" {
		return true
	}
	return false
}
