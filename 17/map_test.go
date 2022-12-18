package main

import "testing"

func rockTestData() (blocks [][][]Space, blockWidths []int) {
	return [][][]Space{ExtractBlock(RockShape0), ExtractBlock(RockShape1), ExtractBlock(RockShape2), ExtractBlock(RockShape3), ExtractBlock(RockShape4)}, []int{4, 3, 3, 1, 2}
}

func TestMap_LastBlockOrFloorIdx(t *testing.T) {
	m := Map{}
	m.ExtendMap(5)

	m.Map[3][0] = RockStill

	test := m.LastBlockOrFloorIdx()
	want := 3
	if test != want {
		t.Fatalf("Incorrect Idx %d vs %d", test, want)
	}
}

func TestMap_HasRock(t *testing.T) {
	m := Map{}
	m.ExtendMap(5)

	m.Map[3][0] = RockStill

	if !m.HasRock(3) {
		t.Fatalf("does have rock")
	}

	if m.HasRock(2) {
		t.Fatalf("doesn't have rock")
	}
}

func TestMap_PutBlock(t *testing.T) {
	m := Map{}
	m.ExtendMap(5)

	blocks, _ := rockTestData()

	m.PutBlock(Pos{X: 0, Y: 0}, blocks[0], RockStill) // ####

	defer func() {
		if r := recover(); r != nil {
			// TODO fatal if not recovered
		}
	}()
	m.PutBlock(Pos{X: 0, Y: 0}, blocks[0], RockStill) // --> panic
}

func TestMap_IsFree(t *testing.T) {
	m := Map{}
	m.ExtendMap(5)

	blocks, blockWidths := rockTestData()

	m.PutBlock(Pos{X: 0, Y: 0}, blocks[0], RockStill) // ####
	m.Print()

	// ###

	p := Pos{X: 0, Y: 0}
	test := m.IsFree(p, blocks[0], blockWidths[0])

	if test {
		t.Fatalf("Shouldn't be free - 1")
	}

	// .#.
	// ###
	// .#.

	p = Pos{X: 0, Y: 1}
	test = m.IsFree(p, blocks[1], blockWidths[1])

	if test {
		t.Fatalf("Shouldn't be free - 2")
	}

	// ..#
	// ..#
	// ###

	p = Pos{X: 0, Y: 1}
	test = m.IsFree(p, blocks[2], blockWidths[2])

	if test {
		t.Fatalf("Shouldn't be free - 3")
	}

	//#
	//#
	//#
	//#

	p = Pos{X: 0, Y: 3}
	test = m.IsFree(p, blocks[3], blockWidths[3])

	if test {
		t.Fatalf("Shouldn't be free - 4")
	}

	//##
	//##

	p = Pos{X: 0, Y: 1}
	test = m.IsFree(p, blocks[4], blockWidths[4])

	if test {
		t.Fatalf("Shouldn't be free - 5")
	}
}
