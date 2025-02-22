package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

//go:embed day19input.txt
var inputFile []byte

func main() {
	start := time.Now()
	availableTowels, towelsToCreate := parseFile(inputFile)

	partOneResult, partTwoResult := Solve(towelsToCreate, availableTowels)

	partOne := time.Since(start)
	fmt.Println("Part 1:", partOneResult)
	fmt.Println("Part 2:", partTwoResult)
	fmt.Println("Time: ", partOne)

}

func parseFile(inputFile []byte) ([][]byte, [][]byte) {
	table := bytes.Split(bytes.TrimSpace(inputFile), []byte("\n\n"))
	firstPart := table[0]
	secondPart := table[1]

	tableOfAvailableTowels := bytes.Split(firstPart, []byte(", "))
	tableOfTowerToCreate := bytes.Split(secondPart, []byte("\n"))

	return tableOfAvailableTowels, tableOfTowerToCreate
}

func Solve(targetTowels [][]byte, towels [][]byte) (int, int) {
	canBeSolved := 0
	totalPossibilities := 0
	maxFilteredSize := len(towels)
	filteredTowels := make([][]byte, 0, maxFilteredSize)
	dp := make([]int, 1000)

	// Filter the towels to remove the one we are sure we won't use, because their digits are not
	// in the target towel
	for _, targetTowel := range targetTowels {
		filteredTowels = filteredTowels[:0]
		targetLen := len(targetTowel)

		for _, towel := range towels {
			if len(towel) <= targetLen && bytes.Contains(targetTowel, towel) {
				filteredTowels = append(filteredTowels, towel)
			}
		}

		count := howManyPossibilities(targetTowel, filteredTowels, dp)
		if count > 0 {
			canBeSolved++
		}
		totalPossibilities += count
	}
	return canBeSolved, totalPossibilities
}

func howManyPossibilities(targetTowel []byte, towels [][]byte, dp []int) int {
	targetLen := len(targetTowel)

	// Reuse the dp array for speed ^^'
	for i := 0; i <= targetLen; i++ {
		dp[i] = 0
	}
	dp[0] = 1

	// The logic is that we first compute the first position, and build on top of the result
	for i := 0; i < targetLen; i++ {
		if dp[i] == 0 && i != 0 {
			continue
		}
		for _, towel := range towels {
			towelLen := len(towel)
			end := i + towelLen
			if end <= targetLen && bytes.Equal(towel, targetTowel[i:end]) {
				dp[end] += dp[i]
			}
		}
	}
	return dp[targetLen]
}
