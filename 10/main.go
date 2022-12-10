package main

import (
	"os"
	"strconv"
	"strings"
)

const InputFile = "10/input.txt"

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

	//
	// Part 01
	//
	lines := strings.Split(inputStr, "\n")
	stackH := make([]int, 0)

	result := 1
	for instrIdx := 0; instrIdx < len(lines); instrIdx++ {
		line := lines[instrIdx]

		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			value := check2(strconv.Atoi(parts[1]))
			stackH = append(stackH, result)
			result += value
			stackH = append(stackH, result)
		} else {
			stackH = append(stackH, result)
		}
	}

	//signalStrengh := 0
	//for i, stack := range stackH {
	//  print(i+1, " ", stack)
	//  if ((i+1)-20)%40 == 0 {
	//    change := (i + 1) * stackH[i-1]
	//    print("   <<<<<<", stackH[i-1], " ", change)
	//    signalStrengh += change
	//  }
	//  println()
	//}
	//
	//println("Signal Strength Sum", signalStrengh)

	// Part 02
	for i, _ := range stackH {
		//cycle := i + 1
		crtPosition := i
		linePosition := crtPosition % 40
		//stackEndCycle := stack
		stackDuringCycle := stackH[max(i-1, 0)]
		spriteMiddle := stackDuringCycle
		if linePosition == spriteMiddle-1 || linePosition == spriteMiddle || linePosition == spriteMiddle+1 {
			print("#")
		} else {
			print(".")
		}

		// println(cycle, stackDuringCycle, stackEndCycle)
		if crtPosition%40 == 39 {
			println()
		}
	}
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func check2(a int, e error) int {
	check(e)
	return a
}

// 1: 1
// 2: 1
// 3: 1 -> 4
// 4: 4
// 5: 4 -> -1
