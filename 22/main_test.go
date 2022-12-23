package main

import "testing"

func TestMap_Rotate(t *testing.T) {
	m := Map{
		Direction: 0,
	}

	for i := 0; i < 10; i++ {
		m.Rotate("L")
		println("Rotate L: ", m.Direction)
	}

	for i := 0; i < 10; i++ {
		m.Rotate("R")
		println("Rotate R: ", m.Direction)
	}
}
