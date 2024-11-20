package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Received cancellation signal\n", id)
			return
		default:
			fmt.Printf("Worker %d: Doing work\n", id)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure all paths cancel the context to release resources

	// Start multiple workers
	for i := 0; i < 3; i++ {
		go worker(ctx, i)
	}

	// Simulate some work in the main goroutine
	time.Sleep(7 * time.Second)

	fmt.Println("Main: Exiting")
}
