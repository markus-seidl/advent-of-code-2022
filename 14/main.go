package main

import (
	"os"
	"strconv"
	"strings"
)

const InputFile = "14/input.txt"

type Material int

const (
	Air        Material = 0
	Rock                = 1
	SandMoving          = 2
	SandStill           = 3
)

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

	startPos := Pos{X: 500, Y: 0}
	maxMapSize := Pos{X: 1000, Y: 1000}

	m := Map{
		Type: CreateMap(maxMapSize),
		//  minX, minY, maxX, maxY
		Range: []int{startPos.X, startPos.Y, startPos.X, startPos.Y},
	}

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		println(line)
		paths := strings.Split(line, " -> ")
		for ei := 1; ei < len(paths); ei++ {
			from := ParsePos(paths[ei-1])
			to := ParsePos(paths[ei])

			m.RockLine(from, to)
		}
	}

	// part 02
	DrawPart02Line(&m)
	// /part 02

	m.Draw(startPos)

	sandCount := 0
	for true {
		leftMap := SimulateSingleSand(&m, startPos)
		//m.Draw(startPos)
		if leftMap {
			m.Draw(startPos)
			println("SandCount:", sandCount)
			return
		}

		sandCount++
	}
}

func DrawPart02Line(m *Map) {
	_, _, _, maxY := m.Range[0], m.Range[1], m.Range[2], m.Range[3]
	m.RockLine(Pos{5, maxY + 2}, Pos{995, maxY + 2})
}

func SimulateSingleSand(m *Map, startPos Pos) bool {

	currentPos := startPos
	if !m.IsFree(startPos) {
		return true // part 02
	}

	for true {
		// fall down one step
		tryNextPos := currentPos
		tryNextPos.Y += 1

		// hit "rock" bottom
		minX, minY, maxX, maxY := m.Range[0], m.Range[1], m.Range[2], m.Range[3]
		if tryNextPos.X < minX || tryNextPos.X > maxX || tryNextPos.Y < minY || tryNextPos.Y > maxY {
			m.Type[tryNextPos.X][tryNextPos.Y] = SandMoving
			return true // we move outside the map
		}

		if m.IsFree(tryNextPos) {
			currentPos = tryNextPos
			continue
		}

		// fall down bottom left
		tryNextPos.X -= 1
		if m.IsFree(tryNextPos) {
			currentPos = tryNextPos
			continue
		}

		// fall down bottom right
		tryNextPos.X += 2 // negate the -1 before
		if m.IsFree(tryNextPos) {
			currentPos = tryNextPos
			continue
		}

		break
	}

	m.Type[currentPos.X][currentPos.Y] = SandStill
	return false
}

func CreateMap(maxMapSize Pos) [][]Material {
	ret := make([][]Material, maxMapSize.X)

	for i := 0; i < len(ret); i++ {
		ret[i] = make([]Material, maxMapSize.Y)
	}

	return ret
}

func ParsePos(s string) Pos {
	c := strings.Split(s, ",")
	return Pos{
		X: check2(strconv.Atoi(c[0])),
		Y: check2(strconv.Atoi(c[1])),
	}
}

type Map struct {
	Type  [][]Material
	Range []int // minX, minY, maxX, maxY
}

type Pos struct {
	X int
	Y int
}

func (m *Map) RockLine(from Pos, to Pos) {
	if to.X != from.X && to.Y != from.Y {
		panic("Diagonal!")
	}

	if to.X == from.X {
		for ri := Min(to.Y, from.Y); ri <= Max(to.Y, from.Y); ri++ {
			m.Type[to.X][ri] = Rock
		}
	}
	if to.Y == from.Y {
		for ri := Min(to.X, from.X); ri <= Max(to.X, from.X); ri++ {
			m.Type[ri][to.Y] = Rock
		}
	}

	// Update range
	minX, minY, maxX, maxY := Range(from, to)
	m.Range[0] = Min(m.Range[0], minX)
	m.Range[1] = Min(m.Range[1], minY)
	m.Range[2] = Max(m.Range[2], maxX)
	m.Range[3] = Max(m.Range[3], maxY)
}

func (m *Map) Draw(startPos Pos) {
	minX, minY, maxX, maxY := m.Range[0], m.Range[1], m.Range[2], m.Range[3]

	for x := minX; x <= maxX; x++ {
		if x == startPos.X {
			print("v")
		}
		print(" ")
	}
	println()

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			switch m.Type[x][y] {
			case Air:
				print(".")
			case Rock:
				print("#")
			case SandMoving:
				print("+")
			case SandStill:
				print("o")
			default:
				print("?")
			}
			print()
		}

		println()
	}
}

func (m *Map) IsFree(pos Pos) bool {
	return m.Type[pos.X][pos.Y] == Air
}

// Range minX, minY, maxX, maxY
func Range(a Pos, b Pos) (int, int, int, int) {
	return Min(a.X, b.X), Min(a.Y, b.Y), Max(a.X, b.X), Max(a.Y, b.Y)
}

func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}
