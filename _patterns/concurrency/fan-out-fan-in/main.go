package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func fanOut(input <-chan int, numWorkers int) []<-chan int {
	channels := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		channels[i] = worker(input)
	}
	return channels
}

func worker(input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for n := range input {
			output <- process(n)
		}
	}()
	return output
}

func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			multiplexedStream <- i
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func process(n int) int {
	// Simulate some work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	return n * n
}

func main() {
	input := make(chan int, 100)

	// Fan-out to 5 workers
	workers := fanOut(input, 5)

	// Fan-in the results
	results := fanIn(workers...)

	// Send some input
	go func() {
		for i := 0; i < 100; i++ {
			input <- i
		}
		close(input)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}
