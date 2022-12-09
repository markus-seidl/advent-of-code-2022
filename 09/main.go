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

	headPos := Pos{X: 0, Y: 0}
	tailPos := Pos{X: 0, Y: 0}

	//PrintMap(6, 6, headPos, tailPos, make([]Pos, 0))
	//println("------------")

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
			headPos.Move(dirV, &tailPos)
			tailVisited = append(tailVisited, tailPos)

			//PrintMap(6, 6, headPos, tailPos, tailVisited)
			//println("------------ ", headPos.X, headPos.Y)
		}
	}

	temp := make(map[string]bool)
	for _, pos := range tailVisited {
		key := strconv.Itoa(pos.X) + ";" + strconv.Itoa(pos.Y)

		temp[key] = true
	}

	println("Visited = ", len(temp))
}

func MoveUpDown(headPos *Pos, tailPos *Pos, dir int) {
	headBefore := *headPos
	//diffBefore := ManhattanDistance(*headPos, *tailPos)
	sameXBefore := headPos.X == tailPos.X
	//sameYBefore := headPos.Y == tailPos.Y

	// Move Head
	headPos.Y += dir

	diffAfter := ManhattanDistance(*headPos, *tailPos)
	//sameXAfter := headPos.X == tailPos.X
	sameYAfter := headPos.Y == tailPos.Y

	if diffAfter == 0 {
		// NOP
	} else if diffAfter == 1 {
		if sameYAfter {
			tailPos.Y += dir
		} else if sameXBefore {

		}
	} else if diffAfter == 2 {
		tailPos.X = headBefore.X
		tailPos.Y = headBefore.Y
	} else {
		panic("Error distAfter>2")
	}
}

func MoveRightLeft(headPos *Pos, tailPos *Pos, dir int) {
	headBefore := *headPos

	//diffBefore := ManhattanDistance(*headPos, *tailPos)
	//sameXBefore := headPos.X == tailPos.X
	//sameYBefore := headPos.Y == tailPos.Y

	// Move Head
	headPos.X += dir

	diffAfter := ManhattanDistance(*headPos, *tailPos)
	sameXAfter := headPos.X == tailPos.X
	sameYAfter := headPos.Y == tailPos.Y

	if diffAfter == 0 {
		// NOP
	} else if diffAfter == 1 {
		if sameYAfter {
			tailPos.X += dir
		} else if sameXAfter {

		}
	} else if diffAfter == 2 {
		tailPos.X = headBefore.X
		tailPos.Y = headBefore.Y
	} else {
		panic("Error")
	}
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

func PrintMap(maxX int, maxY int, head Pos, tail Pos, tailVisited []Pos) {
	for y := maxY - 1; y >= 0; y-- {
		for x := 0; x < maxX; x++ {
			if head.IsPos(x, y) && tail.IsPos(x, y) {
				print("h")
			} else if head.IsPos(x, y) {
				print("H")
			} else if tail.IsPos(x, y) {
				print("T")
			} else {
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
