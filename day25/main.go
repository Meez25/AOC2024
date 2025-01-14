package main

import (
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed inputday25.txt
var inputFile []byte

type Key [5]int
type Lock [5]int

func main() {
	keys := make([]Key, 0)
	locks := make([]Lock, 0)

	for _, symbol := range bytes.Split(bytes.TrimSpace(inputFile), []byte("\n\n")) {
		table := []string{}
		for _, line := range bytes.Split(symbol, []byte("\n")) {
			table = append(table, string(line))
		}
		if table[0][0] == '.' {
			// for _, line := range table {
			// 	fmt.Println(line)
			// }
			key := Key{-1, -1, -1, -1, -1}
			for y := range table {
				for x := range table[y] {
					if table[y][x] == '#' {
						key[x]++
					}
				}
			}
			keys = append(keys, key)
		}
		if table[0][0] == '#' {
			// for _, line := range table {
			// 	fmt.Println(line)
			// }
			lock := Lock{-1, -1, -1, -1, -1}
			for y := range table {
				for x := range table[y] {
					if table[y][x] == '#' {
						lock[x]++
					}
				}
			}
			locks = append(locks, lock)
		}
	}

	match := 0

	for _, key := range keys {
		for _, lock := range locks {
			i := 0
			for i < 5 {
				if key[i]+lock[i] > 5 {
					break
				}
				i++
				if i == 5 {
					match++
				}
			}
		}
	}
	fmt.Println(match)
}
