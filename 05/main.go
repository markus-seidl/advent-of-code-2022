package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "05/input.txt"

type Field struct {
	Columns [][]string
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

	var fieldLines []string
	var instrLines []string

	field_sw := false
	for _, line := range lines {
		if line == "" {
			field_sw = true
			continue
		}

		if field_sw {
			instrLines = append(instrLines, line)
		} else {
			fieldLines = append(fieldLines, line)
		}
	}

	field := parseField(fieldLines)

	for _, line := range instrLines {
		// decode instructions

		// move 1 from 2 to 1
		parts01 := strings.Split(line, " from ")
		parts010 := strings.Split(parts01[0], " ")
		moveAmount := check2(strconv.Atoi(parts010[1]))

		parts011 := strings.Split(parts01[1], " to ")
		from := check2(strconv.Atoi(parts011[0])) - 1
		to := check2(strconv.Atoi(parts011[1])) - 1

		fmt.Printf("#%d %d -> %d\n", moveAmount, from, to)

		//field.MoveRepeatedlyV1(moveAmount, from, to)
		field.MoveV2(moveAmount, from, to)
		field.Print()
	}

	for _, column := range field.Columns {
		print(column[len(column)-1])
	}
	println()
}

func check2(a int, e error) int {
	check(e)
	return a
}

func parseField(fieldLines []string) Field {
	var columns [][]string

	for lineIdx := len(fieldLines) - 2; lineIdx >= 0; lineIdx-- {
		line := fieldLines[lineIdx]

		if lineIdx == len(fieldLines)-1 {
			break
		}
		rows := Chunks(line, 4)

		for colIdx, col := range rows {
			col = strings.Trim(col, " []")
			if col == "" {
				continue
			}

			if len(columns) <= colIdx {
				toAdd := []string{col}
				columns = append(columns, toAdd)
			} else {
				columns[colIdx] = append(columns[colIdx], col)
			}
		}
	}

	return Field{
		Columns: columns,
	}
}

func (f *Field) Print() {
	println("------")
	for _, column := range f.Columns {
		for _, element := range column {
			print(element, " ")
		}
		println()
	}
	println("------")
}

func (f *Field) MoveRepeatedlyV1(amount int, from int, to int) {
	for i := 0; i < amount; i++ {
		f.Move(from, to)
	}
}

func (f *Field) MoveV2(amount int, from int, to int) {
	col := f.Columns[from]
	toMove := col[len(col)-amount:]
	col = col[:len(col)-amount]

	f.Columns[from] = col

	toCol := f.Columns[to]
	toCol = append(toCol, toMove...)
	f.Columns[to] = toCol
}

func (f *Field) Move(from int, to int) {
	toMove := f.pop(from)

	toCol := f.Columns[to]
	toCol = append(toCol, toMove)
	f.Columns[to] = toCol
}

func (f *Field) pop(colIdx int) string {
	col := f.Columns[colIdx]
	last := col[len(col)-1]
	col = col[:len(col)-1]

	f.Columns[colIdx] = col

	return last
}

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
