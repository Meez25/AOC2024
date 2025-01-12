package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"time"
)

//go:embed inputday22.txt
var inputFile []byte

func main() {
	sum := 0
	start := time.Now()
	for _, line := range bytes.Split(bytes.TrimSpace(inputFile), []byte("\n")) {
		prices := make([]int, 2001)
		startingSecret, err := strconv.Atoi(string(line))
		prices[0] = startingSecret
		if err != nil {
			panic(err)
		}
		for i := 1; i < 2001; i++ {
			startingSecret = nextSecret(startingSecret)
			prices[i] = startingSecret
		}
		sum += startingSecret
		pricesChanges := convertPricesToPriceChange(prices)
		possibilities := findBestCombinaisons(prices, pricesChanges)
		fmt.Println(possibilities)
	}
	fmt.Println("Part 1", sum, "in", time.Since(start))
}

func findBestCombinaisons(prices, pricesChanges []int) []int {
	prices = prices[1:]
	fmt.Println(prices[0], pricesChanges[0])
	return []int{}
}

func mix(secret, toMixWith int) int {
	return secret ^ toMixWith
}

func prune(secret int) int {
	return secret % 16777216
}

func step1(secret int) int {
	secret = mix(secret, secret*64)
	secret = prune(secret)
	return secret
}

func step2(secret int) int {
	temp := secret / 32
	secret = mix(temp, secret)
	return prune(secret)
}

func step3(secret int) int {
	value := secret * 2048
	secret = mix(value, secret)
	return prune(secret)
}

func nextSecret(secret int) int {
	step1 := step1(secret)
	step2 := step2(step1)
	step3 := step3(step2)
	return step3
}

func convertPricesToPriceChange(price []int) []int {
	priceChange := make([]int, len(price))
	for i := 1; i < len(price)-1; i++ {
		priceChange[i-1] = price[i]%10 - price[i-1]%10
	}
	return priceChange
}
