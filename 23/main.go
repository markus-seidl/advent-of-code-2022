package main

import (
	"math"
	"os"
	"strconv"
	"strings"
)

const InputFile = "23/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

type TileType int

const (
	Ground TileType = 0
	Elf             = 1
)

type Pos struct {
	X int
	Y int
}

type Map struct {
	Elves          []Pos // x = vvvvv; y --->
	FirstDirection int   // 0 North; 1 South; 2 West; 3 East
	ElvesIndex     map[Pos]bool
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	m := Map{
		Elves: make([]Pos, 0),
	}

	for x, line := range lines {
		for y, c := range line {
			if c == '#' {
				elf := Pos{X: x, Y: y}
				m.Elves = append(m.Elves, elf)
			}
		}
	}

	m.Print()
	println()

	for i := 1; i <= 1000; i++ {
		println("== Round ", i, " ==: Direction ", m.FirstDirection)
		elfMoved := m.Move()
		//m.Print()

		if !elfMoved {
			break
		}
	}

	println("== Result ==")
	m.Print()

	println("Result: ", m.CountEmptyTiles())
}

func (m *Map) Move() bool {
	previousElves := len(m.Elves)
	m.UpdateElvesIndex()

	// Propose move for every elf
	proposedMoves := make([]struct {
		Elf    Pos
		NewElf Pos
	}, 0)

	for _, elf := range m.Elves {
		newElf := m.ProposeMoveElf(elf)
		proposedMoves = append(proposedMoves, struct {
			Elf    Pos
			NewElf Pos
		}{Elf: elf, NewElf: newElf})
	}

	// find collisions
	duplicates := make(map[Pos]int)
	for _, m := range proposedMoves {
		duplicates[m.NewElf]++
	}

	elfMoved := false
	// Move every elf
	newElfList := make([]Pos, 0)
	for _, p := range proposedMoves {
		if duplicates[p.NewElf] == 1 {
			newElfList = append(newElfList, p.NewElf)

			if p.Elf != p.NewElf {
				elfMoved = true
			}

		} else {
			newElfList = append(newElfList, p.Elf)
		}
	}

	if previousElves != len(newElfList) {
		panic("WE LOST AN ELF! " + strconv.Itoa(previousElves) + " " + strconv.Itoa(len(newElfList)))
	}

	m.Elves = newElfList

	m.NextDirection()

	return elfMoved
}

func (m *Map) ProposeMoveElf(pos Pos) Pos {
	dir := m.FirstDirection

	// Is there an elf around me?
	if !m.IsElfOn(Pos{X: pos.X, Y: pos.Y - 1}) && !m.IsElfOn(Pos{X: pos.X, Y: pos.Y + 1}) && !m.IsElfOn(Pos{X: pos.X - 1, Y: pos.Y}) && !m.IsElfOn(Pos{X: pos.X + 1, Y: pos.Y}) && (            // x-y
	!m.IsElfOn(Pos{X: pos.X - 1, Y: pos.Y - 1}) && !m.IsElfOn(Pos{X: pos.X - 1, Y: pos.Y + 1}) && !m.IsElfOn(Pos{X: pos.X + 1, Y: pos.Y + 1}) && !m.IsElfOn(Pos{X: pos.X + 1, Y: pos.Y - 1})) { // diagonal
		return pos
	}

	// Yes, move into the first free direction
	for i := m.FirstDirection; i < m.FirstDirection+4; i++ {
		dir = i % 4
		switch dir {
		case 0: // North (up)
			newX := pos.X - 1
			if !m.IsElfOn(Pos{X: newX, Y: pos.Y}) && !m.IsElfOn(Pos{X: newX, Y: pos.Y + 1}) && !m.IsElfOn(Pos{X: newX, Y: pos.Y - 1}) {
				return Pos{X: newX, Y: pos.Y}
			}
		case 1: // South (down)
			newX := pos.X + 1
			if !m.IsElfOn(Pos{X: newX, Y: pos.Y}) && !m.IsElfOn(Pos{X: newX, Y: pos.Y + 1}) && !m.IsElfOn(Pos{X: newX, Y: pos.Y - 1}) {
				return Pos{X: newX, Y: pos.Y}
			}
		case 2: // West (left)
			newY := pos.Y - 1
			if !m.IsElfOn(Pos{X: pos.X, Y: newY}) && !m.IsElfOn(Pos{X: pos.X - 1, Y: newY}) && !m.IsElfOn(Pos{X: pos.X + 1, Y: newY}) {
				return Pos{X: pos.X, Y: newY}
			}
		case 3: // East (right)
			newY := pos.Y + 1
			if !m.IsElfOn(Pos{X: pos.X, Y: newY}) && !m.IsElfOn(Pos{X: pos.X - 1, Y: newY}) && !m.IsElfOn(Pos{X: pos.X + 1, Y: newY}) {
				return Pos{X: pos.X, Y: newY}
			}
		default:
			panic("Error in the matrix")
		}
	}

	return pos // no free direction
}

func (m *Map) UpdateElvesIndex() {
	m.ElvesIndex = make(map[Pos]bool)

	for _, elf := range m.Elves {
		m.ElvesIndex[elf] = true
	}
}

func (m *Map) IsElfOn(pos Pos) bool {
	return m.ElvesIndex[pos]
	//for _, elf := range m.Elves {
	//  if elf == pos {
	//    return true
	//  }
	//}
	//return false
}

func (m *Map) NextDirection() {
	m.FirstDirection = (m.FirstDirection + 1) % 4
}

func (m *Map) CountEmptyTiles() int {
	minX := math.MaxInt64
	maxX := -math.MaxInt64
	minY := math.MaxInt64
	maxY := -math.MaxInt64
	elfMap := make(map[Pos]bool)

	for _, elf := range m.Elves {
		minX = Min(minX, elf.X)
		maxX = Max(maxX, elf.X)
		minY = Min(minY, elf.Y)
		maxY = Max(maxY, elf.Y)
		elfMap[elf] = true
	}

	sum := 0
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if !elfMap[Pos{X: x, Y: y}] {
				sum++
			}
		}
	}

	return sum
}

func (m *Map) Print() {
	minX := math.MaxInt64
	maxX := -math.MaxInt64
	minY := math.MaxInt64
	maxY := -math.MaxInt64
	elfMap := make(map[Pos]bool)

	for _, elf := range m.Elves {
		minX = Min(minX, elf.X)
		maxX = Max(maxX, elf.X)
		minY = Min(minY, elf.Y)
		maxY = Max(maxY, elf.Y)
		elfMap[elf] = true
	}

	println("(", minX, minY, ") -> (", maxX, maxY, ")")

	for x := minX - 3; x <= maxX+3; x++ {
		for y := minY - 3; y <= maxY+3; y++ {
			if elfMap[Pos{X: x, Y: y}] {
				print("#")
			} else {
				print(".")
			}
		}
		println()
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
