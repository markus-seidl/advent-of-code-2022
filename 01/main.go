package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

const InputFile = "01/input.txt"

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
	var elfCalories []int
	var currentElfCalories int
	var maxCalories int
	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		// println(line)

		if line == "" {
			if currentElfCalories > maxCalories {
				maxCalories = currentElfCalories
			}
			elfCalories = append(elfCalories, currentElfCalories)
			currentElfCalories = 0
			continue
		}

		calories, err := strconv.Atoi(line)
		check(err)
		currentElfCalories += calories
	}

	var fatElf int
	for elf, calories := range elfCalories {
		if calories >= maxCalories {
			println(elf+1, calories, "Max calories")
			fatElf = elf + 1
		} else {
			//println(elf+1, calories)
		}
	}
	println("Elf with the most calories", fatElf)

	//
	// Part 2
	//
	sort.Sort(sort.Reverse(sort.IntSlice(elfCalories)))
	sumTopThree := 0
	for i := 0; i < 3; i++ {
		sumTopThree += elfCalories[i]
	}
	println("Sum Top Three", sumTopThree)
}
