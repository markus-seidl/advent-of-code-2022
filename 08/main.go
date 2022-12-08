package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "08/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	matrix := make([][]uint8, len(lines))
	visible := make([][]bool, len(lines))

	for i, line := range lines {
		matrix[i] = make([]uint8, len(line))
		visible[i] = make([]bool, len(line))

		for ii := 0; ii < len(line); ii++ {
			matrix[i][ii] = uint8(check2(strconv.Atoi(line[ii : ii+1])))
		}
	}

	// Part01 - Check visibility
	countVisible := 0
	for x := 1; x < len(matrix)-1; x++ {
		for y := 1; y < len(matrix[x])-1; y++ {
			height := matrix[x][y]
			v1 := checkVisibilityY(matrix, x, y+1, 1, height)
			v2 := checkVisibilityY(matrix, x, y-1, -1, height)

			v3 := checkVisibilityX(matrix, y, x+1, 1, height)
			v4 := checkVisibilityX(matrix, y, x-1, -1, height)

			temp := v1 || v2 || v3 || v4
			visible[x][y] = temp
			if temp {
				countVisible++
			}
		}
	}
	countVisible += 2*len(matrix) + 2*(len(matrix)-2)

	// Print(matrix, visible)
	println("Visible: ", countVisible)

	// Part 02 - Scenic score
	maxScenicScore := 0
	for x := 1; x < len(matrix)-1; x++ {
		for y := 1; y < len(matrix[x])-1; y++ {
			height := matrix[x][y]
			v1 := countUntilBlockingTreeY(matrix, x, y+1, 1, height)
			v2 := countUntilBlockingTreeY(matrix, x, y-1, -1, height)

			v3 := countUntilBlockingTreeX(matrix, y, x+1, 1, height)
			v4 := countUntilBlockingTreeX(matrix, y, x-1, -1, height)

			score := v1 * v2 * v3 * v4
			if score > maxScenicScore {
				maxScenicScore = score
			}
		}
	}

	println("Max Scenic Score = ", maxScenicScore)
}

func countUntilBlockingTreeY(matrix [][]uint8, x int, startIdx int, dir int, height uint8) int {
	ret := 0
	for y := startIdx; y < len(matrix[x]) && y >= 0; y += dir {
		ret++
		if matrix[x][y] >= height {
			break
		}
	}

	return ret
}

func countUntilBlockingTreeX(matrix [][]uint8, y int, startIdx int, dir int, height uint8) int {
	ret := 0
	for x := startIdx; x < len(matrix) && x >= 0; x += dir {
		ret++
		if matrix[x][y] >= height {
			break
		}
	}

	return ret
}

func checkVisibilityY(matrix [][]uint8, x int, startIdx int, dir int, height uint8) bool {
	for y := startIdx; y < len(matrix[x]) && y >= 0; y += dir {
		if matrix[x][y] >= height {
			return false
		}
	}

	return true
}

func checkVisibilityX(matrix [][]uint8, y int, startIdx int, dir int, height uint8) bool {
	for x := startIdx; x < len(matrix) && x >= 0; x += dir {
		if matrix[x][y] >= height {
			return false
		}
	}

	return true
}

func Print(matrix [][]uint8, visible [][]bool) {
	for x := 0; x < len(matrix); x++ {
		for y := 0; y < len(matrix[x]); y++ {
			if visible[x][y] {
				fmt.Printf("[%d]", matrix[x][y])
			} else {
				fmt.Printf(" %d ", matrix[x][y])
			}
		}
		println()
	}

}
