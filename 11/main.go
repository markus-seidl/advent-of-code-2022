package main

import (
	"math/big"
)

func main() {

	monkeys := GenerateInput()

	for round := 1; round <= 10000; round++ {
		//println("Round ", round)

		for _, monkey := range monkeys {
			//println("Monkey ", monkeyIdx)

			tempMonkeyItems := monkey.Items
			monkey.Items = make([]big.Int, 0) // monkeys throw out all of their stuff each turn
			for _, item := range tempMonkeyItems {
				//println("\tMonkey inspects an item with a worry level of", item.Text(10))
				newWorryLevel := UpdateWorryLevel(item, *monkey)
				//newWorryLevel = int(math.Floor(float64(newWorryLevel) / 3.0))
				//println("\t\tMonkey gets bored with item. Worry level is divided by 3 to", newWorryLevel.Text(10))

				mod := big.Int{}
				mod.Mod(&newWorryLevel, big.NewInt(int64(monkey.TestDivision)))

				isDivisible := mod.CmpAbs(big.NewInt(0)) == 0
				if isDivisible {
					//println("\t\tCurrent worry level *is* divisible by", monkey.TestDivision)
					//println("\t\tItem with worry level", newWorryLevel.Text(10), "is thrown to monkey", monkey.TestTrueMonkey)
					Throw(newWorryLevel, monkey.TestTrueMonkey, monkeys)
				} else {
					//println("\t\tCurrent worry level is *not* divisible by", monkey.TestDivision)
					//println("\t\tItem with worry level", newWorryLevel.Text(10), "is thrown to monkey", monkey.TestFalseMonkey)
					Throw(newWorryLevel, monkey.TestFalseMonkey, monkeys)
				}

				monkey.ItemInspections++
			}
		}

		//Statistics(monkeys)
		if round%1000 == 0 || round == 20 || round == 1 {
			println("Round ", round)
			PrintInspectionTimes(monkeys)
			//Statistics(monkeys)
		}
	}

	//171832*172863 = 29703395016

	PrintInspectionTimes(monkeys)
}

func PrintInspectionTimes(monkeys []*Monkey) {
	println("------------------------------------------------------------------")
	for monkeyIdx, monkey := range monkeys {
		println("Monkey ", monkeyIdx, " inspected items ", monkey.ItemInspections, " times")
	}
}

func Statistics(monkeys []*Monkey) {
	// round statistics
	println("------------------------------------------------------------------")
	for monkeyIdx, monkey := range monkeys {
		print("Monkey ", monkeyIdx, ": ")
		for _, worryLevel := range monkey.Items {
			print(worryLevel.Text(10), ", ")
		}
		println()
	}
	println("------------------------------------------------------------------")
}
