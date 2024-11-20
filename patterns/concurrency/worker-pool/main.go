package main

import (
	"fmt"
	"math/rand"
	"time"
)

func workerPool(numWorkers int, jobs <-chan int, results chan<- int) {
	for i := 0; i < numWorkers; i++ {
	 go worker(jobs, results)
	}
   }
   
   func worker(jobs <-chan int, results chan<- int) {
	for j := range jobs {
	 results <- process(j)
	}
   }
   
   func process(job int) int {
	// Simulate some work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	return job * 2
   }
   
   func main() {
	numJobs := 100
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
   
	// Start the worker pool
	workerPool(5, jobs, results)
   
	// Send jobs
	for i := 0; i < numJobs; i++ {
	 jobs <- i
	}
	close(jobs)
   
	// Collect results
	for i := 0; i < numJobs; i++ {
	 result := <-results
	 fmt.Printf("Result: %d\n", result)
	}
	fmt.Println("Done...")
   }