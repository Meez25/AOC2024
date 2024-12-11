package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func dayEleven() {
	input, _ := os.ReadFile("day11input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte(" "))

	stonesMap := make(map[int]int, 0)

	for _, v := range lines {
		asInt, _ := strconv.Atoi(string(v))
		stonesMap[asInt]++
	}

	count := 0
	for i := 0; i < 75; i++ {
		stonesMap = blinkMap(stonesMap)
	}

	for _, v := range stonesMap {
		count += v
	}
	fmt.Println(count)

}

func blinkMap(stonesMap map[int]int) map[int]int {
	secondMap := make(map[int]int, len(stonesMap))

	for k, v := range stonesMap {
		if k == 0 {
			secondMap[1] = v + secondMap[1]
			continue
		}
		if lenItoa(k)%2 == 0 {
			toInsert := cutIntInHalf(k)
			secondMap[toInsert[0]] = v + secondMap[toInsert[0]]
			secondMap[toInsert[1]] = v + secondMap[toInsert[1]]
			continue
		} else {
			secondMap[k*2024] += v
		}
	}
	return secondMap
}

func lenItoa(i int) int {
	return len(strconv.Itoa(i))
}

func cutIntInHalf(input int) []int {
	asString := strconv.Itoa(input)
	firstPart := asString[:len(asString)/2]
	secondPart := asString[len(asString)/2:]

	asIntFirst, _ := strconv.Atoi(firstPart)
	asIntSecond, _ := strconv.Atoi(secondPart)

	return []int{asIntFirst, asIntSecond}
}

func removeWithoutOrder(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
