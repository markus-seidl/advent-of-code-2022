package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const InputFile = "09/example.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

type Pos struct {
	X int
	Y int
}

func (p Pos) IsPos(x int, y int) bool {
	return x == p.X && y == p.Y
}

func (p *Pos) Add(pos Pos) {
	p.X += pos.X
	p.Y += pos.Y
}

func (head *Pos) Move(dir Pos, tailPos *Pos) {
	tempNewPos := *head
	tempNewPos.X += dir.X
	tempNewPos.Y += dir.Y

	if ManhattanDistance(tempNewPos, *tailPos) > 1 {
		// Tail moves first one step to H
		if head.Y == tailPos.Y && dir.X != 0 {
			// same row
			x := tailPos.X - head.X
			if math.Abs(float64(x)) < 1 {
				// on top of each other == NOP
			} else if math.Signbit(float64(x)) {
				tailPos.X += 1
			} else {
				tailPos.X += -1
			}
		} else if head.X == tailPos.X && dir.Y != 0 {
			// same column
			y := tailPos.Y - head.Y
			if math.Abs(float64(y)) < 1 {
				// on top of each other == NOP
			} else if math.Signbit(float64(y)) {
				tailPos.Y += 1
			} else {
				tailPos.Y += -1
			}
		} else if head.IsPos(tailPos.X, tailPos.Y) {
			// overlapping
		} else if ManhattanDistance(*tailPos, *head) == 2 {
			// too far away; move diagonally into the dir of the head
			diffX := head.X - tailPos.X
			diffY := head.Y - tailPos.Y
			tailPos.X += 1 * Sign(diffX)
			tailPos.Y += 1 * Sign(diffY)
		} else {
			//
		}
	}

	head.X += dir.X
	head.Y += dir.Y
}

func Sign(a int) int {
	if a > 0 {
		return 1
	}
	if a < 0 {
		return -1
	}
	return 0
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	const tailLength = 9

	headPos := Pos{X: 0, Y: 0}
	var tailPoss []Pos
	for i := 0; i < tailLength; i++ {
		tailPoss = append(tailPoss, Pos{X: 0, Y: 0})
	}

	PrintMap(6, 6, headPos, tailPoss, make([]Pos, 0))
	println("------------")

	var tailVisited []Pos
	for _, line := range lines {
		dirParts := strings.Split(line, " ")
		dir := dirParts[0]
		size := check2(strconv.Atoi(dirParts[1]))
		fmt.Printf("> %s %d\n", dir, size)

		// Move Part 01
		var dirV Pos
		switch dir {
		case "R":
			dirV = Pos{X: 1, Y: 0}
		case "L":
			dirV = Pos{X: -1, Y: 0}
		case "U":
			dirV = Pos{X: 0, Y: 1}
		case "D":
			dirV = Pos{X: 0, Y: -1}
		}

		for t := 0; t < size; t++ {
			headPos.Move(dirV, &tailPoss[len(tailPoss)-1])
			for i := 1; i < len(tailPoss); i++ {
				tailPoss[i-1].Move(Pos{X: 0, Y: 0}, &tailPoss[i])
			}

			//tailVisited = append(tailVisited, tailPos)

			PrintMap(6, 6, headPos, tailPoss, tailVisited)
			println("------------ ", headPos.X, headPos.Y)
		}
	}

	temp := make(map[string]bool)
	for _, pos := range tailVisited {
		key := strconv.Itoa(pos.X) + ";" + strconv.Itoa(pos.Y)

		temp[key] = true
	}

	println("Visited = ", len(temp))
}

func ManhattanDistance(pos1 Pos, pos2 Pos) int {
	return absDiffInt(pos1.X, pos2.X) + absDiffInt(pos1.Y, pos2.Y)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func PrintMap(maxX int, maxY int, head Pos, tails []Pos, tailVisited []Pos) {
	for y := maxY - 1; y >= 0; y-- {
		for x := 0; x < maxX; x++ {
			alreadySet := false
			for i, tail := range tails {
				if head.IsPos(x, y) && tail.IsPos(x, y) {
					print("h")
					alreadySet = true
					break
				} else if head.IsPos(x, y) {
					print("H")
					alreadySet = true
					break
				} else if tail.IsPos(x, y) {
					print(strconv.Itoa(i))
					alreadySet = true
					break
				}
			}

			if !alreadySet {
				visited := false
				for _, p := range tailVisited {
					if p.X == x && p.Y == y {
						visited = true
						break
					}
				}

				if visited {
					print("#")
				} else {
					print(".")
				}
			}
		}
		println()
	}

}
