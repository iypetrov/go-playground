package main

import (
	"fmt"
	"sync"
	"time"
)

func workerChannel(done chan bool) {
    fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")
    done <- true
}

func workerWaitGroup(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")
}

func workerWaitGroupLatest() {
	fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")
}

// only accepts a channel for sending values
func ping(pings chan<- string, msg string) {
    pings <- msg
}

// accepts one channel for receives (pings) and a second for sends (pongs).
func pong(pings <-chan string, pongs chan<- string) {
    msg := <-pings
    pongs <- msg
}

func main() {
	// unbuffered channel
    msgsUnbuffered := make(chan string)
    go func() { msgsUnbuffered <- "foo" }()
    msg := <-msgsUnbuffered
    fmt.Println(msg)

	// buffered channel
	msgsBuffered := make(chan string, 2)
	msgsBuffered <- "hello"
	msgsBuffered <- "world"
	fmt.Println(<-msgsBuffered)
	fmt.Println(<-msgsBuffered)

	// sync
	// if we want to wait just for a single goroutine, we can use a channel
	done := make(chan bool)
    go workerChannel(done)
    <-done

	// if we want to wait for multiple goroutines, we can use sync.WaitGroup
	var wg sync.WaitGroup

	wg.Add(1)          
	go workerWaitGroup(&wg)     
	wg.Wait()          

	// after Go 1.25, we can use wg.Go
	wg.Go(func() {
        workerWaitGroupLatest()
    })
	wg.Wait()

	// directions
    pings := make(chan string, 1)
    pongs := make(chan string, 1)
    ping(pings, "passed message")
    pong(pings, pongs)
    fmt.Println(<-pongs)

	// select
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for range 2 {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }

	// timeout
    c3 := make(chan string, 1)
    go func() {
        time.Sleep(2 * time.Second)
        c3 <- "result 3"
    }()

    select {
    case res := <-c3:
        fmt.Println(res)
    case <-time.After(3 * time.Second):
        fmt.Println("timeout 3")
    }

	// non-blocking select
	messages := make(chan string)
	signals := make(chan bool)

	// non-blocking receive
	// if there is a message in messages it will be received
	// otherwise the default case is executed
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}
	
	// non-blocking send
	// if messages channel is not full the message will be sent
	// otherwise the default case is executed
	txt := "hi"
	select {
	case messages <- txt:
		fmt.Println("sent message", txt)
	default:
		fmt.Println("no message sent")
	}

	// can use multiple cases above the default clause
	// to implement a multi-way non-blocking select
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
	
	// closing channel
    jobs := make(chan int, 5)
    done = make(chan bool)

    go func() {
        for {
			// this special 2-value form of receive, the more value will be
			// false if jobs has been closed and all values in the channel 
			// have already been received
            j, more := <-jobs
            if more {
                fmt.Println("received job", j)
            } else {
                fmt.Println("received all jobs")
                done <- true
                return
            }
        }
    }()

    for j := 1; j <= 3; j++ {
        jobs <- j
        fmt.Println("sent job", j)
    }
    close(jobs)
    fmt.Println("sent all jobs")

    <-done

    _, ok := <-jobs
    fmt.Println("received more jobs:", ok)

	// for/range over a channel
	queue := make(chan string, 2)
    queue <- "one"
    queue <- "two"
    close(queue)

    for elem := range queue {
        fmt.Println(elem)
    }
}
