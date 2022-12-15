package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "15/input.txt"

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

	m := Map{}

	lines := strings.Split(inputStr, "\n")
	for idx, line := range lines {
		//println(line)
		// Sensor at x=2, y=18: closest beacon is at x=-2, y=15
		temp1 := strings.Split(line, ": closest beacon is at x=")
		sensorTxt := strings.Split(strings.ReplaceAll(temp1[0], "Sensor at x=", ""), ", y=")
		beaconTxt := strings.Split(strings.ReplaceAll(temp1[1], "y=", ""), ",")

		sensor := Pos{
			X: check2(strconv.Atoi(strings.Trim(sensorTxt[0], " "))),
			Y: check2(strconv.Atoi(strings.Trim(sensorTxt[1], " "))),
		}
		beacon := Pos{
			X: check2(strconv.Atoi(strings.Trim(beaconTxt[0], " "))),
			Y: check2(strconv.Atoi(strings.Trim(beaconTxt[1], " "))),
		}

		m.Data = append(m.Data, SensorBeacon{
			Sensor: sensor,
			Beacon: beacon,
			SBDist: ManhattanDistance(sensor, beacon),
		})

		if idx == 0 {
			m.MaxX = sensor.X
			m.MinX = sensor.X
			m.MinY = sensor.Y
			m.MaxY = sensor.Y
		}

		m.ExtendMinMax(sensor)
		m.ExtendMinMax(beacon)
	}

	m.Print()

	if false {
		yCoord := 2000000
		part01Result := 0
		for x := m.MinX - 100_000_000; x < m.MaxX+100_000_000; x++ {
			p := Pos{X: x, Y: yCoord}
			if m.InRange(p) && !m.IsBeacon(p) {
				//print("#")
				part01Result++
			} else {
				//print(".")
			}
		}
		//println()

		// 4107093
		// 1M 5515610
		// 5M 6124805
		// 100M

		println("Part 01:", part01Result)
	}

	// part 02
	max := 4000000
	for x := 0; x < max; x++ {
		if x%100_000 == 0 {
			println("x = ", x)
		}

		// which ranges are not available on this X - row?
		var unavailableRange []Pos // use pos X=min Y=max

		for _, sb := range m.Data {
			diffSX := absDiffInt(sb.Sensor.X, x)
			// .................B........
			// .X......S.................
			d := sb.SBDist - diffSX

			minY := sb.Sensor.Y - d
			maxY := sb.Sensor.Y + d
			r := Pos{X: minY, Y: maxY}
			unavailableRange = append(unavailableRange, r)
		}

		// jump from range to range and find the "empty" spot
		currentY := 0
		for currentY < max {

			inRange := false
			for _, r := range unavailableRange {
				if currentY >= r.X /*minX*/ && currentY <= r.Y {
					currentY = r.Y + 1
					inRange = true
					break // found range, restart search
				}
			}

			if !inRange {
				println(x, currentY)

				if !m.InRange(Pos{X: x, Y: currentY}) { // double check
					println("Double check passed")
					return // 3138881 3364986
				}
			}
		}

		//for y := 0; y < max; y++ {
		//  p := Pos{X: x, Y: y}
		//  if !m.InRange(p) {
		//    println(x, y)
		//  }
		//}
	}
}

type Map struct {
	Data []SensorBeacon
	MinX int
	MinY int
	MaxX int
	MaxY int
}

func (m *Map) Print() {
	for _, sb := range m.Data {
		println(sb.String())
	}

	println("Range: [", m.MinX, m.MinY, "] -> [", m.MaxX, m.MaxY, "]")
}

func (m *Map) IsBeacon(p Pos) bool {
	for _, sb := range m.Data {
		if sb.Beacon.X == p.X && sb.Beacon.Y == p.Y {
			return true
		}
	}

	return false
}

// InRange determine if the position is in the range of sensor-beacon and returns true if that is the case
func (m *Map) InRange(p Pos) bool {
	for _, sb := range m.Data {
		if sb.InRange(p) {
			return true
		}
	}

	return false
}

func (m *Map) ExtendMinMax(p Pos) {
	m.MinX = Min(p.X, m.MinX)
	m.MaxX = Max(p.X, m.MaxX)

	m.MinY = Min(p.Y, m.MinY)
	m.MaxY = Max(p.Y, m.MaxY)
}

type SensorBeacon struct {
	Sensor Pos
	Beacon Pos
	SBDist int
}

// InRange determine if the position is in the range of sensor-beacon and returns true if that is the case
// --> no other beacon can be here if true
func (sb *SensorBeacon) InRange(p Pos) bool {
	return ManhattanDistance(sb.Sensor, p) <= sb.SBDist
}

func (sb *SensorBeacon) String() string {
	return fmt.Sprintf("%s - %s = %d", sb.Sensor.String(), sb.Beacon.String(), sb.SBDist)
}

type Pos struct {
	X int
	Y int
}

func (p *Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
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
