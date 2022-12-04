package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "04/input.txt"

type Range struct {
	Start int
	End   int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")
	contains := 0
	overlaps := 0

	for _, line := range lines {
		elfes := strings.Split(line, ",")

		firstElf := parseElf(elfes[0])
		secondElf := parseElf(elfes[1])

		contain := FullyContains(firstElf, secondElf)
		overlap := Overlap(firstElf, secondElf)

		fmt.Printf("%v %v %v\n", elfes, contain, overlap)
		if contain {
			contains++
		}
		if overlap {
			overlaps++
		}
	}

	println("Fully contains", contains, "Overlaps", overlaps)
}

func parseElf(elf string) Range {
	parts := strings.Split(elf, "-")
	return Range{
		Start: check2(strconv.Atoi(parts[0])),
		End:   check2(strconv.Atoi(parts[1])),
	}
}

func check2(a int, e error) int {
	check(e)
	return a
}

func Overlap(a Range, b Range) bool {
	// ---aaaa---
	// --bbb-----
	e := a.Start >= b.Start && b.End >= a.Start && b.End <= a.End

	// ---aaaa---
	// ----bbbb--
	f := a.Start <= b.Start && a.End <= b.End && b.Start <= a.End

	// ---aaaa---
	// ----b-----
	g := a.Start <= b.Start && a.End >= b.End

	// ---aaaa---
	// --bbbbbb--
	h := a.Start >= b.Start && a.End <= b.End

	return e || f || g || h
}

// FullyContains checks if a contains b fully or the other way around
func FullyContains(a Range, b Range) bool {
	return FullyContains2nd(a, b) || FullyContains2nd(b, a)
}

// FullyContains2nd checks if the second range is fully contained in the main range
func FullyContains2nd(main Range, second Range) bool {
	return main.Start <= second.Start && main.End >= second.End
}
