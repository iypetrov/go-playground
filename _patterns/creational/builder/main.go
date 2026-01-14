package main

import (
	"errors"
	"fmt"
)

type DatabasePool struct {
	Addr string
	MaxConn int
}

type DatabasePoolBuilder struct {
	maxConn *int
}

func (b *DatabasePoolBuilder) MaxConn(maxConn int) *DatabasePoolBuilder {
	b.maxConn = &maxConn
	return b
}

func (b *DatabasePoolBuilder) Build() (DatabasePool, error) {
	dbpool := DatabasePool{}
	if b.maxConn == nil {
		dbpool.MaxConn = 10
	} else if(*b.maxConn < 0) {
		return DatabasePool{}, errors.New("max connection should be positive")
	} else {
		dbpool.MaxConn = *b.maxConn
	}
	return dbpool, nil
}

func main() {
	fmt.Println("hello builder")
}
