package main

import (
	"fmt"
	"testing"
	"time"
)

func TestApi_parallel_subtests(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("subtest_%d", i), func(t *testing.T) {
			t.Parallel()
			simulateSlowCall(1 * time.Second)
		})
	}
}
