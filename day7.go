package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

func daySeven() {
	start := time.Now()
	input, _ := os.ReadFile("day7input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))

	result := 0
	resultP2 := 0

	for _, line := range lines {
		toGet, numbers := parseInput(line)
		if isPossible(toGet, numbers) {
			result += toGet
		}
	}
	fmt.Println("part1 :", result)
	elapsed := time.Since(start)
	fmt.Println(elapsed)

	for _, line := range lines {
		toGet, numbers := parseInput(line)
		if isPossibleP2(toGet, numbers) {
			resultP2 += toGet
		}
	}
	fmt.Println("part2 :", resultP2)
	elapsed = time.Since(start)
	fmt.Println(elapsed)
}

func parseInput(input []byte) (int, []int) {
	numbersToArrange := make([]int, 0)
	table := bytes.Split(input, []byte(":"))
	trimmedNumbers := bytes.TrimSpace(table[1])
	numbers := bytes.Split(trimmedNumbers, []byte(" "))

	for _, number := range numbers {
		asInt, _ := strconv.Atoi(string(number))
		numbersToArrange = append(numbersToArrange, asInt)
	}
	toGet, _ := strconv.Atoi(string(table[0]))

	return toGet, numbersToArrange
}

func isPossibleP2(target int, numbers []int) bool {
	permutationSymbol := generateCombinationsPart2(len(numbers))

	for i := 0; i < len(permutationSymbol); i++ {
		result := numbers[0]
		for j := 0; j < len(numbers)-1; j++ {
			if permutationSymbol[i][j] == "+" {
				result = result + numbers[j+1]
			}
			if permutationSymbol[i][j] == "x" {
				result = result * numbers[j+1]
			}
			if permutationSymbol[i][j] == "||" {
				temp := strconv.Itoa(result) + strconv.Itoa(numbers[j+1])
				asInt, err := strconv.Atoi(temp)
				if err != nil {
					fmt.Println(err)
				}
				result = asInt
			}
		}
		if result == target {
			return true
		}
	}

	return false
}

func isPossible(target int, numbers []int) bool {
	permutationSymbol := generateCombinations(len(numbers))

	for i := 0; i < len(permutationSymbol); i++ {
		result := numbers[0]
		for j := 0; j < len(numbers)-1; j++ {
			if permutationSymbol[i][j] == "+" {
				result = result + numbers[j+1]
			}
			if permutationSymbol[i][j] == "x" {
				result = result * numbers[j+1]
			}
			if result == target {
				return true
			}
		}
	}

	return false
}

func generateCombinations(n int) [][]string {
	if n <= 1 {
		return [][]string{{}}
	}

	ops := []string{"+", "x"}
	var result [][]string

	prev := generateCombinations(n - 1)
	for _, p := range prev {
		for _, op := range ops {
			newComb := append(append([]string{}, p...), op)
			result = append(result, newComb)
		}
	}
	if n == 2 {
		fmt.Println(result)
	}
	return result
}

func generateCombinationsPart2(n int) [][]string {
	if n <= 1 {
		return [][]string{{}}
	}

	ops := []string{"+", "x", "||"}
	var result [][]string

	prev := generateCombinationsPart2(n - 1)
	for _, p := range prev {
		for _, op := range ops {
			newComb := append(append([]string{}, p...), op)
			result = append(result, newComb)
		}
	}
	return result
}
