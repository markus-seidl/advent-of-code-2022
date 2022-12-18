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
