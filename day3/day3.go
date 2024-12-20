package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	input, _ := os.ReadFile("day3input.txt")
	lines := bytes.Split(input, []byte("\n"))

	combined := make([]byte, 0, len(input))
	for _, line := range lines {
		combined = append(combined, line...)
	}

	total := processInput(combined, false)
	elapsed := time.Since(start)
	fmt.Println("Part 1:", total, "in", elapsed)

	start = time.Now()
	total = processInput(combined, true)
	elapsed = time.Since(start)
	fmt.Println("Part 2:", total, "in", elapsed)
}

func processInput(input []byte, p2 bool) int {
	total := 0
	ignore := false
	i := 0

	for i < len(input) {
		if p2 && i+4 <= len(input) {
			if bytes.Equal(input[i:i+4], []byte("do()")) {
				ignore = false
				i += 4
				continue
			}
		}
		if p2 && i+7 <= len(input) {
			if bytes.Equal(input[i:i+7], []byte("don't()")) {
				ignore = true
				i += 7
				continue
			}
		}

		if !ignore && i+4 <= len(input) {
			if bytes.Equal(input[i:i+3], []byte("mul")) && input[i+3] == '(' {
				num1, newPos := parseNumber(input, i+4)
				if newPos < len(input) && input[newPos] == ',' {
					num2, endPos := parseNumber(input, newPos+1)
					if endPos < len(input) && input[endPos] == ')' {
						total += num1 * num2
						i = endPos + 1
						continue
					}
				}
			}
		}
		i++
	}
	return total
}

func parseNumber(input []byte, start int) (int, int) {
	num := 0
	i := start
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		num = num*10 + int(input[i]-'0')
		i++
	}
	return num, i
}
