package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	defer file.Close()

	listOne := make([]int, 0)
	listTwo := make([]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) >= 2 {
			valOne, err1 := strconv.Atoi(parts[0])
			valTwo, err2 := strconv.Atoi(parts[1])

			if err1 == nil && err2 == nil {
				listOne = append(listOne, valOne)
				listTwo = append(listTwo, valTwo)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file")
	}

	slices.Sort(listOne)
	slices.Sort(listTwo)

	dist := 0
	for i := 0; i < len(listOne); i++ {
		dist += Abs(listOne[i] - listTwo[i])
	}
	fmt.Println(dist)
}
