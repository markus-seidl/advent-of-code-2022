package main

type State struct {
	Moves       int
	WorldMap    Map
	GoalReached bool
}

func Solve(startingMap Map, goal Pos) State {
	start := State{
		Moves:    0,
		WorldMap: startingMap,
	}
	queue := Queue[State]{
		Lifo: true,
	}
	queue.Push(start)

	dirs := []Pos{{X: 0, Y: 1}, {X: 0, Y: -1}, {X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 0}}

	minMoves := ManhattanDistance(startingMap.CurrentPos, goal) * 10
	var minState State

	i := 0
	for queue.Size() > 0 {
		state := queue.Pop()
		if i%100_000 == 0 {
			println("Queue size: ", queue.Size())
		}
		i++

		if state.GoalReached {
			if state.Moves < minMoves {
				println("Found new min moves: ", state.Moves)
				minState = state
				minMoves = state.Moves
			}

			continue
		}

		if state.Moves > minMoves {
			continue // discard too long paths
		}
		//startupHelper := false //minMoves >= math.MaxInt64 // no min moves yet, so we discard seen positions to find at least one path

		newMap := state.WorldMap.Copy()
		newMap.SimulateBlizzards()

		// decide where to go
		for _, dir := range dirs {
			newPos := Add(state.WorldMap.CurrentPos, dir)
			if !newMap.CanBe(newPos) {
				continue
			}

			nextStateMap := newMap.Copy()
			nextStateMap.CurrentPos = newPos
			nextState := State{
				Moves:       state.Moves + 1,
				WorldMap:    nextStateMap,
				GoalReached: newPos == goal,
			}

			queue.Push(nextState)
		}
	}

	return minState
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

func (m *Map) CanBe(pos Pos) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= len(m.Tiles) || pos.Y >= len(m.Tiles[0]) {
		return false
	}

	tile := m.Tiles[pos.X][pos.Y]

	if tile.IsWall {
		return false
	}

	return len(tile.Blizzards) == 0
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
