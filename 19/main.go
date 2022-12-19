package main

import (
	"os"
	"strconv"
	"strings"
)

const InputFile = "19/input.txt"
const SimulationLength = 24

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

type Blueprint struct {
	No                      int
	OreRobot                int   // ore
	ClayRobot               int   // ore
	ObsidianRobot           []int // ore, clay
	GeodeRobot              []int // ore, obsidian
	maxOreNeededForAnyRobot int
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	var blueprints []Blueprint

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		// Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 12 clay. Each geode robot costs 4 ore and 19 obsidian.
		parts := strings.Split(line, ".")
		part0 := strings.Split(parts[0], ":")
		blueprint := strings.ReplaceAll(part0[0], "Blueprint ", "")

		oreRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(part0[1], "Each ore robot costs", ""), " ore", ""), " ")
		clayRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(parts[1], "Each clay robot costs", ""), " ore", ""), " ")

		parts2 := strings.Split(parts[2], "ore and")
		obsidianOreRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(parts2[0], "Each obsidian robot costs", ""), "", ""), " ")
		obsidianClayRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(parts2[1], "clay", ""), "", ""), " ")

		parts3 := strings.Split(parts[3], "ore and")
		geodeOreRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(parts3[0], "Each geode robot costs", ""), "", ""), " ")
		geodeClayRobot := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(parts3[1], "obsidian", ""), "", ""), " ")

		b := Blueprint{
			No:            check2(strconv.Atoi(blueprint)),
			OreRobot:      check2(strconv.Atoi(oreRobot)),
			ClayRobot:     check2(strconv.Atoi(clayRobot)),
			ObsidianRobot: []int{check2(strconv.Atoi(obsidianOreRobot)), check2(strconv.Atoi(obsidianClayRobot))},
			GeodeRobot:    []int{check2(strconv.Atoi(geodeOreRobot)), check2(strconv.Atoi(geodeClayRobot))},
		}
		b.maxOreNeededForAnyRobot = Max4(b.OreRobot, b.ClayRobot, b.ObsidianRobot[0], b.GeodeRobot[0])
		b.Print()
		blueprints = append(blueprints, b)
	}

	// Part 01
	waitChannel := make(chan int)

	for _, b := range blueprints {
		temp := b
		go func() {
			maxGeodes := Evaluate(temp)
			println("Blueprint", temp.No, " Geodes:", maxGeodes, "Quality:", temp.No*maxGeodes)
			waitChannel <- temp.No * maxGeodes
		}()
	}

	sum := 0
	for _, _ = range blueprints {
		res := <-waitChannel
		sum += res
	}

	println("Part 01: ", sum)
}

func Max4(a int, b int, c int, d int) int {
	ab := Max(a, b)
	cd := Max(c, d)
	return Max(ab, cd)
}

type State struct {
	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int

	ore      int
	clay     int
	obsidian int
	geode    int
}

func EvaluateTest(b Blueprint) {
	s := State{
		oreRobots: 1,
	}

	for step := 1; step <= SimulationLength; step++ {
		println("== Minute", step, "==")

		newOreRobot := false
		newClayRobot := false
		newObsidianRobot := false
		newGeodeRobot := false

		//
		// Spend
		//
		if shouldBuildGeodeRobot(step, s, b) {
			println("Spend", b.GeodeRobot[0], "ore and", b.GeodeRobot[1], "obsidian to start building a geode-collecting robot.")
			s.ore -= b.GeodeRobot[0]
			s.obsidian -= b.GeodeRobot[1]
			newGeodeRobot = true
		}

		if shouldBuildObsidianRobot(step, s, b) {
			println("Spend", b.ObsidianRobot[0], "ore and", b.ObsidianRobot[1], "clay to start building a obsidian-collecting robot.")
			s.ore -= b.ObsidianRobot[0]
			s.clay -= b.ObsidianRobot[1]
			newObsidianRobot = true
		}

		if shouldBuildClayRobot(step, s, b) {
			println("Spend", b.ClayRobot, "to start building a clay-collecting robot.")
			s.ore -= b.ClayRobot
			newClayRobot = true
		}

		if shouldBuildOreRobot(step, s, b) {
			println("Spend", b.OreRobot, "to start building a ore-collecting robot.")
			s.ore -= b.OreRobot
			newOreRobot = true
		}

		//
		// Collect
		//
		s.ore += s.oreRobots
		s.clay += s.clayRobots
		s.obsidian += s.obsidianRobots
		s.geode += s.geodeRobots
		println(s.ore, "ore", s.clay, "clay", s.obsidian, "obsidian", s.geode, "geode")

		//
		// Finish building
		//
		if newOreRobot {
			s.oreRobots += 1
			println("The new ore-collecting robot is ready; you now have", s.oreRobots, "of them.")
		}
		if newClayRobot {
			s.clayRobots += 1
			println("The new clay-collecting robot is ready; you now have", s.clayRobots, "of them.")
		}
		if newObsidianRobot {
			s.obsidianRobots += 1
			println("The new obsidian-collecting robot is ready; you now have", s.obsidianRobots, "of them.")
		}
		if newGeodeRobot {
			s.geodeRobots += 1
			println("The new geode-collecting robot is ready; you now have", s.geodeRobots, "of them.")
		}

		println()
	}
}

func shouldBuildGeodeRobot(step int, s State, b Blueprint) bool {
	if s.ore >= b.GeodeRobot[0] && s.obsidian >= b.GeodeRobot[1] {
		return true
	}
	return false
}

func shouldBuildObsidianRobot(step int, s State, b Blueprint) bool {
	if s.ore >= b.ObsidianRobot[0] && s.clay >= b.ObsidianRobot[1] {
		return true
	}
	return false
}

func shouldBuildOreRobot(step int, s State, b Blueprint) bool {
	return false
}

func shouldBuildClayRobot(step int, s State, b Blueprint) bool {
	if s.ore < b.ClayRobot {
		return false
	}

	if s.clayRobots <= 3 && s.obsidianRobots == 0 {
		return true // must build because we need clay(+robots)!
	}

	if s.ore > Max(b.GeodeRobot[0], b.ObsidianRobot[0]) {
		return true // build if we have spare ore
	}

	return false
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b *Blueprint) Print() {
	println("Blueprint", b.No, ":")
	println("\tEach ore robot costs", b.OreRobot, "ore")
	println("\tEach clay robot costs", b.ClayRobot, "ore")
	println("\tEach obsidian robot costs", b.ObsidianRobot[0], "ore and", b.ObsidianRobot[1], "clay")
	println("\tEach geode robot costs", b.GeodeRobot[0], "ore and", b.GeodeRobot[1], "obsidian")
	println("\tMax ore needed for any robot", b.maxOreNeededForAnyRobot)
}
