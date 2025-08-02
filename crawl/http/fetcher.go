package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FetchResult struct {
	StatusCode int
	Body       io.ReadCloser // caller must close Body if non-nil
	Location   string        // redirect target URL if applicable
	Err        error
}

type Fetcher struct {
	Client *http.Client
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch makes a GET request for the given URL and returns a FetchResult.
// This version does NOT perform retries; it returns immediately with what it gets
// http call from fetch tries to mimic chrome browser to not get `202` status code in some cases.
//
// ## improvements
// - retrying is left out atm to not block threads.
func (f *Fetcher) Fetch(rawUrl string) FetchResult {
	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return FetchResult{Err: err}
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36") // Or mimic your version
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FetchResult{Err: err}
	}

	switch resp.StatusCode {
	case http.StatusOK, 201, 203, 204, 206:
		return FetchResult{
			StatusCode: resp.StatusCode,
			Body:       resp.Body,
			Err:        nil,
		}
	case 301, 302, 303, 307, 308:
		location := resp.Header.Get("Location")
		resp.Body.Close()
		if location == "" {
			return FetchResult{
				StatusCode: resp.StatusCode,
				Err:        fmt.Errorf("redirect status %d without Location header", resp.StatusCode),
			}
		}
		return FetchResult{
			StatusCode: resp.StatusCode,
			Location:   location,
			Err:        nil,
		}
	case 202:
		err = errors.New("received HTTP 202 Accepted - processing not complete")
	case 429:
		err = errors.New("received HTTP 429 Too Many Requests")
	default:
		err = errors.New(fmt.Sprintf("failed with status %d", resp.StatusCode))
	}
	resp.Body.Close()
	return FetchResult{
		StatusCode: resp.StatusCode,
		Err:        err,
	}
}
