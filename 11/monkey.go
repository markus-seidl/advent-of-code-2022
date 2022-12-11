package main

import (
	"strconv"
	"strings"
)

type OperationType int64

const (
	Multiply OperationType = 0
	Add                    = 1
	Pow                    = 2
)

type Monkey struct {
	Name            int
	OperationType   OperationType // 0 = multiply, 1 = add, 2 = pow
	OperationValue  int
	TestDivision    int
	TestTrueMonkey  int
	TestFalseMonkey int
	//
	Items           []int
	ItemInspections int
}

func GenerateInput() []*Monkey {
	return []*Monkey{
		{
			Name:            0,
			OperationType:   Multiply,
			OperationValue:  7,
			TestDivision:    3,
			TestTrueMonkey:  4,
			TestFalseMonkey: 1,
			Items:           []int{64, 89, 65, 95},
		},
		{
			Name:            1,
			OperationType:   Add,
			OperationValue:  5,
			TestDivision:    13,
			TestTrueMonkey:  7,
			TestFalseMonkey: 3,
			Items:           []int{76, 66, 74, 87, 70, 56, 51, 66},
		},
		{
			Name:            2,
			OperationType:   Pow,
			OperationValue:  0,
			TestDivision:    2,
			TestTrueMonkey:  6,
			TestFalseMonkey: 5,
			Items:           []int{91, 60, 63},
		},
		{
			Name:            3,
			OperationType:   Add,
			OperationValue:  6,
			TestDivision:    11,
			TestTrueMonkey:  2,
			TestFalseMonkey: 6,
			Items:           []int{92, 61, 79, 97, 79},
		},
		{
			Name:            4,
			OperationType:   Multiply,
			OperationValue:  11,
			TestDivision:    5,
			TestTrueMonkey:  1,
			TestFalseMonkey: 7,
			Items:           []int{93, 54},
		},
		{
			Name:            5,
			OperationType:   Add,
			OperationValue:  8,
			TestDivision:    17,
			TestTrueMonkey:  4,
			TestFalseMonkey: 0,
			Items:           []int{60, 79, 92, 69, 88, 82, 70},
		},
		{
			Name:            6,
			OperationType:   Add,
			OperationValue:  1,
			TestDivision:    19,
			TestTrueMonkey:  0,
			TestFalseMonkey: 5,
			Items:           []int{64, 57, 73, 89, 55, 53},
		},
		{
			Name:            7,
			OperationType:   Add,
			OperationValue:  4,
			TestDivision:    7,
			TestTrueMonkey:  3,
			TestFalseMonkey: 2,
			Items:           []int{62},
		},
	}
}

func GenerateExample() []*Monkey {
	return []*Monkey{
		{
			Name:            0,
			OperationType:   Multiply,
			OperationValue:  19,
			TestDivision:    23,
			TestTrueMonkey:  2,
			TestFalseMonkey: 3,
			Items:           []int{79, 98},
		},
		{
			Name:            1,
			OperationType:   Add,
			OperationValue:  6,
			TestDivision:    19,
			TestTrueMonkey:  2,
			TestFalseMonkey: 0,
			Items:           []int{54, 65, 75, 74},
		},
		{
			Name:            2,
			OperationType:   Pow,
			OperationValue:  13,
			TestDivision:    13,
			TestTrueMonkey:  1,
			TestFalseMonkey: 3,
			Items:           []int{79, 60, 97},
		},
		{
			Name:            3,
			OperationType:   Add,
			OperationValue:  3,
			TestDivision:    17,
			TestTrueMonkey:  0,
			TestFalseMonkey: 1,
			Items:           []int{74},
		},
	}
}

func UpdateWorryLevel(item int, monkey Monkey) int {
	switch monkey.OperationType {
	case Multiply:
		ret := item * monkey.OperationValue
		println("\t\tWorry level is multiplicated by", monkey.OperationValue, "to", ret)
		return ret
	case Add:
		ret := item + monkey.OperationValue
		println("\t\tWorry level increases by", monkey.OperationValue, "to", ret)
		return ret
	case Pow:
		ret := item * item
		println("\t\tWorry level is multiplicated by itself", "to", ret)
		return ret
	}

	panic(strings.Join([]string{"Unknown operation ", strconv.FormatInt(int64(monkey.OperationType), 10)}, " "))
}

func Throw(worryLevel int, throwToMonkey int, monkeys []*Monkey) {
	monkeys[throwToMonkey].Items = append(monkeys[throwToMonkey].Items, worryLevel)
}
