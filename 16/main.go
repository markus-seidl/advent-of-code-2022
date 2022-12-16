package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const InputFile = "16/example.txt"

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
	m.ValveMap = make(map[string]*Valve, 0)

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		parts01 := strings.Split(line, "has flow rate=")
		valveName := strings.ReplaceAll(parts01[0], "Valve ", "")
		parts02 := strings.Split(parts01[1], ";")

		flowRate := parts02[0]
		parts03 := strings.ReplaceAll(strings.ReplaceAll(parts02[1], " tunnels lead to valves", ""), " tunnel leads to valve ", "")
		nextValves := strings.Split(parts03, ", ")

		m.AddValve(&Valve{
			Name:           strings.Trim(valveName, " "),
			FlowRate:       check2(strconv.Atoi(flowRate)),
			Open:           false,
			NextValves:     nil,
			NextValveNames: nextValves,
		})
	}

	// fix nextValves
	for _, v := range m.Valves {
		for _, nextValveName := range v.NextValveNames {
			v.NextValves = append(v.NextValves, m.ValveMap[strings.Trim(nextValveName, " ")])
		}
	}
	//

	println(m.String())
	println()
	println()

	currentValve := m.ValveMap["AA"]
	for move := 1; move <= 30; move++ {
		println("== Minute ", move, " ==")
		m.PrintStatus()
		println("Current Valve: ", currentValve.Name)

		// Plan next move
		if currentValve.FlowRate > 0 {
			println("You open", currentValve.Name, "for", currentValve.FlowRate)
			currentValve.Open = true
		}
		candidates := m.CandidateValves(currentValve)
		//candidates = Filter(candidates, 3)
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Path[1].FlowRate > candidates[j].Path[1].FlowRate
		})
		println("\tCandidates:")
		for _, c := range candidates {
			println("\t\t", c.String())
		}

		nextValve := candidates[0].Path[1]
		println("You move to", nextValve.Name)
		currentValve = nextValve

		// /Plan next move

		println()
	}

	//
}

func Filter(candidates []Candidate, maxPathLength int) []Candidate {
	ret := make([]Candidate, 0)
	for _, c := range candidates {
		if len(c.Path) <= maxPathLength {
			ret = append(ret, c)
		}
	}
	return ret
}

type Map struct {
	Valves   []*Valve
	ValveMap map[string]*Valve
}

func (m *Map) AddValve(v *Valve) {
	m.Valves = append(m.Valves, v)
	m.ValveMap[v.Name] = v
}

func (m *Map) String() string {
	ret := ""
	for _, v := range m.Valves {
		ret = fmt.Sprintf("%s\n%s", ret, v.String())
	}
	return ret
}

func (m *Map) CandidateValves(currentValve *Valve) []Candidate {
	var ret []Candidate
	valvesVisited := make(map[*Valve]bool)
	currentPath := []*Valve{currentValve}
	m._candidateValves(currentValve, &ret, valvesVisited, currentPath, 0)

	return ret
}

func (m *Map) _candidateValves(currentValve *Valve, candidates *[]Candidate, valvesVisited map[*Valve]bool, currentPath []*Valve, pressureSoFar int) {

	valvesVisited[currentValve] = true

	for _, v := range currentValve.NextValves {
		if _, ok := valvesVisited[v]; ok {
			continue // prevent loops
		}

		path := append(currentPath, v)
		if v.FlowRate > 0 && !v.Open {
			*candidates = append(*candidates, Candidate{
				Pressure: v.FlowRate + pressureSoFar,
				Path:     path,
			})
		}

		copyPath := make([]*Valve, len(path))
		copy(copyPath, path)

		copyVisited := make(map[*Valve]bool, len(valvesVisited))
		for k, v := range valvesVisited {
			copyVisited[k] = v
		}

		m._candidateValves(v, candidates, copyVisited, copyPath, pressureSoFar+v.FlowRate)
	}
}

func (m *Map) PrintStatus() {
	p, vs := m.PressureRelease()
	print("Valves ")
	for _, v := range vs {
		print(v.Name, ",")
	}
	println(" are open, releasing", p, "pressure")
}

func (m *Map) PressureRelease() (int, []*Valve) {
	var openValves []*Valve
	releasedPressure := 0
	for _, v := range m.Valves {
		if v.Open {
			openValves = append(openValves, v)
			releasedPressure += v.FlowRate
		}
	}

	return releasedPressure, openValves
}

type Candidate struct {
	Pressure int
	Path     []*Valve
}

func (c *Candidate) String() string {
	p := ""
	for _, v := range c.Path {
		p = fmt.Sprintf("%s -> %s(%d)", p, v.Name, v.FlowRate)
	}

	return fmt.Sprintf("%d %s", c.Pressure, p)
}

type Valve struct {
	Name           string
	FlowRate       int
	Open           bool
	NextValves     []*Valve
	NextValveNames []string
}

func (v *Valve) String() string {
	return fmt.Sprintf("%s (%d, %v, %s [%d])", v.Name, v.FlowRate, v.Open, strings.Join(v.NextValveNames, ", "), len(v.NextValves))
}
