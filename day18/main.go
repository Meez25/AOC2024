package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"slices"
	"strings"
	"time"
)

//go:embed day18input.txt
var inputFile string

type Node struct {
	x    int
	y    int
	path []Position
}

type Position struct {
	x int
	y int
}

func newNode(x, y int) *Node {
	path := make([]Position, 0)
	return &Node{x: x, y: y, path: path}
}

func (n *Node) AddPosition(position Position) {
	n.path = append(n.path, position)
}

func main() {
	start := time.Now()
	grid := make([][]string, 71)
	for i := range grid {
		grid[i] = make([]string, 71)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	for _, coord := range strings.Split(strings.TrimSpace(inputFile), "\n")[:1024] {
		parts := strings.Split(coord, ",")
		x, y := 0, 0
		fmt.Sscanf(parts[0]+","+parts[1], "%d,%d", &x, &y)
		grid[y][x] = "#"
	}

	// Printing grid
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}

	// BFS to find exit
	successNodes := BFS(grid)
	path := successNodes[0].path

	err := generateGIF(grid, path, inputFile)
	if err != nil {
		fmt.Printf("Error generating GIF: %v\n", err)
	}

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", len(successNodes[0].path), "in", elapsed)

	// Part 2, doing BFS for each new corrupted memory
	firstPosition := Part2(inputFile, path, grid)

	elapsed = time.Since(start)
	fmt.Println("Part 2 :", firstPosition, elapsed)

	// Printing end grid
	for y := range grid {
		for x := range grid[y] {
			for _, position := range path {
				if y == position.y && x == position.x {
					grid[position.y][position.x] = "O"
				}
			}
		}
	}
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}
}

func Part2(inputFile string, path []Position, grid [][]string) Position {
	pathAsPosition := make([]Position, len(path))

	var firstPositionToMatch Position
	var allInputFalling []Position

	for i, pos := range path {
		position := Position{x: pos.x, y: pos.y}
		pathAsPosition[i] = position
	}

	for _, coord := range strings.Split(strings.TrimSpace(inputFile), "\n") {
		parts := strings.Split(coord, ",")
		x, y := 0, 0
		fmt.Sscanf(parts[0]+","+parts[1], "%d,%d", &x, &y)
		position := Position{x: x, y: y}
		allInputFalling = append(allInputFalling, position)
	}

	i := 1024
	for i := i; i < len(allInputFalling); i++ {
		grid[allInputFalling[i].y][allInputFalling[i].x] = "#"
		if slices.Contains(path, allInputFalling[i]) {
			result := BFS(grid)
			if len(result) > 0 {
				path = result[0].path
			} else {
				firstPositionToMatch = allInputFalling[i]
				break
			}
		}
	}
	return firstPositionToMatch
}

func BFS(grid [][]string) []Node {
	startNode := newNode(0, 0)
	visited := make(map[Position]bool)
	queue := make([]*Node, 0)
	successNode := make([]Node, 0)

	queue = append(queue, startNode)
	visited[Position{x: 0, y: 0}] = true
	startNode.path = append(startNode.path, Position{x: 0, y: 0})

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v.x == len(grid[0])-1 && v.y == len(grid)-1 {
			successNode = append(successNode, *v)
			continue
		}
		// Explore 4 sides
		// Up
		if v.y > 0 {
			up := grid[v.y-1][v.x]
			if up != "#" {
				upNode := newNode(v.x, v.y-1)
				upNode.path = make([]Position, len(v.path))
				copy(upNode.path, v.path)
				nextPosition := Position{x: upNode.x, y: upNode.y}
				upNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, upNode)
					visited[nextPosition] = true
				}
			}
		}
		// Down
		if v.y < len(grid)-1 {
			down := grid[v.y+1][v.x]
			if down != "#" {
				downNode := newNode(v.x, v.y+1)
				downNode.path = make([]Position, len(v.path))
				copy(downNode.path, v.path)
				nextPosition := Position{x: downNode.x, y: downNode.y}
				downNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, downNode)
					visited[nextPosition] = true
				}
			}
		}
		// Right
		if v.x < len(grid[0])-1 {
			right := grid[v.y][v.x+1]
			if right != "#" {
				rightNode := newNode(v.x+1, v.y)
				rightNode.path = make([]Position, len(v.path))
				copy(rightNode.path, v.path)
				nextPosition := Position{x: rightNode.x, y: rightNode.y}
				rightNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, rightNode)
					visited[nextPosition] = true
				}
			}
		}
		// Left
		if v.x > 0 {
			left := grid[v.y][v.x-1]
			if left != "#" {
				leftNode := newNode(v.x-1, v.y)
				leftNode.path = make([]Position, len(v.path))
				copy(leftNode.path, v.path)
				nextPosition := Position{x: leftNode.x, y: leftNode.y}
				leftNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, leftNode)
					visited[nextPosition] = true
				}
			}
		}
	}
	return successNode
}

func drawPoint(img *image.Paletted, x, y, size int, c color.Color) {
	for dy := -size; dy <= size; dy++ {
		for dx := -size; dx <= size; dx++ {
			if dx*dx+dy*dy <= size*size {
				img.Set(x+dx, y+dy, c)
			}
		}
	}
}

