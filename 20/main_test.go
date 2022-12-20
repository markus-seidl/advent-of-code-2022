package main

import "testing"

func TestMove(t *testing.T) {

	println("Move 5 -2")

	s := CreateDummy()
	Print(s)
	Move(&s, 5, -2)
	Print(s)
	if s[3].Number != 5 {
		t.Fatal("5 not on position.")
	}

	println("Move 5 10")

	s = CreateDummy()
	Print(s)
	Move(&s, 5, 10)
	Print(s)
	if s[5].Number != 5 {
		t.Fatal("5 not on position.")
	}

	println("Move 5 12")

	s = CreateDummy()
	Print(s)
	Move(&s, 5, 12)
	Print(s)
	if s[7].Number != 5 {
		t.Fatal("5 not on position.")
	}

	println("Move 1 -3")

	s = CreateDummy()
	Print(s)
	Move(&s, 1, -3)
	Print(s)
	if s[7].Number != 1 {
		t.Fatal("1 not on position.")
	}

}

func CreateDummy() []*SortingPos {
	ret := make([]*SortingPos, 0)
	for i := 0; i < 10; i++ {
		ret = append(ret, &SortingPos{
			OriginalPosition: i,
			Number:           i,
		})
	}

	return ret
}
