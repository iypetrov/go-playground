package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func fetchURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, url)
	}
	fmt.Printf("Successfully fetched %s\n", url)
	return nil
}

func main() {
	var g errgroup.Group
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.nonexistent-url.com", // This URL will cause an error
	}

	for _, url := range urls {
		url := url
		g.Go(func() error {
			return fetchURL(url)
		})
	}

	// Wait for all goroutines to complete and check for errors
	if err := g.Wait(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All URLs fetched successfully")
	}
}
