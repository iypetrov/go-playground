package main

import "fmt"

func main() {
	p := PointerReceiver{count: 0}
	p.Increment()
	fmt.Println(p.count)
}
