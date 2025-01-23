package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Generator struct {
	mu     sync.Mutex
	count  int
	ticker *time.Ticker
	file   *os.File
}

func New(filename string) (*Generator, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create file: %w", err)
	}

	g := &Generator{
		count:  0,
		ticker: time.NewTicker(1 * time.Minute),
		file:   file,
	}

	go func() {
		for range g.ticker.C {
			g.mu.Lock()
			fmt.Println(g.count)
			g.count = 0
			g.mu.Unlock()

			err := g.file.Truncate(0)
			if err != nil {
				fmt.Printf("failed to truncate file: %v\n", err)
			}
			_, err = g.file.Seek(0, 0)
			if err != nil {
				fmt.Printf("Failed to seek file: %v\n", err)
			}
		}
	}()

	return g, nil
}

func (g *Generator) Increment() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.count++
	text := fmt.Sprintf("increment %d at %s\n", g.count, time.Now().Format(time.RFC3339))

	_, err := g.file.WriteString(text)
	if err != nil {
		fmt.Printf("failed to write to file: %v\n", err)
	}
}

func (g *Generator) Close() {
	g.ticker.Stop()
	g.file.Close()
}

func main() {
	generator, err := New("output.txt")
	if err != nil {
		fmt.Printf("failed to create generator: %v\n", err)
		return
	}
	defer generator.Close()

	go func() {
		for {
			generator.Increment()
		}
	}()

	select {}
}
