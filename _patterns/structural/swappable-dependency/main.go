package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type SwappableInt struct {
	mu    sync.RWMutex
	value int
	ready bool
}

func (s *SwappableInt) Get() (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value, s.ready
}

func (s *SwappableInt) Swap(v int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = v
	s.ready = true
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &SwappableInt{}

	go func() {
		select {
		case <-time.After(5 * time.Second):
			s.Swap(5)
		case <-ctx.Done():
			fmt.Println("Swap goroutine exiting early due to signal")
			return
		}
	}()

	for {
		v, ok := s.Get()
		if ok {
			fmt.Println("Ready! Value is", v)
			break
		}

		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("Waiting for value...")
		case <-ctx.Done():
			fmt.Println("Main loop exiting early due to signal")
			return
		}
	}
}
