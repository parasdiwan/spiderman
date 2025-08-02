## Spiderman

This is a crawler which crawls and publishes all url links for a particular website. 
It takes input of a url and finds all the links within that html page. -> also recursively crawls on the links on that website and prints the links for that page as well.

### Usage:

Make sure you have the build of the project in the current directory

```shell
  ./spider http://monzo.com
```

if you want to export the results to a file

```shell
  ./spider https://www.duckduckgo.com > results.txt
```

- Make build
```shell
  make build
```
it creates a binary in the current directory called `spider`
which can be run as mentioned above.

- To run tests
```shell
  make test 
```

**It's possible the tests could fail because it hits wesbites which could be down or could block the IP**

## Architecture & Design Decisions

### Project Structure
``` 
spiderman/
├── main.go                 # CLI entry point
├── crawl/
│   ├── crawler.go         # Main crawler logic
│   ├── crawler_test.go    # Crawler tests
│   ├── queue.go           # Queue implementations
│   ├── filters/
│   │   ├── filters.go     # Link filtering logic
│   │   └── filters_test.go
│   ├── links/
│   │   ├── parser.go      # HTML link extraction
│   │   ├── parser_test.go
│   │   └── link_extractors.go
│   └── http/
│       ├── fetcher.go     # HTTP client wrapper
│       └── fetcher_test.go
└── publish/
    ├── publisher.go       # Output handling
    └── test_helpers.go    # Test utilities
```

### Core Components

The system is built around several key abstractions that promote separation of concerns and extensibility:

#### 1. **Crawler** (`crawl/crawler.go`)
- **Responsibility**: Orchestrates the crawling process and manages recursion
- **Design Choice**: Supports both sequential (`Crawl()`) and parallel (`CrawlParallel()`) modes
- **Trade-off**: Parallel mode sacrifices some determinism for performance gains

#### 2. **Parser** (`crawl/links/parser.go`)
- **Responsibility**: Extracts links from HTML pages
- **Design Choice**: Pluggable extractor system (supports `<a href>`, `<link href>`, `<area href>`)
- **Trade-off**: Could extract more link types, but focused on most common use cases

#### 3. **Publisher** (`publish/publisher.go`)
- **Responsibility**: Handles output formatting and statistics
- **Design Choice**: Interface-based design allows for different output formats
- **Current Implementation**: Console output with crawl statistics
- **Extensibility**: Easy to add file, database, or web socket publishers

#### 4. **Filters** (`crawl/filters/filters.go`)
- **Responsibility**: Determines which links should be crawled
- **Design Choice**: Chain of responsibility pattern with composable filters
- **Current Filters**:
    - `InternalLink`: Only crawl same-domain links
    - `NoEmpty`: Skip empty/invalid links
    - `NoFragment`: Skip anchor fragments (#section)
    - `NoFile`: Skip file downloads (.pdf, .jpg, etc.)
    - `NoMailLink`: Skip mailto: links
    - `NoTelephone`: Skip tel: links

### Dependencies
- **Core**: Standard library only (net/http, html parser)
- **Testing**: for assertions `github.com/stretchr/testify`
- **HTML Parsing**: for robust HTML parsing `golang.org/x/net/html`


## Key Trade-offs & Assumptions
### **Internal Links Only**
- **Decision**: Only crawl links within the same domain
- **Rationale**: Prevents infinite crawling and respects website boundaries
- **Implementation**: and treated as same domain `www.example.com``example.com`

### **No Robots.txt Compliance**
- **Decision**: Does not check robots.txt
- **Trade-off**: Simpler implementation vs. web etiquette

### **Limited Retry Logic**
- **Decision**: No automatic retries for failed requests
- **Rationale**: Prevents worker threads from blocking on slow/failing endpoints
- **Alternative**: Could add exponential backoff/circuit breakers for production use

### **Memory vs. Performance**
- **Decision**: Keep all discovered URLs in memory for deduplication
- **Trade-off**: Fast lookups vs. memory usage on large sites

### **File Type Filtering**
- **Decision**: Skip common file types (.pdf, .jpg, .js, etc.)
- **Rationale**: Focus on HTML pages that might contain more links
- **Implementation**: Extensible list of file extensions

### **Queue** (`crawl/queue.go`)
- **Sequential**: Simple FIFO queue with deduplication
- **Parallel**: Thread-safe channel-based queue with visited tracking
- **Trade-off**: Parallel queue uses more memory but enables concurrency

### Key Concurrency Features:
- **Worker Pool**: Fixed number of goroutines processing URLs
- **WaitGroup Coordination**: Ensures all work completes before termination
- **Thread-Safe Queue**: Uses `sync.Map` for visited tracking and channels for work distribution
- **Graceful Shutdown**: Proper channel closing prevents goroutine leaks


### Future improvements:
If I had more time, I would have done the following:
- **Improve filtering**: There could be more filters to be added both to the publisher and for crawled links.
  - Add filters according to `robots.txt`
  - Some query params can be added to the url to filter the results
- **Improve error handling**:
  - Collect errors separately from the valid urls
  - Adding retries or other resilience mechanisms
  - Graceful handling of invalid cases
  - Revisit http status code like 20x, 30x to check if some responses are valid.
- **Expand publisher**:
    - Publisher is not thread-safe. It should be.
    - Collect more statistics on the crawling process
    - `publisher.RecordError` function is not integrated with the crawler.
- **Improve concurrency**:
  - Find the perfect ratio between number of workers and url channel buffer size.
  - Add a limit to the number of workers to avoid memory issues.
  - Add a limit to the number of urls to be crawled so we don't run into memory issues.
- **Place limits when websites end-up being too big**: Add a limit to the amount of urls to be crawled so we don't run into memory issues.
- Parser can be improved to become more generic.
