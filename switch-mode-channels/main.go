package main

import (
	"fmt"
	"os"
	"time"
)

func workerA(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Println("Stopping Worker A")
			return
		default:
			fmt.Println("Worker A")
			time.Sleep(time.Second)
		}
	}
}

func workerB(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Println("Stopping Worker B")
			return
		default:
			fmt.Println("Worker B")
			time.Sleep(time.Second)
		}
	}
}

func isWorkerB() bool {
	_, err := os.Stat("/var/lib/worker-b")
	return err == nil
}

func main() {
	stop := make(chan struct{})

	go workerA(stop)

	for {
		if isWorkerB() {
			fmt.Println("Switching to Worker B")

			close(stop)
			stop = make(chan struct{})
			go workerB(stop)
			break
		}

		time.Sleep(1 * time.Second)
	}

	select {}
}
