package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IsReportSafe(report []int) bool {
	diffs := make([]int, len(report)-1)
	for i, v := range report {
		if i == 0 {
			continue
		}

		prev := report[i-1]
		diffs[i-1] = (v - prev)
	}

	isSafe := true
	for i, diff := range diffs {
		abs := Abs(diff)
		if abs < 1 || abs > 3 {
			isSafe = false
			break
		}

		if i == 0 {
			continue
		}

		prev := diffs[i-1]
		if diff < 0 && prev > 0 {
			isSafe = false
			break
		} else if diff > 0 && prev < 0 {
			isSafe = false
			break
		}
	}

	return isSafe
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	defer file.Close()

	cnt := 0
	data := make(map[int][]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				fmt.Println("error converting string to int")
				return
			}
			data[cnt] = append(data[cnt], num)
		}
		cnt++
	}

	safeReports := 0
	for _, report := range data {
		if IsReportSafe(report) {
			safeReports++
		}
	}

	fmt.Println(safeReports)
}
