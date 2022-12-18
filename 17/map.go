package main

func (m *Map) PrintWMR(topLeft Pos, movingRock [][]Space) {
	c := m.Copy()
	c.PutBlock(topLeft, movingRock, RockMoving)

	c.Print()
}

func (m *Map) Copy() Map {
	duplicate := make([][]Space, len(m.Map))
	for i := range m.Map {
		duplicate[i] = make([]Space, len(m.Map[i]))
		copy(duplicate[i], m.Map[i])
	}

	return Map{Map: duplicate}
}

func (m *Map) Print() {

	println("|vvvvvvv| -- top")

	for y := len(m.Map) - 1; y >= 0; y-- {
		print("|")
		for x := 0; x < len(m.Map[y]); x++ {
			switch m.Map[y][x] {
			case Air:
				print(".")
			case RockStill:
				print("#")
			case RockMoving:
				print("@")
			}
		}
		println("|")
	}

	println("|-------| -- bottom")
}

func (m *Map) LastBlockOrFloorIdx() int {
	for y := len(m.Map) - 1; y >= 0; y-- {
		if m.HasRock(y) {
			return y
		}
	}
	return 0
}

func (m *Map) PutBlock(topLeft Pos, block [][]Space, asSpace Space) {
	for y := 0; y < len(block); y++ {
		for x := 0; x < len(block[y]); x++ {
			mY := topLeft.Y - y
			mX := topLeft.X + x

			if block[y][x] != Air {
				if m.Map[mY][mX] != Air && block[y][x] != Air {
					panic("unexpected hit")
				}
				m.Map[mY][mX] = asSpace
			}
		}
	}
}

func (m *Map) HitRock(topLeft Pos, block [][]Space, blockWidth int) bool {
	return !m.IsFree(topLeft, block, blockWidth)
}

func (m *Map) IsFree(topLeft Pos, block [][]Space, blockWidth int) bool {
	if topLeft.Y+1-len(block) < 0 {
		return false // would be below bottom
	}

	if topLeft.X+blockWidth > ChamberWidth {
		return false // block would be outside of chamber
	}
	if topLeft.X < 0 {
		return false // block would be outside of chamber
	}

	for y := 0; y < len(block); y++ {
		for x := 0; x < len(block[y]); x++ {
			b := block[y][x]

			if b != Air && m.Map[topLeft.Y-y][topLeft.X+x] != Air {
				return false
			}
		}
	}

	return true
}

func (m *Map) HasRock(y int) bool {
	for x := 0; x < len(m.Map[y]); x++ {
		if m.Map[y][x] != Air {
			return true
		}
	}
	return false
}

func (m *Map) ExtendMap(lines int) {
	if lines == 0 {
		return
	}

	if lines < 0 {
		m.Map = m.Map[:len(m.Map)-(-lines)]
	}

	if lines > 0 {
		for i := 0; i < lines; i++ {
			m.Map = append(m.Map, make([]Space, ChamberWidth))
		}
	}
}
