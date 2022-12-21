package main

type State struct {
	Minute           int
	StepPosition     *Valve
	OpenValves       map[string]bool
	ValveSteps       []*Valve
	ReleasedPressure int
}

func Evaluate02(m Map, maxMinutes int) {
	startValve := m.ValveMap["AA"]
	startState := State{
		Minute:       0,
		StepPosition: startValve,
		OpenValves:   make(map[string]bool),
		ValveSteps:   make([]*Valve, maxMinutes+1),
	}

	queue := Queue[State]{
		Lifo: true,
	}
	queue.Push(startState)

	bestPressureSoFar := -1
	bestValveCombination := make([]*Valve, maxMinutes+1)

	for queue.Size() > 0 {
		state := queue.Pop()
		if state.Minute > maxMinutes {
			pressure := CalculatePressure(state.ValveSteps)
			if pressure > bestPressureSoFar {
				println("Best so far:", pressure)
				bestPressureSoFar = pressure
				copy(bestValveCombination, state.ValveSteps)
			}
			continue
		}

		OpenValve(&state)

		// move to every connected valve, first the not opened ones with flowRate > 0
		iterated := make(map[string]bool)
		for _, nextValve := range state.StepPosition.NextValves {
			if nextValve.FlowRate > 0 {
				if _, ok := state.OpenValves[nextValve.Name]; !ok {
					iterated[nextValve.Name] = true
					queue.Push(PrepareNextStateFrom(&state, nextValve))
				}
			}
		}

		// all other next valves
		for _, nextValve := range state.StepPosition.NextValves {
			if _, ok := iterated[nextValve.Name]; ok {
				continue
			}

			queue.Push(PrepareNextStateFrom(&state, nextValve))
		}
	}

	println("Found best: ", bestPressureSoFar)
	println("Via: ")
	for _, v := range bestValveCombination {
		if v == nil {
			print(" -> - ")
		} else {
			print(" -> ", v.Name)
		}
	}
	println()
}

func OpenValve(state *State) {
	if _, ok := state.OpenValves[state.StepPosition.Name]; ok {
		return // valve already open
	}
	state.OpenValves[state.StepPosition.Name] = true
	state.ValveSteps[state.Minute] = state.StepPosition

}

func PrepareNextStateFrom(state *State, valve *Valve) State {
	ret := State{
		Minute:       state.Minute + 1,
		StepPosition: valve,
		OpenValves:   make(map[string]bool, len(state.ValveSteps)),
		ValveSteps:   make([]*Valve, len(state.ValveSteps)),
	}

	for k, v := range state.OpenValves {
		ret.OpenValves[k] = v
	}

	copy(ret.ValveSteps, state.ValveSteps)

	return ret
}

func CalculatePressure(valveStep []*Valve) int {
	sum := 0
	change := 0
	for _, v := range valveStep {
		if v == nil {
			continue
		}

		change += v.FlowRate
		sum += change
	}

	return sum
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
