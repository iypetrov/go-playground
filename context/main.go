package main

import (
	"context"
	"fmt"
	"time"
)

type requestIDKeyType string

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, requestIDKeyType("requestID"), "abc-123")

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	done := make(chan struct{})

	go worker(ctx, done)

	select {
	case <-done:
		cancel()
		fmt.Println("main done successfully")
	case <-ctx.Done():
		fmt.Println("main done with error:", ctx.Err())
	}
}

func worker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	reqID := ctx.Value(requestIDKeyType("requestID"))
	fmt.Println("worker started, requestID:", reqID)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("worker finished work")
	case <-ctx.Done():
		fmt.Println("worker canceled:", ctx.Err())
	}
}
