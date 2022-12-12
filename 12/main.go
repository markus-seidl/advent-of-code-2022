package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "12/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type HeightMap struct {
	Height [][]int
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	hm := HeightMap{
		Height: make([][]int, len(lines)),
	}

	var start Pos
	var end Pos
	for x, line := range lines {
		elements := make([]int, len(line))
		for y, element := range line {
			elements[y] = int(element)

			if element == int32('S') {
				start = Pos{X: x, Y: y}
			}
			if element == int32('E') {
				end = Pos{X: x, Y: y}
			}
		}
		hm.Height[x] = elements
	}
	hm.Print()

	fmt.Printf("Start: %v\n", start)
	fmt.Printf("End: %v\n", end)

	// Part 1
	//path := a_star(hm, start, end)
	//println("Found path with", len(path)-1, "steps")
	//
	//tempMap := map[Pos]bool{}
	//for _, p := range path {
	//  print(p.X, "/", p.Y, "/", hm.At(p), " - > ")
	//  tempMap[p] = true
	//}
	//println()
	//
	//for x, line := range hm.Height {
	//  for y, _ := range line {
	//    if _, ok := tempMap[Pos{X: x, Y: y}]; ok {
	//      print("X")
	//    } else {
	//      print(".")
	//    }
	//  }
	//  println()
	//}

	// Part 2
	// search all a for the shortest path
	minSteps := 10000
	var minStepsStartPoint Pos
	var minPath []Pos
	for x, line := range hm.Height {
		for y, v := range line {
			if v == 'a' {
				start = Pos{X: x, Y: y}
				path := a_star(hm, start, end)
				if len(path) == 0 {
					continue
				}

				if minSteps > len(path) {
					println("Starting point:", start.String(), "needs", len(path)-1, "steps")
					minSteps = len(path)
					minStepsStartPoint = start
					minPath = path
				}
			}
		}
	}

	println("Minimum:", len(minPath)-1, "from", minStepsStartPoint.String())
}

func a_star(hm HeightMap, start Pos, end Pos) []Pos {
	openSet := map[Pos]float32{}
	openSet[start] = 1 // value unused (set...)

	cameFrom := map[Pos]Pos{}

	gScore := map[Pos]float32{}
	gScore[start] = 0

	fScore := map[Pos]float32{}
	fScore[start] = h(start, end)

	for len(openSet) > 0 {
		current := lowestScore(openSet, fScore)
		if current == end {
			println("Found solution with gScore", withDefaultInf(gScore, current), "fScore", fScore[current])
			return reconstructPath(cameFrom, current, start)
		}

		delete(openSet, current)
		for _, neighbor := range neighbors(hm, current) {
			tentative_gScore := withDefaultInf(gScore, current) + 1 // + d(current, neighbor) == 1 because we only allow direct neighbors and only care about length
			if tentative_gScore < withDefaultInf(gScore, neighbor) {
				// good path
				cameFrom[neighbor] = current
				gScore[neighbor] = tentative_gScore
				fScore[neighbor] = tentative_gScore + h(current, end)

				if _, ok := openSet[neighbor]; !ok {
					openSet[neighbor] = 1 // value unused (set...)
				}
			}
		}
	}

	return []Pos{}
}

func lowestScore(candidates map[Pos]float32, scores map[Pos]float32) Pos {
	minScore := float32(1000)
	var minScorePos Pos
	for k := range candidates {
		if minScore > scores[k] {
			minScorePos = k
			minScore = scores[k]
		}
	}

	return minScorePos
}

func reconstructPath(cameFrom map[Pos]Pos, current Pos, start Pos) []Pos {
	path := []Pos{current}
	for _, ok := cameFrom[current]; ok; {
		current = cameFrom[current]
		path = append([]Pos{current}, path...)

		if current == start {
			break
		}
	}

	return path
}

func h(start Pos, end Pos) float32 {
	return float32(ManhattanDistance(start, end))
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

func withDefaultInf(m map[Pos]float32, at Pos) float32 {
	if v, ok := m[at]; ok {
		return v
	}
	return 10000.0
}

func neighbors(hm HeightMap, center Pos) []Pos {
	centerHeight := hm.At(center)
	if string(rune(centerHeight)) == "S" {
		centerHeight = int('a')
	}
	if string(rune(centerHeight)) == "E" {
		centerHeight = int('z')
	}

	//     _ X _
	//     X 0 X
	//     _ X _

	rel := []Pos{
		/*    _    */ {X: -1, Y: 0}, /*    _    */
		{X: 0, Y: -1} /*    0    */, {X: 0, Y: 1},
		/*    _    */ {X: 1, Y: 0}, /*    _    */
	}

	var ret []Pos
	for _, pos := range rel {
		temp := pos
		temp.Add(center)
		if temp.X < 0 || temp.Y < 0 || temp.X >= len(hm.Height) || temp.Y >= len(hm.Height[temp.X]) {
			continue
		}

		heightAtPos := hm.At(temp)
		if centerHeight+1 >= heightAtPos {
			ret = append(ret, temp)
		}
	}

	return ret
}

func (hm HeightMap) At(pos Pos) int {
	temp := hm.Height[pos.X][pos.Y]
	if string(rune(temp)) == "S" {
		return int('a')
	}
	if string(rune(temp)) == "E" {
		return int('z')
	}
	return temp
}

func anyPos(posMap map[Pos]float32) Pos {
	for k := range posMap {
		return k
	}

	panic("No more positions in map!")
}

type Pos struct {
	X int
	Y int
}

func (p *Pos) IsPos(x int, y int) bool {
	return x == p.X && y == p.Y
}

func (p *Pos) Add(pos Pos) {
	p.X += pos.X
	p.Y += pos.Y
}

func (p Pos) String() string {
	return strings.Join([]string{strconv.Itoa(p.X), strconv.Itoa(p.Y)}, "/")
}

func (hm HeightMap) Print() {
	for _, line := range hm.Height {
		for _, v := range line {
			print(string(rune(v)))
		}
		println()
	}
}
