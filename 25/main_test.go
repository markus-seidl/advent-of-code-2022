package main

import "testing"

func allNumbers() map[string]int {
	ret := make(map[string]int)
	ret["0"] = 0
	ret["1"] = 1
	ret["2"] = 2
	ret["1="] = 3
	ret["1-"] = 4
	ret["10"] = 5
	ret["11"] = 6
	ret["12"] = 7
	ret["2="] = 8
	ret["2-"] = 9
	ret["20"] = 10
	ret["1=0"] = 15
	ret["1-0"] = 20
	ret["1=11-2"] = 2022
	ret["1-0---0"] = 12345
	ret["1121-1110-1=0"] = 314159265

	ret["1=-0-2"] = 1747
	ret["12111"] = 906
	ret["2=0="] = 198
	ret["21"] = 11
	ret["2=01"] = 201
	ret["111"] = 31
	ret["20012"] = 1257
	ret["112"] = 32
	ret["1=-1="] = 353
	ret["1-12"] = 107
	ret["12"] = 7
	ret["1="] = 3
	ret["122"] = 37

	return ret
}

func TestToSnafu(t *testing.T) {
	for expected, decimal := range allNumbers() {
		result := toSnafu(decimal)
		if result != expected {
			t.Fatalf("Wrong snafu %s != %s", expected, result)
		}
	}
}

func TestFromSnafu(t *testing.T) {
	for snafu, expected := range allNumbers() {
		result := fromSnafu(snafu)
		if result != expected {
			t.Fatalf("Wrong decimal %d != %d", expected, result)
		}
	}
}
