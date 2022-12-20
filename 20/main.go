package main

import (
	"os"
	"strconv"
	"strings"
)

const InputFile = "20/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

type SortingPos struct {
	OriginalPosition int
	Number           int
	TargetPosition   int
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	var input []int

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		input = append(input, check2(strconv.Atoi(line)))
	}

	var inputS []*SortingPos
	for idx, v := range input {
		inputS = append(inputS, &SortingPos{
			OriginalPosition: idx,
			Number:           v,
		})
	}
	outputS := make([]*SortingPos, len(inputS))
	copy(outputS, inputS)
	//Print(inputS)

	for _, orig := range inputS {

		// find index in output
		curIdx := 0
		for i, v := range outputS {
			if v == orig {
				curIdx = i
				break
			}
		}

		//println(orig.Number, " moved: ")

		Move(&outputS, curIdx, outputS[curIdx].Number)
		//Print(outputS)
	}

	println("Result: ")
	//Print(outputS)

	println()
	whereIsZero := 0
	for i, v := range outputS {
		if v.Number == 0 {
			whereIsZero = i
			break
		}
	}

	a, b, c := outputS[(whereIsZero+1000)%len(outputS)].Number, outputS[(whereIsZero+2000)%len(outputS)].Number, outputS[(whereIsZero+3000)%len(outputS)].Number
	println("1000th = ", a)
	println("2000th = ", b)
	println("3000th = ", c)

	println(a + b + c)
}

func Print(s []*SortingPos) {
	for _, v := range s {
		print(v.Number, ", ")
	}
	println()
}

func Move(arr *[]*SortingPos, pos int, delta int) {
	newPos := ModInRange(pos+delta, len(*arr))
	if pos == newPos {
		return // NoOP
	}
	if newPos > pos {
		newPos = ModInRange(newPos+1, len(*arr))
	}

	no := (*arr)[pos]

	beforeNewPos := make([]*SortingPos, newPos)
	copy(beforeNewPos, (*arr)[:newPos])

	afterNewPos := make([]*SortingPos, len(*arr)-newPos)
	copy(afterNewPos, (*arr)[newPos:])

	if pos < newPos {
		beforeNewPos = append(beforeNewPos[:pos], beforeNewPos[pos+1:]...)
	} else if pos > newPos {
		tempPos := pos - newPos
		afterNewPos = append(afterNewPos[:tempPos], afterNewPos[tempPos+1:]...)
	} else {
		panic("?")
	}

	*arr = append(beforeNewPos, no)
	*arr = append(*arr, afterNewPos...)
}

func ModInRange(a int, m int) int {
	ret := a % m
	if ret <= 0 {
		ret += m - 1
	}
	return ret
}
