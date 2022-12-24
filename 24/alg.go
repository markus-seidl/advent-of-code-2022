package main

type State struct {
	Moves       int
	WorldMap    Map
	GoalReached bool
	NoopMoves   int
}

type BlizzardState struct {
	Tiles [][]Tile
}

type SolveStateCache struct {
	Move int
	Pos  Pos
}

func Solve(startingMap Map, goal Pos) State {
	start := State{
		Moves:     0,
		WorldMap:  startingMap,
		NoopMoves: 0,
	}
	queue := Queue[State]{
		Lifo: true,
	}
	queue.Push(start)

	dirs := []Pos{{X: 0, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: 0}}

	var blizzardStateCache = make(map[int]BlizzardState)
	var solveStateCache = make(map[SolveStateCache]bool)
	revisitMovePrevented := 0

	startMinMoves := ManhattanDistance(startingMap.CurrentPos, goal) * 3
	minMoves := startMinMoves
	println("Manhattan Distance:", ManhattanDistance(startingMap.CurrentPos, goal))
	var minState State

	pruned := 0
	i := 0
	for queue.Size() > 0 {
		state := queue.Pop()
		if i%10_000_000 == 0 {
			print("Queue size: ", queue.Size(), " Cache Size:", len(blizzardStateCache), " Revisit:", revisitMovePrevented)
			if queue.Size() > 0 {
				first := queue.elements[0]
				last := queue.elements[len(queue.elements)-1]
				print(" First: ", first.Moves, "(", first.NoopMoves, ") Last: ", last.Moves, "(", last.NoopMoves, ")")
			}
			println()
		}
		i++

		if state.GoalReached {
			if state.Moves < minMoves {
				println("Found new min moves: ", state.Moves, " (", state.NoopMoves, ")")
				minState = state
				minMoves = state.Moves
			}

			continue
		}

		if state.Moves > minMoves {
			pruned++
			continue // discard too long paths
		}

		var newMap Map
		if tileMap, ok := blizzardStateCache[state.Moves]; ok {
			newMap = Map{
				Tiles: tileMap.Tiles,
			}
		} else {
			newMap = state.WorldMap.SimulateBlizzardsIntoNewMap()
			blizzardStateCache[state.Moves] = BlizzardState{
				Tiles: newMap.Tiles,
			}
		}

		// decide where to go
		for _, dir := range dirs {
			newPos := Add(state.WorldMap.CurrentPos, dir)
			if !newMap.CanBe(newPos) {
				continue
			}

			stateCacheKey := SolveStateCache{Move: state.Moves + 1, Pos: newPos}
			if solveStateCache[stateCacheKey] {
				revisitMovePrevented++
				continue
			}
			solveStateCache[stateCacheKey] = true

			isNoopMove := 0
			if dir.X == 0 && dir.Y == 0 {
				isNoopMove = 1
			}
			nextState := State{
				Moves: state.Moves + 1,
				WorldMap: Map{
					Tiles:      newMap.Tiles,
					CurrentPos: newPos,
				},
				GoalReached: newPos == goal,
				NoopMoves:   state.NoopMoves + isNoopMove,
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
