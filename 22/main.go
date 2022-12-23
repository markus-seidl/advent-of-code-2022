package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const InputFile = "22/input.txt"

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
	Nothing TileType = 0
	Space            = 1
	Wall             = 2
)

type Pos struct {
	X int
	Y int
}

type Map struct {
	Tile       [][]TileType // x, y
	WalkedPath [][]string   // x, y
	Position   Pos
	Direction  int // 0 right; 1 down; 2 left; 3 up
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	m := Map{
		Tile:       make([][]TileType, len(lines)-2),
		WalkedPath: make([][]string, len(lines)-2),
	}

	commandList := ParseInput(lines, &m)
	// fill with empty
	maxRowLength := -1
	for _, row := range m.Tile {
		if maxRowLength < len(row) {
			maxRowLength = len(row)
		}
	}
	for i := 0; i < len(m.Tile); i++ {
		for ii := len(m.Tile[i]); ii < maxRowLength; ii++ {
			m.Tile[i] = append(m.Tile[i], Nothing)
			m.WalkedPath[i] = append(m.WalkedPath[i], "")
		}
	}

	println("Before: ")
	m.Print()
	println()

	m.FindStartPosition()
	println("Start position: ", m.Position.X, m.Position.Y)

	for _, cmd := range commandList {
		switch cmd.(type) {
		case string:
			m.Rotate(cmd.(string))
		case int:
			m.Go(cmd.(int))
		}
	}

	println("After:")
	m.Print()
	println("End position: ", m.Position.X, m.Position.Y)

	println("Part 01: ", (m.Position.X+1)*1000+4*(m.Position.Y+1)+m.Direction)
}

func (m *Map) FindStartPosition() {
	for y, t := range m.Tile[0] {
		if t != Nothing && t != Wall {
			m.Position = Pos{X: 0, Y: y}
			return
		}
	}
	panic("No start position")
}

func (m *Map) Go(steps int) {
	curPos := m.Position
	for i := 0; i < steps; i++ {
		if m.Direction == 0 || m.Direction == 2 {
			change := IIF(m.Direction == 0, 1, -1)
			curPos.Y = m.AdvanceYPosition(curPos, change)

		} else if m.Direction == 3 || m.Direction == 1 {
			change := IIF(m.Direction == 1, 1, -1)
			curPos.X = m.AdvanceXPosition(curPos, change)

		} else {
			panic("Unknown direction")
		}

		m.MarkWalked(curPos)
	}

	m.Position = curPos
}

func (m *Map) MarkWalked(p Pos) {
	var c string
	switch m.Direction {
	case 0:
		c = ">"
	case 1:
		c = "v"
	case 2:
		c = "<"
	case 3:
		c = "^"
	default:
		panic(fmt.Sprintf("Unknown direction %d", m.Direction))
	}
	m.WalkedPath[p.X][p.Y] = c
}

func (m *Map) AdvanceXPosition(p Pos, change int) int {
	oldX := p.X
	newX := oldX + change

	// border wrap
	newX = newX % len(m.Tile)
	if newX < 0 {
		newX += len(m.Tile)
	}

	// check if wall
	switch m.Tile[newX][p.Y] {
	case Wall:
		return oldX // hit a wall, do not advance
	case Nothing:
		// wrap!
		if change > 0 {
			for i := 0; i < len(m.Tile); i++ {
				if m.Tile[i][p.Y] != Nothing {
					newX = i
					break
				}
			}
		} else {
			for i := len(m.Tile) - 1; i >= 0; i-- {
				if m.Tile[i][p.Y] != Nothing {
					newX = i
					break
				}
			}
		}

		if m.Tile[newX][p.Y] == Wall {
			return oldX // hit a wall while wrapping --> cannot wrap!
		}
	}

	if m.Tile[newX][p.Y] == Nothing {
		panic("Hit nothing!")
	}

	return newX
}

func (m *Map) AdvanceYPosition(p Pos, change int) int {
	row := m.Tile[p.X]

	oldY := p.Y
	newY := oldY + change

	// border wrap
	newY = newY % len(row)
	if newY < 0 {
		newY += len(row)
	}

	// check if wall
	switch row[newY] {
	case Wall:
		return oldY // hit a wall, do not advance
	case Nothing:
		// wrap!
		if change > 0 {
			for i := 0; i < len(row); i++ {
				if row[i] != Nothing {
					newY = i
					break
				}
			}
		} else {
			for i := len(row) - 1; i >= 0; i-- {
				if row[i] != Nothing {
					newY = i
					break
				}
			}
		}
		if row[newY] == Wall {
			return oldY // hit a wall while wrapping --> cannot wrap!
		}
	}

	if row[newY] == Nothing {
		panic("Hit nothing!")
	}

	return newY
}

func IIF[T any](d bool, trueCase T, falseCase T) T {
	if d {
		return trueCase
	}
	return falseCase
}

func (m *Map) Rotate(dir string) {
	switch dir {
	case "R":
		m.Direction += 1
	case "L":
		m.Direction -= 1
	}
	if m.Direction > 3 {
		m.Direction -= 4
	}
	if m.Direction < 0 {
		m.Direction += 4
	}
}

func (m *Map) Print() {
	for rowIdx, row := range m.Tile {
		for colIdx, col := range row {
			wp := m.WalkedPath[rowIdx][colIdx]
			if wp != "" {
				print(wp)
			} else {
				switch col {
				case Nothing:
					print(" ")
				case Space:
					print(".")
				case Wall:
					print("#")
				}
			}
		}

		println()
	}
}

func ParseInput(lines []string, m *Map) []any {
	// Parse Map
	for lineIdx, line := range lines {
		if line == "" {
			break
		}

		row := make([]TileType, len(line))
		for tIdx, t := range line {
			var tile TileType
			switch t {
			case ' ':
				tile = Nothing
			case '.':
				tile = Space
			case '#':
				tile = Wall
			default:
				panic(fmt.Sprintf("Unknown type: %s", string(t)))
			}

			row[tIdx] = tile
		}

		m.Tile[lineIdx] = row
		m.WalkedPath[lineIdx] = make([]string, len(line))
	}

	lastIdx := -1
	commandsStr := lines[len(lines)-1]
	commandList := make([]any, 0)
	for idx, ch := range commandsStr {
		if unicode.IsDigit(ch) && lastIdx < 0 {
			lastIdx = idx
		} else if unicode.IsDigit(ch) {
			// nop - wait
		} else {
			if lastIdx >= 0 {
				number := commandsStr[lastIdx:idx]
				commandList = append(commandList, check2(strconv.Atoi(number)))
				lastIdx = -1
			}
			commandList = append(commandList, string(ch))
		}
	}

	if lastIdx >= 0 {
		commandList = append(commandList, check2(strconv.Atoi(commandsStr[lastIdx:len(commandsStr)])))
	}

	for _, cmd := range commandList {
		switch cmd.(type) {
		case int:
			print(fmt.Sprintf("%d", cmd))
		case string:
			print(fmt.Sprintf("%s", cmd))
		}
	}
	println()

	return commandList
}
