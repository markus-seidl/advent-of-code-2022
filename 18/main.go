package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "18/example.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

type Pos3 struct {
	X int
	Y int
	Z int
}

type Map struct {
	RockMap [][][]bool
	Rocks   []Pos3
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	m := Map{}

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		posStr := strings.Split(line, ",")
		pos := Pos3{
			X: check2(strconv.Atoi(posStr[0])),
			Y: check2(strconv.Atoi(posStr[1])),
			Z: check2(strconv.Atoi(posStr[2])),
		}

		m.Rocks = append(m.Rocks, pos)
	}
	//m.PrintRocks()

	// put rocks into map
	for _, rock := range m.Rocks {
		m.ExtendSpace(rock)

		m.RockMap[rock.X][rock.Y][rock.Z] = true
	}

	if false {
		totalNotCovered := 0
		for _, rock := range m.Rocks {
			fmt.Printf("%s has ", rock.String())

			// 6 sides
			// x y z

			covered := 0
			// -1 0 0
			covered += m.IsCoveredPart01(rock, Pos3{X: -1, Y: 0, Z: 0})

			// 0 -1 0
			covered += m.IsCoveredPart01(rock, Pos3{X: 0, Y: -1, Z: 0})

			// 0 0 -1
			covered += m.IsCoveredPart01(rock, Pos3{X: 0, Y: 0, Z: -1})

			// 1 0 0
			covered += m.IsCoveredPart01(rock, Pos3{X: 1, Y: 0, Z: 0})

			// 0 1 0
			covered += m.IsCoveredPart01(rock, Pos3{X: 0, Y: 1, Z: 0})

			// 0 0 1
			covered += m.IsCoveredPart01(rock, Pos3{X: 0, Y: 0, Z: 1})

			println(covered, " not covered sides")
			totalNotCovered += 6 - covered
		}
		println("Total not covered", totalNotCovered)
	}

	// part 02
	if false {
		//sum := Pos3{}
		//for _, rock := range m.Rocks {
		//  sum = Add(rock, sum)
		//}
		//center := Pos3{X: sum.X / len(m.Rocks), Y: sum.Y / len(m.Rocks), Z: sum.Z / len(m.Rocks)}
		//println("Sum: ", sum.String(), " Mean ", center.String())

		totalNotCovered := 0
		for _, rock := range m.Rocks {
			fmt.Printf("%s has ", rock.String())

			// 6 sides
			// x y z

			notCovered := 0
			// -1 0 0
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: -1, Y: 0, Z: 0})

			// 0 -1 0
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: 0, Y: -1, Z: 0})

			// 0 0 -1
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: 0, Y: 0, Z: -1})

			// 1 0 0
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: 1, Y: 0, Z: 0})

			// 0 1 0
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: 0, Y: 1, Z: 0})

			// 0 0 1
			notCovered += m.IsNotCoveredPart02(rock, Pos3{X: 0, Y: 0, Z: 1})

			println(notCovered, " not covered sides")
			totalNotCovered += notCovered
		}
		println("Total not covered", totalNotCovered)
	}

	if true {

		m.FillSpace(Pos3{X: 10, Y: 10, Z: 10})

		for x := 0; x < len(m.RockMap); x++ {
			for y := 0; y < len(m.RockMap[x]); y++ {
				for z := 0; z < len(m.RockMap[x][y]); z++ {
					println(x, ",", y, ",", z)
				}
			}
		}

	}
}

func (m *Map) FillSpace(max Pos3) {
	for i := len(m.RockMap); i <= max.X; i++ {

	}
}

func (m *Map) IsNotCoveredPart02(rock Pos3, lookPos Pos3) int {
	// a part is not covered, if
	// - the end of the array is reached
	// - the direction doesn't have any cubes

	// where is outside?

	for t := 1; t < 100; t++ {
		test := Add(rock, Multiply(lookPos, t))

		if m.HasRock(test) {
			return 0
		}
	}

	return 1
}

func (m *Map) HasRock(p Pos3) bool {
	if len(m.RockMap) <= p.X || p.X < 0 {
		return false
	}
	if len(m.RockMap[p.X]) <= p.Y || p.Y < 0 {
		return false
	}
	if len(m.RockMap[p.X][p.Y]) <= p.Z || p.Z < 0 {
		return false
	}

	return m.RockMap[p.X][p.Y][p.Z]
}

func Multiply(p Pos3, t int) Pos3 {
	return Pos3{
		X: p.X * t,
		Y: p.Y * t,
		Z: p.Z * t,
	}
}

func Add(a Pos3, b Pos3) Pos3 {
	return Pos3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// IsCoveredPart01 note the switched returns contrary to the function name
func (m *Map) IsCoveredPart01(rock Pos3, lookPos Pos3) int {
	x := rock.X + lookPos.X
	y := rock.Y + lookPos.Y
	z := rock.Z + lookPos.Z

	covered := 1
	notCovered := 0

	if len(m.RockMap) <= x || x < 0 {
		return notCovered
	}
	if len(m.RockMap[x]) <= y || y < 0 {
		return notCovered
	}
	if len(m.RockMap[x][y]) <= z || z < 0 {
		return notCovered
	}

	if m.RockMap[x][y][z] {
		return covered
	}
	return notCovered
}

func (m *Map) ExtendSpace(pos Pos3) {
	for x := len(m.RockMap); x <= pos.X; x++ {
		m.RockMap = append(m.RockMap, make([][]bool, 0))
	}
	for y := len(m.RockMap[pos.X]); y <= pos.Y; y++ {
		m.RockMap[pos.X] = append(m.RockMap[pos.X], make([]bool, 0))
	}
	if len(m.RockMap[pos.X][pos.Y]) <= pos.Z {
		m.RockMap[pos.X][pos.Y] = append(m.RockMap[pos.X][pos.Y], make([]bool, pos.Z+1)...)
	}
}

func (m *Map) PrintRocks() {
	for _, pos := range m.Rocks {
		println(pos.String())
	}
}

func (p *Pos3) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}
