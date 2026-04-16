package main

import "C"

//export MyFoo
func MyFoo() {
	println("Hello World!")
}

func main() {}
