package main

import "fmt"

type Config struct {
	Port int
}

func New(addr string, cfg Config) {}

func main() {
	fmt.Println("hello config")
}