func generateGIF(grid [][]string, initialPath []Position, inputFile string) error {
	// Constants for the GIF
	cellSize := 8
	width := len(grid[0]) * cellSize
	height := len(grid) * cellSize
	delay := 1
	pathSize := 3
	endpointSize := 4
	criticalPause := 100 // Long pause before path gets cut

	palette := []color.Color{
		color.White,
		color.Black,
		color.RGBA{255, 0, 0, 255},
		color.RGBA{64, 192, 87, 255},
		color.RGBA{250, 82, 82, 255},
		color.RGBA{255, 165, 0, 255},
	}

	var frames []*image.Paletted
	var delays []int

	baseImg := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	drawBackground(baseImg, grid, cellSize, palette)
	currentPath := initialPath
	drawFullPath(baseImg, currentPath, cellSize, palette, pathSize, endpointSize)

	initialFrame := image.NewPaletted(baseImg.Bounds(), palette)
	copy(initialFrame.Pix, baseImg.Pix)
	frames = append(frames, initialFrame)
	delays = append(delays, 50)

	currentGrid := make([][]string, len(grid))
	for i := range grid {
		currentGrid[i] = make([]string, len(grid[i]))
		copy(currentGrid[i], grid[i])
	}

	var allFallingBlocks []Position
	for _, coord := range strings.Split(strings.TrimSpace(inputFile), "\n")[1024:] {
		parts := strings.Split(coord, ",")
		x, y := 0, 0
		fmt.Sscanf(parts[0]+","+parts[1], "%d,%d", &x, &y)
		allFallingBlocks = append(allFallingBlocks, Position{x: x, y: y})
	}

	foundCriticalBlock := false
	for _, fallingBlock := range allFallingBlocks {
		// If this block will cut the path, add extra pause frames
		if slices.Contains(currentPath, fallingBlock) {
			result := BFS(currentGrid)
			if len(result) == 0 {
				// This is the critical block - add pause frames
				pauseFrame := image.NewPaletted(baseImg.Bounds(), palette)
				copy(pauseFrame.Pix, baseImg.Pix)
				frames = append(frames, pauseFrame)
				delays = append(delays, criticalPause)
				foundCriticalBlock = true
			}
		}

		highlightFrame := image.NewPaletted(baseImg.Bounds(), palette)
		copy(highlightFrame.Pix, baseImg.Pix)
		drawBlock(highlightFrame, fallingBlock.x, fallingBlock.y, cellSize, palette[5])
		frames = append(frames, highlightFrame)
		delays = append(delays, delay)

		wallFrame := image.NewPaletted(baseImg.Bounds(), palette)
		copy(wallFrame.Pix, baseImg.Pix)
		drawBlock(wallFrame, fallingBlock.x, fallingBlock.y, cellSize, palette[1])
		frames = append(frames, wallFrame)
		delays = append(delays, delay)

		currentGrid[fallingBlock.y][fallingBlock.x] = "#"

		if slices.Contains(currentPath, fallingBlock) {
			newFrame := image.NewPaletted(baseImg.Bounds(), palette)
			drawBackground(newFrame, currentGrid, cellSize, palette)

			result := BFS(currentGrid)
			if len(result) > 0 {
				currentPath = result[0].path
				drawFullPath(newFrame, currentPath, cellSize, palette, pathSize, endpointSize)
			}

			copy(baseImg.Pix, newFrame.Pix)
			frames = append(frames, newFrame)
			delays = append(delays, delay)
		} else {
			copy(baseImg.Pix, wallFrame.Pix)
		}

		// If we've found and processed the critical block, break
		if foundCriticalBlock {
			break
		}
	}

	// Add final pause
	finalFrame := image.NewPaletted(baseImg.Bounds(), palette)
	copy(finalFrame.Pix, frames[len(frames)-1].Pix)
	frames = append(frames, finalFrame)
	delays = append(delays, 50)

	f, err := os.Create("path_visualization.gif")
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, &gif.GIF{
		Image: frames,
		Delay: delays,
	})
}

// Helper functions remain the same except for drawFullPath which needs to accept size parameters

func drawSquare(img *image.Paletted, x, y, cellSize int, c color.Color) {
	// Draw a filled square centered at x,y
	size := cellSize - 2 // Slightly smaller than cell size to leave a small gap
	startX := x - size/2
	startY := y - size/2

	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			img.Set(startX+dx, startY+dy, c)
		}
	}
}

func drawFullPath(img *image.Paletted, path []Position, cellSize int, palette []color.Color, pathSize, endpointSize int) {
	if len(path) == 0 {
		return
	}

	// Draw path with solid squares
	for _, pos := range path {
		drawSquare(img, pos.x*cellSize+cellSize/2, pos.y*cellSize+cellSize/2, cellSize, palette[2])
	}

	// Draw start and end points as larger squares
	drawSquare(img, path[0].x*cellSize+cellSize/2, path[0].y*cellSize+cellSize/2, cellSize+2, palette[3])
	drawSquare(img, path[len(path)-1].x*cellSize+cellSize/2, path[len(path)-1].y*cellSize+cellSize/2, cellSize+2, palette[4])
}

// Helper function to draw a block (wall or falling block)
func drawBlock(img *image.Paletted, x, y, cellSize int, c color.Color) {
	for dy := 0; dy < cellSize; dy++ {
		for dx := 0; dx < cellSize; dx++ {
			img.Set(x*cellSize+dx, y*cellSize+dy, c)
		}
	}
}

func drawBackground(img *image.Paletted, grid [][]string, cellSize int, palette []color.Color) {
	// Fill background
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, palette[0])
		}
	}

	// Draw walls
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == "#" {
				drawBlock(img, x, y, cellSize, palette[1])
			}
		}
	}
}
