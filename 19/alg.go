package main

type State01 struct {
	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int

	ore      int
	clay     int
	obsidian int
	geode    int

	newOreRobot      bool
	newClayRobot     bool
	newObsidianRobot bool
	newGeodeRobot    bool
}

func Evaluate01(b Blueprint) int {
	start := State01{
		oreRobots: 1,
	}

	firstGeodeStep := -1 // first step a geode is encountered
	bestEndState := EvaluateStep01(b, start, SimulationLength, &firstGeodeStep)
	return bestEndState.geode
}

func EvaluateStep01(b Blueprint, state State01, step int, firstGeodeStep *int) State01 {
	if step <= 0 {
		return state
	}

	// 5 Options for each step: build O, C, OB, G robot  and NOTHING
	var rstateO, rstateC, rstateOB, rstateG, rstateNothing State01

	if state.ore >= b.GeodeRobot[0] && state.obsidian >= b.GeodeRobot[1] {
		stateG := state
		stateG.newGeodeRobot = true
		updateState01(&stateG, b)
		rstateG = EvaluateStep01(b, stateG, step-1, firstGeodeStep)
	} else {
		if state.ore >= b.ObsidianRobot[0] && state.clay >= b.ObsidianRobot[1] {
			stateOB := state
			stateOB.newObsidianRobot = true
			updateState01(&stateOB, b)
			rstateOB = EvaluateStep01(b, stateOB, step-1, firstGeodeStep)
		}
		if state.ore >= b.ClayRobot && state.clayRobots <= b.ObsidianRobot[1] {
			stateC := state
			stateC.newClayRobot = true
			updateState01(&stateC, b)
			rstateC = EvaluateStep01(b, stateC, step-1, firstGeodeStep)
		}
		if state.ore >= b.OreRobot && state.oreRobots <= b.maxOreNeededForAnyRobot {
			stateO := state
			stateO.newOreRobot = true
			updateState01(&stateO, b)
			rstateO = EvaluateStep01(b, stateO, step-1, firstGeodeStep)
		}
	}

	stateNothing := state
	updateState01(&stateNothing, b)
	rstateNothing = EvaluateStep01(b, stateNothing, step-1, firstGeodeStep)

	states := []State01{rstateO, rstateC, rstateOB, rstateG, rstateNothing}
	maxGeodes := 0
	var retstate State01
	for _, state := range states {
		if state.geode > maxGeodes {
			retstate = state
			maxGeodes = state.geode
		}
	}

	return retstate
}

func updateState01(s *State01, b Blueprint) {
	if s.newGeodeRobot {
		s.ore -= b.GeodeRobot[0]
		s.obsidian -= b.GeodeRobot[1]
	}

	if s.newObsidianRobot {
		s.ore -= b.ObsidianRobot[0]
		s.clay -= b.ObsidianRobot[1]
	}

	if s.newClayRobot {
		s.ore -= b.ClayRobot
	}

	if s.newOreRobot {
		s.ore -= b.OreRobot
	}
	panicOnNegative(s.ore, s.clay, s.obsidian, s.geode)

	//
	// Collect
	//
	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geode += s.geodeRobots

	//
	// Finish building
	//
	if s.newOreRobot {
		s.oreRobots += 1
	}
	if s.newClayRobot {
		s.clayRobots += 1
	}
	if s.newObsidianRobot {
		s.obsidianRobots += 1
	}
	if s.newGeodeRobot {
		s.geodeRobots += 1
	}

	// reset for next cycle
	s.newClayRobot = false
	s.newOreRobot = false
	s.newObsidianRobot = false
	s.newGeodeRobot = false
}

func panicOnNegative(ore int, clay int, obsidian int, geode int) {
	if ore < 0 {
		panic("Ore negative")
	}
	if clay < 0 {
		panic("Clay negative")
	}
	if obsidian < 0 {
		panic("Obsidian negative")
	}
	if geode < 0 {
		panic("Geode negative")
	}
}
