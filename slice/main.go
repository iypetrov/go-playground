package main

import "fmt"

type Foo struct{}

type Bar struct{}

func fooToBar(foo Foo) Bar {
	return Bar{}
}

func convertEmptySlice(foos []Foo) []Bar {
	bars := make([]Bar, 0)

	for _, foo := range foos {
		bars = append(bars, fooToBar(foo))
	}
	return bars
}

func convertGivenCapacity(foos []Foo) []Bar {
	n := len(foos)
	bars := make([]Bar, 0, n)

	for _, foo := range foos {
		bars = append(bars, fooToBar(foo))
	}
	return bars
}

func convertGivenLength(foos []Foo) []Bar {
	n := len(foos)
	bars := make([]Bar, n)

	for i, foo := range foos {
		bars[i] = fooToBar(foo)
	}
	return bars
}

func f(s []int) {
	_ = append(s, 10)
}

func listing1() {
	s := []int{1, 2, 3}

	f(s[:2])
	fmt.Println(s)
}

func listing2() {
	s := []int{1, 2, 3}
	sCopy := make([]int, 2)
	copy(sCopy, s)

	f(sCopy)
	result := append(sCopy, s[2])
	fmt.Println(result)
}

func listing3() {
	s := []int{1, 2, 3}
	f(s[:2:2])
	fmt.Println(s)
}

func main() {
	// // s1: ptr1 3 6
	// //     0 0 0 | - - -
	// s1 := make([]int, 3, 6)
	// // s2: ptr1 2 5
	// //     0 0 | - - -
	// s2 := s1[1:3]
	// fmt.Println(s1, s2)
	// 
	// // s1: ptr1 3 6
	// //     0 1 0 | - - -
	// // s2: ptr1 2 5
	// //     1 0 | - - -
	// s1[1] = 1 // s2[0] = 1
	// fmt.Println(s1, s2)
	// 
	// // s1: ptr1 3 6
	// //     0 1 0 | - - -
	// // s2: ptr1 3 5
	// //     1 0 2 | - - -
	// s2 = append(s2, 2)
	// fmt.Println(s1, s2)

	// // s1: ptr1 4 6
	// //     0 1 0 3 | - -
	// // s2: ptr1 3 5
	// //     1 0 3 | - -
	// s1 = append(s1, 3)
	// fmt.Println(s1, s2)
	// 
	// // s1: ptr1 3 6
	// //     0 1 0 3 | 4 5
	// // s2: ptr2 6 10
	// //     1 0 3 4 5 6 | - - -
	// s2 = append(s2, 4)
	// s2 = append(s2, 5)
	// s2 = append(s2, 6)
	// fmt.Println(s1, s2)
	// s1 := []int{1, 2, 3}
	// s2 := s1[1:2]
	// fmt.Println(s1, s2)
	// s3 := append(s2, 10)
	// fmt.Println(s1, s2, s3)
	listing1()
	listing2()
	listing3()
}
