package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "13/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type PackagePair struct {
	Left  *Packet
	Right *Packet
}

type Packet struct {
	Data []*PacketElement
}

type PacketElement struct {
	IntVal int
	Packet *Packet
}

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	lines := strings.Split(inputStr, "\n")

	var packetPairs [][]string
	for i := 0; i < len(lines); i += 3 {
		packetPairs = append(packetPairs, []string{lines[i], lines[i+1]})
	}

	// parse input
	orderSum := 0
	for no, stringPair := range packetPairs {
		println("== Pair", no+1, " ==")

		result := Compare("", ParsePacket(stringPair[0]), ParsePacket(stringPair[1]))
		println("Result: ", result, "Pair Index: ", no+1)
		if result != -1 {
			orderSum += no + 1
		}
	}

	println("Sum of pairs in order: ", orderSum)
}

// 1 = true
// 0 = equal
// -1 = false
func Compare(prefix string, pl *Packet, pr *Packet) int {
	println(prefix, "- Compare", pl.String(), "vs", pr.String())
	prefix += "  "

	if len(pl.Data) == 0 && len(pr.Data) == 0 {
		println(prefix, "- Empty lists are equal")
		return 0
	}

	for pli, l := range pl.Data {
		if len(pr.Data) <= pli {
			println(prefix, "- Right side ran out of items, so inputs are *not* in the right order")
			return -1
		}

		r := pr.Data[pli]
		println(prefix, "- Compare", l.String(), "vs", r.String())

		if l.IsInt() && r.IsInt() {
			if l.IntVal < r.IntVal {
				println(prefix, "- Left side is smaller, so inputs are in the right order")
				return 1
			}
			if l.IntVal > r.IntVal {
				println(prefix, "- Right side is smaller, so inputs are *not* in the right order")
				return -1
			}
		} else if !l.IsInt() && !r.IsInt() {
			tempR := Compare(fmt.Sprintf("  %s", prefix), l.Packet, r.Packet)
			if tempR != 0 {
				return tempR
			}
		} else {
			var temp *Packet
			if l.IsInt() {
				temp = &Packet{Data: []*PacketElement{{IntVal: l.IntVal}}}
				println(prefix, "- Mixed types; convert left to ", temp.String(), " and retry comparison")
				tempR := Compare(fmt.Sprintf("  %s", prefix), temp, r.Packet)
				if tempR != 0 {
					return tempR
				}
			} else {
				temp = &Packet{Data: []*PacketElement{{IntVal: r.IntVal}}}
				println(prefix, "- Mixed types; convert right to ", temp.String(), " and retry comparison")
				tempR := Compare(fmt.Sprintf("  %s", prefix), l.Packet, temp)
				if tempR != 0 {
					return tempR
				}
			}
		}
	}

	if prefix == "  " {
		println("---")
	}

	return 0
}

func ParsePacket(packet string) *Packet {
	println("Parsing", packet)

	var stack []*Packet
	stack = append(stack, &Packet{})

	for i := 1; i < len(packet)-1; i++ {
		c := packet[i]
		se := stack[len(stack)-1]
		char := string(rune(c))
		if char == "[" {
			// push
			deeper := &Packet{}
			se.Data = append(se.Data, &PacketElement{Packet: deeper})
			stack = append(stack, deeper)
		} else if char == "]" {
			// pop
			stack = stack[:len(stack)-1]
		} else if char == "," {
			// ignore
		} else {
			// read int fully
			nextCommaA := strings.Index(packet[i:], ",")
			nextCommaB := strings.Index(packet[i:], "]")
			nextComma := customMin(nextCommaA, nextCommaB)

			ci, err := strconv.Atoi(packet[i : i+nextComma])
			check(err)
			se.Data = append(se.Data, &PacketElement{IntVal: ci})

			i += nextComma - 1
		}
	}

	println("  --> ", stack[0].String())

	if strings.Trim(packet, " ") != strings.Trim(stack[0].String(), " ") {
		panic("Strings are not equal")
	}

	return stack[0]
}

func customMin(a int, b int) int {
	if a >= 0 && b >= 0 {
		if a > b {
			return b
		}
		return a
	}

	// protect -1
	if a < 0 {
		return b
	}
	return a
}

func (p *PacketElement) IsInt() bool {
	return p.Packet == nil
}

func (p *PacketElement) String() string {
	if p.IsInt() {
		return strconv.Itoa(p.IntVal)
	}
	return p.Packet.String()
}

func (p *Packet) String() string {
	if len(p.Data) == 0 {
		return "[]"
	}

	var elements []string
	for _, e := range p.Data {
		elements = append(elements, e.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ","))
}
