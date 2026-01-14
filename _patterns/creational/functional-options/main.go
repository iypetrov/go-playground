package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type options struct {
	maxConn *int
}

type Option func(options *options) error

func WithMaxConn(maxConn int) Option {
	return func(options *options) error {
		if maxConn < 0 {
			return errors.New("max connection should be positive")
		}
		options.maxConn = &maxConn
		return nil
	}
}

func New(addr string, opts ...Option) (*sql.DB, error){
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	// here we can implement logic for filling missing values or sth like this
	return nil, nil
}

func main() {
	fmt.Println("hello functional-options")
}
