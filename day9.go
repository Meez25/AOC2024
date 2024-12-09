package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func dayNine() {
	input, _ := os.ReadFile("day9input.txt")
	trimmed := bytes.TrimSpace(input)
	// lines := bytes.Split(trimmed, []byte("\n"))

	inputString := string(trimmed)
	createdDisk := createDisk(inputString)
	// step1(createdDisk)
	fmt.Println(step2(createdDisk))

}

func createDisk(inputString string) []string {
	disk := make([]string, 0)

	// Create slices
	i := 0
	for j := 0; j < len(inputString); j++ {
		if j%2 == 0 {
			count, _ := strconv.Atoi(string(inputString[j]))
			slicesToAppend := make([]string, count)
			asString := strconv.Itoa(i)
			for i := range slicesToAppend {
				slicesToAppend[i] = asString
			}
			disk = append(disk, slicesToAppend...)
			i++
		} else {
			count, _ := strconv.Atoi(string(inputString[j]))
			slicesToAppend := make([]string, count)
			for i := range slicesToAppend {
				slicesToAppend[i] = "."
			}
			disk = append(disk, slicesToAppend...)
		}
	}
	return disk
}

func step1(disk []string) {

	// Compact disk
	i := 0
	j := len(disk) - 1
	for i != j {
		if disk[i] == "." && disk[j] != "." {
			disk[i] = disk[j]
			disk[j] = "."
			i++
			j--
		} else if disk[i] == "." && disk[j] == "." {
			j--
		} else {
			i++
		}
	}

	// Compute result
	result := 0
	for i, v := range disk {
		if v != "." {
			asInt, _ := strconv.Atoi(v)
			result = result + i*asInt
		}
	}
	fmt.Println(result)
}

type FileInfo struct {
	id       string
	size     int
	position int
}

func step2(disk []string) int {
	files := make([]FileInfo, 0)
	currentPosition := 0

	for currentPosition < len(disk) {
		if disk[currentPosition] != "." {
			// Found a file, collect its info
			fileID := disk[currentPosition]
			size := 1
			// Count consecutive blocks of the same ID
			for currentPosition+size < len(disk) && disk[currentPosition+size] == fileID {
				size++
			}
			files = append(files, FileInfo{
				id:       fileID,
				size:     size,
				position: currentPosition,
			})
			currentPosition += size
		} else {
			currentPosition++
		}
	}

	// Process files in reverse order (highest ID to lowest)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		// Look for the leftmost space that can fit this file
		bestPosition := -1
		currentPosition := 0

		for currentPosition < file.position {
			if disk[currentPosition] == "." {
				// Big enough ?
				spaceSize := 0
				for pos := currentPosition; pos < len(disk) && disk[pos] == "." && spaceSize < file.size; pos++ {
					spaceSize++
				}

				// Big enough !
				if spaceSize >= file.size {
					bestPosition = currentPosition
					break
				}

				// Skip past this space
				currentPosition += spaceSize
			} else {
				// Skip past this file
				currentID := disk[currentPosition]
				for currentPosition < len(disk) && disk[currentPosition] == currentID {
					currentPosition++
				}
			}
		}

		// If we found a suitable position, move the file
		if bestPosition != -1 {
			// Create temporary storage for the move
			fileContent := make([]string, file.size)
			for j := range fileContent {
				fileContent[j] = file.id
			}

			// Place dots in the original location
			for j := 0; j < file.size; j++ {
				disk[file.position+j] = "."
			}

			// Place the file in its new location
			for j := 0; j < file.size; j++ {
				disk[bestPosition+j] = file.id
			}
		}
	}

	// Calculate checksum
	checksum := 0
	for pos, value := range disk {
		if value != "." {
			fileID, _ := strconv.Atoi(value)
			checksum += pos * fileID
		}
	}

	return checksum
}
