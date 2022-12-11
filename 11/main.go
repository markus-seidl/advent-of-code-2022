package main

import "math"

func main() {

	monkeys := GenerateInput()

	for round := 1; round <= 20; round++ {
		println("Round ", round)

		for monkeyIdx, monkey := range monkeys {
			println("Monkey ", monkeyIdx)

			tempMonkeyItems := monkey.Items
			monkey.Items = make([]int, 0) // monkeys throw out all of their stuff each turn
			for _, item := range tempMonkeyItems {
				println("\tMonkey inspects an item with a worry level of", item)
				newWorryLevel := UpdateWorryLevel(item, *monkey)
				newWorryLevel = int(math.Floor(float64(newWorryLevel) / 3.0))
				println("\t\tMonkey gets bored with item. Worry level is divided by 3 to", newWorryLevel)

				isDivisible := newWorryLevel%monkey.TestDivision == 0
				if isDivisible {
					println("\t\tCurrent worry level *is* divisible by", monkey.TestDivision)
					println("\t\tItem with worry level", newWorryLevel, "is thrown to monkey", monkey.TestTrueMonkey)
					Throw(newWorryLevel, monkey.TestTrueMonkey, monkeys)
				} else {
					println("\t\tCurrent worry level is *not* divisible by", monkey.TestDivision)
					println("\t\tItem with worry level", newWorryLevel, "is thrown to monkey", monkey.TestFalseMonkey)
					Throw(newWorryLevel, monkey.TestFalseMonkey, monkeys)
				}

				monkey.ItemInspections++
			}
		}

		println("Round ", round)
		Statistics(monkeys)
	}

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
			print(worryLevel, ", ")
		}
		println()
	}
	println("------------------------------------------------------------------")
}
