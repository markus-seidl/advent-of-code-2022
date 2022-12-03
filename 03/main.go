package main

import (
	"fmt"
	"os"
	"strings"
)

const InputFile = "03/input.txt"

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
	totalPrio := int32(0)
	if false {
		for _, line := range lines {
			llen := len(line)
			parts := Chunks(line, llen/2)
			firstC, secondC := parts[0], parts[1]

			println(firstC, "-", secondC)

			firstSet := createSet(firstC)

			var prio int32
			for _, s := range secondC {
				if firstSet[s] {
					// duplicate
					prio = calculatePrio(s)
					totalPrio += prio
					fmt.Printf("%c:%d = %d\n", s, s, prio)
					break
				}
			}
		}

		println("Part 1 Prio: ", totalPrio)
	}

	for i := 0; i < len(lines); i += 3 {
		firstSet := createSet(lines[i])
		secondSet := createSet(lines[i+1])
		thirdSet := createSet(lines[i+2])

		allSet := mergeSet(firstSet, secondSet)
		allSet = mergeSet(allSet, thirdSet)

		for s, v := range firstSet {
			if secondSet[s] && thirdSet[s] && v {
				prio := calculatePrio(s)
				totalPrio += prio
				fmt.Printf("%c:%d = %d\n", s, s, prio)
			}
		}
		println("---")
	}

	println("Part 1+2 Prio: ", totalPrio)
}

func calculatePrio(s int32) int32 {
	var prio int32
	if s >= 97 {
		prio = s - 96
	} else {
		prio = 26 + s - 64
	}
	return prio
}

func mergeSet(a map[int32]bool, b map[int32]bool) map[int32]bool {
	ret := map[int32]bool{}
	for k, v := range a {
		ret[k] = v
	}
	for k, v := range b {
		temp := ret[k] && v
		if temp {
			ret[k] = temp
		}
	}

	return ret
}

func createSet(s string) map[int32]bool {
	firstSet := make(map[int32]bool)
	for _, s := range s {
		firstSet[s] = true
	}
	return firstSet
}

// a = 97
// A = 65

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
