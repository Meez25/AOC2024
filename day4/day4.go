package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	input, _ := os.ReadFile("day4input.txt")
	lines := bytes.Split(input, []byte("\n"))

	lines = lines[:len(lines)-1]

	match := 0

	for i := 0; i < len(lines); i++ {
		for y := 0; y < len(lines[i]); y++ {
			if lines[i][y] == 'X' {
				// Check left
				if y > 2 && lines[i][y-1] == 'M' && lines[i][y-2] == 'A' && lines[i][y-3] == 'S' {
					match++
				}
				// Check right
				if y < len(lines[i])-3 && lines[i][y+1] == 'M' && lines[i][y+2] == 'A' && lines[i][y+3] == 'S' {
					match++
				}
				// Check up
				if i > 2 && lines[i-1][y] == 'M' && lines[i-2][y] == 'A' && lines[i-3][y] == 'S' {
					match++
				}
				// Check down
				if i < len(lines)-3 && lines[i+1][y] == 'M' && lines[i+2][y] == 'A' && lines[i+3][y] == 'S' {
					match++
				}
				// Check up right
				if i > 2 && y < len(lines[i])-3 && lines[i-1][y+1] == 'M' && lines[i-2][y+2] == 'A' && lines[i-3][y+3] == 'S' {
					match++
				}
				// Check up left
				if i > 2 && y > 2 && lines[i-1][y-1] == 'M' && lines[i-2][y-2] == 'A' && lines[i-3][y-3] == 'S' {
					match++
				}
				// Check down right
				if i < len(lines)-3 && y < len(lines[i])-3 && lines[i+1][y+1] == 'M' && lines[i+2][y+2] == 'A' && lines[i+3][y+3] == 'S' {
					match++
				}
				// Check down left
				if i < len(lines)-3 && y > 2 && lines[i+1][y-1] == 'M' && lines[i+2][y-2] == 'A' && lines[i+3][y-3] == 'S' {
					match++
				}
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", match, "in", elapsed)
	matchPart2 := 0
	start = time.Now()

	for i := 0; i < len(lines); i++ {
		for y := 0; y < len(lines[i]); y++ {
			if lines[i][y] == 'A' {
				// Check up right
				if i > 0 && y > 0 && i < len(lines)-1 && y < len(lines[i])-1 {
					if (lines[i-1][y-1] == 'M' && lines[i+1][y+1] == 'S') || (lines[i-1][y-1] == 'S' && lines[i+1][y+1] == 'M') {
						if (lines[i+1][y-1] == 'M' && lines[i-1][y+1] == 'S') || (lines[i+1][y-1] == 'S' && lines[i-1][y+1] == 'M') {
							matchPart2++
						}
					}
				}
			}
		}
	}

	elapsed = time.Since(start)
	fmt.Println("Part 2 :", matchPart2, "in", elapsed)
}
