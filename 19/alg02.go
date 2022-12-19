package main

type State02 struct {
	time int

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

	firstTimeGeodeBroken int
}

func Evaluate02(b Blueprint) int {
	start := State02{
		time:      SimulationLength,
		oreRobots: 1,
	}
	queue := Queue[State02]{
		Lifo: true,
	}
	queue.Push(start)

	seenStates := make(map[State02]bool, 100_000_000)

	maxGeodes := -1
	var maxState State02

	for queue.Size() > 0 {
		state := queue.Pop()
		if maxGeodes < state.geode {
			maxState = state
			maxGeodes = state.geode
			if state.time <= 0 {
				println("Max Geodes found at end: ", maxGeodes)
			}
		}

		if state.time <= 0 {
			continue // end reached
		}
		if state.time < 6 && state.geode == 0 {
			continue // not worth searching further
		}

		// can we beat the best (maxGeodes) still?
		if state.time < 20 && state.geodeRobots > 0 {
			bestGuessAtMax := state.time*(state.geodeRobots+10) + state.geode
			if bestGuessAtMax < maxGeodes {
				continue // not possible
			}
		}

		if _, ok := seenStates[state]; ok {
			continue // we have seen the same time and ore constellation already
		}
		if len(seenStates) > 100_000_000 {
			seenStates = make(map[State02]bool, 100_000_000)
		}
		seenStates[state] = true

		if state.ore >= b.GeodeRobot[0] && state.obsidian >= b.GeodeRobot[1] {
			stateG := state
			stateG.newGeodeRobot = true
			updateState02(&stateG, b)
			queue.Push(stateG)
		} else {
			if state.ore >= b.ObsidianRobot[0] && state.clay >= b.ObsidianRobot[1] {
				stateOB := state
				stateOB.newObsidianRobot = true
				updateState02(&stateOB, b)
				queue.Push(stateOB)
			}
			if state.ore >= b.ClayRobot && state.clayRobots <= b.ObsidianRobot[1] {
				stateC := state
				stateC.newClayRobot = true
				updateState02(&stateC, b)
				queue.Push(stateC)
			}
			if state.ore >= b.OreRobot && state.oreRobots <= b.maxOreNeededForAnyRobot {
				stateO := state
				stateO.newOreRobot = true
				updateState02(&stateO, b)
				queue.Push(stateO)
			}
		}

		stateNothing := state
		updateState02(&stateNothing, b)
		queue.Push(stateNothing)
	}

	println(maxState.time, maxState.geode, maxState.firstTimeGeodeBroken)

	return -1
}

func updateState02(s *State02, b Blueprint) {
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

	if s.geode > 0 && s.firstTimeGeodeBroken == 0 {
		s.firstTimeGeodeBroken = s.time
	}

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

	s.time -= 1
}

type Queue[T any] struct {
	elements []T
	Lifo     bool
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

func (q *Queue[T]) Push(s T) {
	q.elements = append(q.elements, s)
}

func (q *Queue[T]) Pop() T {
	if q.Lifo {
		temp := q.elements[len(q.elements)-1]
		q.elements = q.elements[:len(q.elements)-1]
		return temp
	} else {
		temp := q.elements[0]
		q.elements = q.elements[1:]

		return temp
	}
}
