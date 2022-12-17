package main

import (
	"os"
	"strings"
)

const InputFile = "17/input.txt"
const ChamberWidth int = 7

type Space int

const (
	Air        Space = 0
	RockMoving       = 1
	RockStill        = 2
)

const (
	RockShape0 = "####"
	RockShape1 = ".#..###..#.."
	RockShape2 = "..#...#.###."
	RockShape3 = "#...#...#...#..."
	RockShape4 = "##..##.."
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

type Map struct {
	Map [][]Space
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")
	jetForecast := lines[0]

	const simulationLength int = 2023
	m := Map{
		Map: make([][]Space, 0),
	}

	blocks := [][][]Space{ExtractBlock(RockShape0), ExtractBlock(RockShape1), ExtractBlock(RockShape2), ExtractBlock(RockShape3), ExtractBlock(RockShape4)}
	blockWidths := []int{4, 3, 3, 1, 2}
	blockIdx := 0

	//
	// Rocks fall Y up!
	// Y  X --->
	// |
	// |
	// V
	//
	m.ExtendMap(4)

	jetStep := 0
	for rockCount := 1; rockCount < simulationLength; rockCount++ {
		println("Simulate rock: ", rockCount)

		// spawn block
		newRock := blocks[blockIdx%len(blocks)]
		newRockWdith := blockWidths[blockIdx%len(blocks)]
		blockIdx++

		// extend map
		lastBlockIdx := m.LastBlockOrFloorIdx()
		if rockCount != 1 {
			m.ExtendMap(lastBlockIdx + len(newRock) + 4 - len(m.Map))
		}

		rockTopLeft := Pos{X: 2, Y: len(m.Map) - 1}

		for true {
			jetStep = jetStep % len(jetForecast)

			// Apply Jet
			//println(jetStep, string(jetForecast[jetStep]), jetForecast)
			moveLeft := string(jetForecast[jetStep]) == "<"
			jetStep++

			m.PrintWMR(rockTopLeft, newRock)

			if moveLeft {
				temp := Pos{X: rockTopLeft.X - 1, Y: rockTopLeft.Y}
				if m.IsFree(temp, newRock, newRockWdith) {
					rockTopLeft = temp
				}
			} else {
				temp := Pos{X: rockTopLeft.X + 1, Y: rockTopLeft.Y}
				if m.IsFree(temp, newRock, newRockWdith) {
					rockTopLeft = temp
				}
			}

			// Fall one step
			fallPos := Pos{X: rockTopLeft.X, Y: rockTopLeft.Y - 1}
			if m.HitRock(fallPos, newRock, newRockWdith) {
				m.PutBlock(rockTopLeft, newRock, RockStill)
				m.Print()
				break
			} else {
				rockTopLeft = fallPos // ok to fall down one step
				m.PrintWMR(rockTopLeft, newRock)
			}
			m.Print()
		}
	}

	println("-----------------------")
	m.Print()
	println("-----------------------")
	println("len(m.Map) = ", len(m.Map), "lastRockIdx", m.LastBlockOrFloorIdx())
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func ExtractBlock(rockShape string) [][]Space {
	const rockWidth = 4
	lines := len(rockShape) / rockWidth
	ret := make([][]Space, lines)
	for l := 0; l < lines; l++ {
		ret[l] = make([]Space, rockWidth)
		for c := 0; c < rockWidth; c++ {
			var s Space
			temp := rockShape[l*rockWidth+c : l*rockWidth+c+1]
			switch temp {
			case "#":
				s = RockMoving
			default:
				s = Air
			}

			ret[l][c] = s
		}
	}
	return ret
}

func (m *Map) PrintWMR(topLeft Pos, movingRock [][]Space) {
	c := m.Copy()
	c.PutBlock(topLeft, movingRock, RockMoving)

	c.Print()
}

func (m *Map) Copy() Map {
	duplicate := make([][]Space, len(m.Map))
	for i := range m.Map {
		duplicate[i] = make([]Space, len(m.Map[i]))
		copy(duplicate[i], m.Map[i])
	}

	return Map{Map: duplicate}
}

func (m *Map) Print() {

	println("|vvvvvvv| -- top")

	for y := len(m.Map) - 1; y >= 0; y-- {
		print("|")
		for x := 0; x < len(m.Map[y]); x++ {
			switch m.Map[y][x] {
			case Air:
				print(".")
			case RockStill:
				print("#")
			case RockMoving:
				print("@")
			}
		}
		println("|")
	}

	println("|-------| -- bottom")
}

func (m *Map) LastBlockOrFloorIdx() int {
	for y := len(m.Map) - 1; y >= 0; y-- {
		if m.HasRock(y) {
			return y
		}
	}
	return 0
}

func (m *Map) PutBlock(topLeft Pos, block [][]Space, asSpace Space) {
	for y := 0; y < len(block); y++ {
		for x := 0; x < len(block[y]); x++ {
			mY := topLeft.Y - y
			mX := topLeft.X + x

			if block[y][x] != Air {
				if m.Map[mY][mX] != Air && block[y][x] != Air {
					panic("unexpected hit")
				}
				m.Map[mY][mX] = asSpace
			}
		}
	}
}

func (m *Map) HitRock(topLeft Pos, block [][]Space, blockWidth int) bool {
	if !m.IsFree(topLeft, block, blockWidth) {
		return true
	}

	////if topLeft.Y-1 >= 0 {
	//y := len(block) - 1
	//for x := 0; x < len(block[y]); x++ {
	//  b := block[y][x]
	//
	//  if b != Air {
	//    // check if there is a rock under this rock
	//    if m.Map[topLeft.Y-y][topLeft.X+x] != Air {
	//      return true
	//    }
	//  }
	//}
	////}

	return false
}

func (m *Map) IsFree(topLeft Pos, block [][]Space, blockWidth int) bool {
	if topLeft.Y+1-len(block) < 0 {
		return false // would be below bottom
	}

	if topLeft.X+blockWidth > ChamberWidth {
		return false // block would be outside of chamber
	}
	if topLeft.X < 0 {
		return false // block would be outside of chamber
	}

	for y := 0; y < len(block); y++ {
		for x := 0; x < len(block[y]); x++ {
			b := block[y][x]

			if b != Air && m.Map[topLeft.Y-y][topLeft.X+x] != Air {
				return false
			}
		}
	}

	return true
}

type Pos struct {
	X int
	Y int
}

func (m *Map) HasRock(y int) bool {
	for x := 0; x < len(m.Map[y]); x++ {
		if m.Map[y][x] != Air {
			return true
		}
	}
	return false
}

func (m *Map) ExtendMap(lines int) {
	if lines == 0 {
		return
	}

	if lines < 0 {
		m.Map = m.Map[:len(m.Map)-(-lines)]
	}

	if lines > 0 {
		for i := 0; i < lines; i++ {
			m.Map = append(m.Map, make([]Space, ChamberWidth))
		}
	}
}
