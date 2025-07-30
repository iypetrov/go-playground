package main

import (
	"fmt"
	"sync"
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
	s := &SwappableInt{}

	go func() {
		time.Sleep(5 * time.Second)
		s.Swap(5)
	}()

	for {
		v, ok := s.Get()
		if ok {
			fmt.Println("Ready! Value is", v)
			break
		}
		fmt.Println("Waiting for value...")
		time.Sleep(500 * time.Millisecond)
	}
}
