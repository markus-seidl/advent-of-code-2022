package main

import (
	"os"
	"strings"
)

const InputFile = "24/example.txt"

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
	Tiles      [][]Tile
	CurrentPos Pos
}

type Tile struct {
	IsWall    bool
	Blizzards []Pos
}

type Pos struct {
	X int
	Y int
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	m := Map{
		Tiles: make([][]Tile, len(lines)),
	}

	for x, line := range lines {
		row := make([]Tile, len(line))
		for y, c := range line {
			if c == '#' {
				row[y].IsWall = true
			} else if c == '.' {
				row[y].IsWall = false // "nop"
			} else if c == '>' {
				row[y].Blizzards = append(row[y].Blizzards, Pos{X: 0, Y: 1})
			} else if c == '<' {
				row[y].Blizzards = append(row[y].Blizzards, Pos{X: 0, Y: -1})
			} else if c == '^' {
				row[y].Blizzards = append(row[y].Blizzards, Pos{X: -1, Y: 0})
			} else if c == 'v' {
				row[y].Blizzards = append(row[y].Blizzards, Pos{X: 1, Y: 0})
			} else {
				panic("Unknown char")
			}
		}
		m.Tiles[x] = row
	}

	println("Initial state:")
	m.Print()
	println()

	m.CurrentPos = Pos{X: 0, Y: 1}
	if !m.CanBe(m.CurrentPos) {
		panic("Invalid starting position")
	}
	goal := Pos{
		X: len(m.Tiles) - 1,
		Y: len(m.Tiles[0]) - 2,
	}
	if !m.CanBe(goal) {
		panic("Invalid goal position")
	}

	solution := Solve(m, goal)

	solution.WorldMap.Print()
}

func (m *Map) Print() {
	for _, row := range m.Tiles {
		for _, tile := range row {
			if tile.IsWall {
				print("#")
			} else if len(tile.Blizzards) == 0 {
				print(".")
			} else if len(tile.Blizzards) > 1 {
				print(len(tile.Blizzards))
			} else {
				switch tile.Blizzards[0] {
				case Pos{X: 0, Y: 1}:
					print(">")
				case Pos{X: 0, Y: -1}:
					print("<")
				case Pos{X: -1, Y: 0}:
					print("^")
				case Pos{X: 1, Y: 0}:
					print("v")
				default:
					panic("Unknown blizzard direction")
				}
			}
		}
		println()
	}
}

func (m *Map) Copy() Map {
	newTiles := make([][]Tile, len(m.Tiles))
	for x, row := range m.Tiles {
		newTiles[x] = make([]Tile, len(row))
		for y, tile := range row {
			tempBlizzards := make([]Pos, len(tile.Blizzards))
			copy(tempBlizzards, tile.Blizzards)

			newTiles[x][y] = tile
			newTiles[x][y].Blizzards = tempBlizzards
		}
	}
	return Map{
		Tiles:      newTiles,
		CurrentPos: m.CurrentPos,
	}
}

func (m *Map) SimulateBlizzards() {
	// copy m.Tiles without blizzards
	newTiles := make([][]Tile, len(m.Tiles))
	for x, row := range m.Tiles {
		newTiles[x] = make([]Tile, len(row))
		for y, tile := range row {
			tile.Blizzards = make([]Pos, 0)
			newTiles[x][y] = tile
		}
	}

	for x, row := range m.Tiles {
		for y, tile := range row {
			for _, blizzard := range tile.Blizzards {
				newPos := Add(Pos{X: x, Y: y}, blizzard)
				if m.Tiles[newPos.X][newPos.Y].IsWall {
					if newPos.X == 0 && newPos.Y != 0 {
						// top wall, wrap to bottom
						newPos.X = len(m.Tiles) - 2
					} else if newPos.X == len(m.Tiles)-1 && newPos.Y != 0 {
						// bottom wall, wrap to top
						newPos.X = 1
					} else if newPos.X != 0 && newPos.Y == 0 {
						// left wall, wrap to right
						newPos.Y = len(m.Tiles[0]) - 2
					} else if newPos.X != 0 && newPos.Y == len(m.Tiles[0])-1 {
						// right wall, wrap to left
						newPos.Y = 1
					} else {
						panic("Unhandled direction")
					}
				}

				// should work because we always modify another blizzard in another tile
				newTiles[newPos.X][newPos.Y].Blizzards = append(newTiles[newPos.X][newPos.Y].Blizzards, blizzard)
			}
		}
	}

	m.Tiles = newTiles
}

func (m *Map) SimulateBlizzardsIntoNewMap() Map {
	// copy m.Tiles without blizzards
	newTiles := make([][]Tile, len(m.Tiles))
	for x, row := range m.Tiles {
		newTiles[x] = make([]Tile, len(row))
		copy(newTiles[x], row)

		for _, tile := range row {
			tile.Blizzards = nil // reset blizzards, will be filled in next loop
		}
	}

	for x, row := range m.Tiles {
		for y, tile := range row {
			for _, blizzard := range tile.Blizzards {
				newPos := Add(Pos{X: x, Y: y}, blizzard)
				WrapNewBlizzardPosition(m, &newPos)

				// should work because we always modify another blizzard in another tile
				newTiles[newPos.X][newPos.Y].Blizzards = append(newTiles[newPos.X][newPos.Y].Blizzards, blizzard)
			}
		}
	}

	return Map{
		Tiles:      newTiles,
		CurrentPos: m.CurrentPos,
	}
}

func WrapNewBlizzardPosition(m *Map, newPos *Pos) {
	if m.Tiles[newPos.X][newPos.Y].IsWall {
		if newPos.X == 0 && newPos.Y != 0 {
			// top wall, wrap to bottom
			newPos.X = len(m.Tiles) - 2
		} else if newPos.X == len(m.Tiles)-1 && newPos.Y != 0 {
			// bottom wall, wrap to top
			newPos.X = 1
		} else if newPos.X != 0 && newPos.Y == 0 {
			// left wall, wrap to right
			newPos.Y = len(m.Tiles[0]) - 2
		} else if newPos.X != 0 && newPos.Y == len(m.Tiles[0])-1 {
			// right wall, wrap to left
			newPos.Y = 1
		} else {
			panic("Unhandled direction")
		}
	}
}

func NegateDir(dir Pos) Pos {
	return Pos{X: -dir.X, Y: -dir.Y}
}

func Add(a Pos, b Pos) Pos {
	return Pos{X: a.X + b.X, Y: a.Y + b.Y}
}
