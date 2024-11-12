package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetchGitHub(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN")),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data for github user with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func getYear(dateStr string) (int, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return 0, err
	}
	return date.Year(), nil
}
