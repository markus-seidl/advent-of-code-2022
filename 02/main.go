package main

import (
	"os"
	"strings"
)

type Hand int64

const (
	Rock     Hand = 1
	Paper         = 2
	Scissors      = 3
)

const InputFile = "02/input.txt"

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")
	total := int64(0)
	for _, line := range lines {
		if line == "" {
			break
		}

		outcome := determineResult(line)
		println(line, outcome)

		total += outcome
	}

	println("Total", total)
}

func determineResult(line string) int64 {
	parts := strings.Split(line, " ")

	leftStr := parts[0]
	var leftHand Hand

	if leftStr == "A" {
		leftHand = Rock
	} else if leftStr == "B" {
		leftHand = Paper
	} else if leftStr == "C" {
		leftHand = Scissors
	} else {
		panic(line)
	}

	rightStr := parts[1]
	var rightHand Hand

	if false { // uncomment for first part
		if rightStr == "X" {
			rightHand = Rock
		} else if rightStr == "Y" {
			rightHand = Paper
		} else if rightStr == "Z" {
			rightHand = Scissors
		} else {
			panic(line)
		}
	} else { // second part of the answer
		if rightStr == "X" {
			// right losses
			if leftHand == Rock {
				rightHand = Scissors
			} else if leftHand == Paper {
				rightHand = Rock
			} else if leftHand == Scissors {
				rightHand = Paper
			}
		} else if rightStr == "Y" {
			rightHand = leftHand
		} else if rightStr == "Z" {
			// right wins
			if leftHand == Rock {
				rightHand = Paper
			} else if leftHand == Paper {
				rightHand = Scissors
			} else if leftHand == Scissors {
				rightHand = Rock
			}
		} else {
			panic(line)
		}
	}

	if leftHand == rightHand {
		// draw
		return 3 + int64(rightHand)
	}
	// wins
	if leftHand == Rock && rightHand == Paper {
		return 6 + int64(rightHand)
	}
	if leftHand == Paper && rightHand == Scissors {
		return 6 + int64(rightHand)
	}
	if leftHand == Scissors && rightHand == Rock {
		return 6 + int64(rightHand)
	}

	// losses
	return 0 + int64(rightHand)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
