package main

import (
	"fmt"
	"time"
)

func mergeWithoutNilChannel(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)
	ch1Closed := false
	ch2Closed := false

	go func() {
		for {
			select {
			case v, open := <-ch1:
				if !open {
					ch1Closed = true
					break
				}
				ch <- v
			case v, open := <-ch2:
				if !open {
					ch2Closed = true
					break
				}
				ch <- v
			}

			if ch1Closed && ch2Closed {
				close(ch)
				return
			}
		}
	}()

	return ch
}

func mergeWithNilChannel(ch1, ch2 <-chan int) <-chan int {
	ch := make(chan int, 1)

	go func() {
		for ch1 != nil || ch2 != nil {
			select {
			case v, open := <-ch1:
				if !open {
					ch1 = nil
					break
				}
				ch <- v
			case v, open := <-ch2:
				if !open {
					ch2 = nil
					break
				}
				ch <- v
			}
		}
		close(ch)
	}()

	return ch
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Producer for ch1
	go func() {
		defer close(ch1)
		for i := 1; i <= 5; i++ {
			ch1 <- i
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Producer for ch2
	go func() {
		defer close(ch2)
		for i := 100; i <= 103; i++ {
			ch2 <- i
			time.Sleep(350 * time.Millisecond)
		}
	}()

	merged := mergeWithoutNilChannel(ch1, ch2)

	for v := range merged {
		fmt.Println(v)
	}

	fmt.Println("done")
}
