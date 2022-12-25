package main

import (
	"math"
	"os"
	"strings"
)

const InputFile = "25/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int {
	check(e)
	return i
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	sum := 0
	for _, line := range lines {
		sum += fromSnafu(line)
	}

	println("Sum = ", sum)
	println("Snafu Sum = ", toSnafu(sum))
}

func fromSnafu(snafu string) int {
	rsnafu := reverse(snafu)

	res := 0
	for power, s := range rsnafu {
		d := int(math.Pow(5, float64(power)))
		res += d * decodeSnafuChar(s)
	}

	return res
}

func decodeSnafuChar(c int32) int {
	switch rune(c) {
	case '2':
		return 2
	case '1':
		return 1
	case '0':
		return 0
	case '-':
		return -1
	case '=':
		return -2
	default:
		panic("Unrecognized snafu " + string(rune(c)))
	}
}

func toSnafu(d int) string {
	res := ""

	if d == 0 {
		return "0"
	}
	for d > 0 {
		res = string("=-012"[(d+2)%5]) + res
		d = (d + 2) / 5
	}

	return res
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
