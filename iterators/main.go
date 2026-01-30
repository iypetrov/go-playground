package main

import (
	"fmt"
	"iter"
)

// Set holds a set of elements.
type Set[E comparable] struct {
    m map[E]struct{}
}

// New returns a new [Set].
func New[E comparable]() *Set[E] {
    return &Set[E]{m: make(map[E]struct{})}
}

// Add adds an element to a set.
func (s *Set[E]) Add(v E) {
    s.m[v] = struct{}{}
}

// Contains reports whether an element is in a set.
func (s *Set[E]) Contains(v E) bool {
    _, ok := s.m[v]
    return ok
}

// Union returns the union of two sets.
func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
    r := New[E]()
    // Note for/range over internal Set field m.
    // We are looping over the maps in s1 and s2.
    for v := range s1.m {
        r.Add(v)
    }
    for v := range s2.m {
        r.Add(v)
    }
    return r
}

// Option 1: Push Set elements
func (s *Set[E]) Push(f func(E) bool) {
    for v := range s.m {
        if !f(v) {
            return
        }
    }
}

func PrintAllElementsPush[E comparable](s *Set[E]) {
    s.Push(func(v E) bool {
        fmt.Println(v)
        return true
    })
}

// Option 2: Pull Set elements
// Pull returns a next function that returns each
// element of s with a bool for whether the value
// is valid. The stop function should be called
// when finished calling the next function.
func (s *Set[E]) Pull() (func() (E, bool), func()) {
    ch := make(chan E)
    stopCh := make(chan bool)

    go func() {
        defer close(ch)
        for v := range s.m {
            select {
            case ch <- v:
            case <-stopCh:
                return
            }
        }
    }()

    next := func() (E, bool) {
        v, ok := <-ch
        return v, ok
    }

    stop := func() {
        close(stopCh)
    }

    return next, stop
}

func PrintAllElementsPull[E comparable](s *Set[E]) {
    next, stop := s.Pull()
    defer stop()
    for v, ok := next(); ok; v, ok = next() {
        fmt.Println(v)
    }
}

// Option 3: Standard/Push iterator
func (s *Set[E]) All() iter.Seq[E] {
    return func(yield func(E) bool) {
        for v := range s.m {
            if !yield(v) {
                return
            }
        }
    }
}

func PrintAllElements[E comparable](s *Set[E]) {
    for v := range s.All() {
        fmt.Println(v)
    }
}

// Option 4: Pull iterator
func EqSeq[E comparable](s1, s2 iter.Seq[E]) bool {
    next1, stop1 := iter.Pull(s1)
    defer stop1()
    next2, stop2 := iter.Pull(s2)
    defer stop2()
    for {
        v1, ok1 := next1()
        v2, ok2 := next2()
        if !ok1 {
            return !ok2
        }
        if ok1 != ok2 || v1 != v2 {
            return false
        }
    }
}

// Filter returns a sequence that contains the elements
// of s for which f returns true.
func Filter[V any](f func(V) bool, s iter.Seq[V]) iter.Seq[V] {
    return func(yield func(V) bool) {
        for v := range s {
            if f(v) {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

func isEven(n int) bool {
    return n%2 == 0
}

func main() {
	s1 := New[int]()
	s1.Add(1)
	s1.Add(2)
	s1.Add(3)

	s2 := New[int]()
	s2.Add(3)
	s2.Add(4)
	s2.Add(5)

	PrintAllElements(s1)
	fmt.Printf("%t \n", EqSeq(s1.All(), s1.All()))

	evens := Filter(isEven, s1.All())
    for v := range evens {
        fmt.Println(v)
    }
}
