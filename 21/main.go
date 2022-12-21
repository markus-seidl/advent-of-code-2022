package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const InputFile = "21/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check2(i int, e error) int64 {
	check(e)
	return int64(i)
}

type Expression struct {
	NumberL   int64
	StringL   string
	Operation string
	NumberR   int64
	StringR   string
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	monkeys := make(map[string]*Expression)

	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		monkey := parts[0]

		var temp Expression
		exprs := strings.Split(strings.Trim(parts[1], " "), " ")
		if len(exprs) <= 2 {
			temp = Expression{
				NumberL: check2(strconv.Atoi(exprs[0])),
			}
		} else {
			temp = Expression{
				StringL:   exprs[0],
				Operation: exprs[1],
				StringR:   exprs[2],
			}
		}

		monkeys[monkey] = &temp
	}

	monkeys["root"].Operation = "-"
	AlternateSortFind(monkeys)

	lastResult := int64(1)
	value := int64(1) // manually helping it to converge and not oscilate
	change := int64(2)
	multiplier := float64(2)
	i := 0
	for true {
		// COPY
		copyMonkeys := make(map[string]*Expression, len(monkeys))
		for k, v := range monkeys {
			copyExpression := *v
			copyMonkeys[k] = &copyExpression
		}

		copyMonkeys["humn"].NumberL = value
		calculateMonkeysFaster(copyMonkeys)

		rootMonkey := copyMonkeys["root"]
		result := rootMonkey.NumberL
		if result == 0 {
			println("Found: ", value, result, change)
			println("Took", i, "iterations")
			return
		}
		if (result < 0 && lastResult > 0) || (result > 0 && lastResult < 0) {
			if change < 0 {
				change = 1
			} else {
				change = -1
			}
			multiplier -= 0.1
			println(value, result, change, multiplier)
		} else { // no sign change --> speed up
			if change < 0 {
				change = int64(math.Max(math.Floor(float64(change)*multiplier), -100_000_000_000))
			} else {
				change = int64(math.Min(math.Ceil(float64(change)*multiplier), 100_000_000_000))
			}
		}
		if i%1_000 == 0 {
			println(value, result, change, multiplier)
		}

		value += change
		lastResult = result
		i++
	}

	println("Result: ", monkeys["root"].NumberL)
}

func AlternateSortFind(monkeys map[string]*Expression) {
	// need a large value that is producing the opposite sign from Evaluate(x, 0) for this to work
	pos, found := sort.Find(int((1<<63)/2), func(val int) int {
		temp := Evaluate(monkeys, val)
		println("sort.Find = ", temp, "at", val)
		return temp
	})

	println("sort.Find result = ", pos, found)
}

func Evaluate(monkeys map[string]*Expression, at int) int {
	copyMonkeys := make(map[string]*Expression, len(monkeys))
	for k, v := range monkeys {
		copyExpression := *v
		copyMonkeys[k] = &copyExpression
	}

	copyMonkeys["humn"].NumberL = int64(at)
	calculateMonkeysFaster(copyMonkeys)

	rootMonkey := copyMonkeys["root"]
	return int(rootMonkey.NumberL)
}

func calculateMonkeysFaster(monkeys map[string]*Expression) {
	rootMonkey := monkeys["root"]

	calculateMonkey(&monkeys, rootMonkey.StringL)
	calculateMonkey(&monkeys, rootMonkey.StringR)

	mL := (monkeys)[rootMonkey.StringL]
	mR := (monkeys)[rootMonkey.StringR]

	if mL.IsNumber() && mR.IsNumber() {
		result := Calculate(mL.NumberL, rootMonkey.Operation, mR.NumberL)
		rootMonkey.StoreResult(result)
	} else {
		panic(fmt.Sprintf("Calculated %s %s but no result", rootMonkey.StringL, rootMonkey.StringR))
	}
}

func calculateMonkey(monkeys *map[string]*Expression, monkey string) {
	toCalculate := (*monkeys)[monkey]
	if toCalculate.IsNumber() {
		return
	}

	calculateMonkey(monkeys, toCalculate.StringL)
	calculateMonkey(monkeys, toCalculate.StringR)

	mL := (*monkeys)[toCalculate.StringL]
	mR := (*monkeys)[toCalculate.StringR]

	if mL.IsNumber() && mR.IsNumber() {
		result := Calculate(mL.NumberL, toCalculate.Operation, mR.NumberL)
		toCalculate.StoreResult(result)
	} else {
		panic(fmt.Sprintf("Calculated %s %s but no result", toCalculate.StringL, toCalculate.StringR))
	}
}

func FirstNonNumberMonkey(monkeys map[string]*Expression) (string, *Expression, error) {
	for monkeyName, v := range monkeys {
		if monkeyName != "root" && !v.IsNumber() {
			return monkeyName, v, nil
		}
	}

	return "", nil, fmt.Errorf("Not found")
}

func calculateMonkeys(monkeys map[string]*Expression) {
	foundNonNumberMonkey := true
	for foundNonNumberMonkey {
		foundNonNumberMonkey = false
		for monkeyName, v := range monkeys {
			if !v.IsNumber() && monkeyName != "root" {
				// try to calculate
				monkeyL := monkeys[v.StringL]
				monkeyR := monkeys[v.StringR]

				if monkeyL.IsNumber() && monkeyR.IsNumber() {
					result := Calculate(monkeyL.NumberL, v.Operation, monkeyR.NumberL)
					v.StoreResult(result)
					//println("Calculated", monkeyName, " = ", result)
				}
				foundNonNumberMonkey = true
			}
		}
	}
}

func (e *Expression) IsNumber() bool {
	return e.Operation == ""
}

func (e *Expression) StoreResult(result int64) {
	e.Operation = ""
	e.StringR = ""
	e.StringL = ""

	e.NumberL = result
}

func Calculate(numberL int64, operation string, numberR int64) int64 {
	switch operation {
	case "+":
		return numberL + numberR
	case "-":
		return numberL - numberR
	case "/":
		return numberL / numberR
	case "*":
		return numberL * numberR
	}
	panic(fmt.Sprintf("Unknown operation: %s", operation))
}
