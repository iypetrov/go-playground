package main

import (
	"fmt"
	"math/rand"
	"time"
)

func simulateSlowCall(sleepTime time.Duration) {
	time.Sleep(sleepTime + (time.Duration(rand.Intn(1000)) * time.Millisecond))
}

func main() {
	fmt.Println("hello parallel-tests")
}
