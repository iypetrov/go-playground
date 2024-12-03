package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	defer file.Close()

	result := big.NewInt(0)
	scanner := bufio.NewScanner(file)
	pattern := `mul\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)`
	re := regexp.MustCompile(pattern)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) == 3 {
				num1, err := strconv.Atoi(match[1])
				if err != nil {
					fmt.Println("error converting num1:", err)
					return
				}

				num2, err := strconv.Atoi(match[2])
				if err != nil {
					fmt.Println("error converting num2:", err)
					return
				}

				mul := big.NewInt(int64(num1 * num2))
				result.Add(result, mul)
			}
		}
	}

	fmt.Println(result)
}
