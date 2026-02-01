package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// channel should be buffered
    sigs := make(chan os.Signal, 1)
	// registers the given channel to receive 
	// notifications of the specified signals
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    done := make(chan bool, 1)
    go func() {
        sig := <-sigs
        fmt.Println()
		fmt.Printf("Received signal: %v\n", sig)
        done <- true
    }()

    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}
