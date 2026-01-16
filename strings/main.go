package main

import (
	"fmt"
	"strings"
)

func concatV1(values []string) string {
	s := ""
	for _, value := range values {
		s += value
	}
	return s
}

func concatV2(values []string) string {
	sb := strings.Builder{}
	for _, value := range values {
		_, _ = sb.WriteString(value)
	}
	return sb.String()
}

func concatV3(values []string) string {
	total := 0
	for i := range values {
		total += len(values[i])
	}

	sb := strings.Builder{}
	sb.Grow(total)
	for _, value := range values {
		_, _ = sb.WriteString(value)
	}
	return sb.String()
}

func main() {
	// s := "hêllo"
	// for i, v := range s {
	// 	fmt.Printf("position %d: %c\n", i, v)
	// }
	// runes := []rune(s)
	// for i, r := range runes {
	// 	fmt.Printf("position %d: %c\n", i, r)
	// }
	// fmt.Printf("Number of runes:%d\n", utf8.RuneCountInString(s))

	// s := "xoo123oxo"
	// fmt.Println(strings.Trim(s, "xo"))       // 123
	// fmt.Println(strings.TrimLeft(s, "xo"))   // 123oxo
	// fmt.Println(strings.TrimRight(s, "xo"))  // xoo123
	// fmt.Println(strings.TrimPrefix(s, "xo")) // o123oxo
	// fmt.Println(strings.TrimSuffix(s, "xo")) // xoo123o

	// values := []string{"hello", " ", "world", "!"}
	// resultV1 := concatV1(values)
	// fmt.Println(resultV1)
	// resultV2 := concatV2(values)
	// fmt.Println(resultV2)
	// resultV3 := concatV3(values)
	// fmt.Println(resultV3)

	s1 := "Hêllo, World!"
	s2 := s1[:5]
	fmt.Println(s2)
	s3 := string([]rune(s1)[:5])
	fmt.Println(s3)
}
