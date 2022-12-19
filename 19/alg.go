package main

type StateAlg struct {
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

func Evaluate(b Blueprint) int {
	start := StateAlg{
		oreRobots: 1,
	}

	firstGeodeStep := -1 // first step a geode is encountered
	bestEndState := EvaluateStep(b, start, SimulationLength, &firstGeodeStep)
	return bestEndState.geode
}

func EvaluateStep(b Blueprint, state StateAlg, step int, firstGeodeStep *int) StateAlg {
	if step <= 0 {
		return state
	}

	// 5 Options for each step: build O, C, OB, G robot  and NOTHING
	var rstateO, rstateC, rstateOB, rstateG, rstateNothing StateAlg

	if state.ore >= b.GeodeRobot[0] && state.obsidian >= b.GeodeRobot[1] {
		stateG := state
		stateG.newGeodeRobot = true
		updateState(&stateG, b)
		rstateG = EvaluateStep(b, stateG, step-1, firstGeodeStep)
	} else {
		if state.ore >= b.ObsidianRobot[0] && state.clay >= b.ObsidianRobot[1] {
			stateOB := state
			stateOB.newObsidianRobot = true
			updateState(&stateOB, b)
			rstateOB = EvaluateStep(b, stateOB, step-1, firstGeodeStep)
		}
		if state.ore >= b.ClayRobot && state.clayRobots <= b.ObsidianRobot[1] {
			stateC := state
			stateC.newClayRobot = true
			updateState(&stateC, b)
			rstateC = EvaluateStep(b, stateC, step-1, firstGeodeStep)
		}
		if state.ore >= b.OreRobot && state.oreRobots <= b.maxOreNeededForAnyRobot {
			stateO := state
			stateO.newOreRobot = true
			updateState(&stateO, b)
			rstateO = EvaluateStep(b, stateO, step-1, firstGeodeStep)
		}
	}

	stateNothing := state
	updateState(&stateNothing, b)
	rstateNothing = EvaluateStep(b, stateNothing, step-1, firstGeodeStep)

	states := []StateAlg{rstateO, rstateC, rstateOB, rstateG, rstateNothing}
	maxGeodes := 0
	var retstate StateAlg
	for _, state := range states {
		if state.geode > maxGeodes {
			retstate = state
			maxGeodes = state.geode
		}
	}

	return retstate
}

func updateState(s *StateAlg, b Blueprint) {
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
