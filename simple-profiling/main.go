package main

import (
	"fmt"
	"time"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.GoroutineProfile, profile.ProfilePath(".")).Stop()

	fmt.Println("Starting app...")

	for i := 0; i < 100; i++ {
		go worker(i)
	}

	time.Sleep(30 * time.Second)
	fmt.Println("Exiting app...")
}

func worker(id int) {
	ch := make(chan struct{})
	fmt.Printf("Leaking goroutine %d started\n", id)
	<-ch
}
