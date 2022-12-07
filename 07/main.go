package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const InputFile = "07/input.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Node struct {
	Name        string
	IsDirectory bool
	Parent      *Node
	Size        int64
	children    []*Node
}

const PrefixCD = "$ cd "
const PrefixLS = "$ ls"
const PrefixDir = "dir"

func main() {
	// Load file
	inputBytes, err := os.ReadFile(InputFile)
	check(err)
	inputStr := string(inputBytes)

	//
	// Part 01
	//
	lines := strings.Split(inputStr, "\n")

	rootNode := Node{
		Name:        "/",
		IsDirectory: true,
		Parent:      nil,
		Size:        -1,
		children:    make([]*Node, 0),
	}
	currentDir := &rootNode

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line[:1] != "$" {
			panic("found non-command")
		}

		// decode command
		if strings.HasPrefix(line, PrefixLS) {
			// decode until next command
			nextCommand := i + 1
			for ; nextCommand < len(lines); nextCommand++ {
				if lines[nextCommand][:1] == "$" {
					break
				}
			}
			startListing := i
			i = nextCommand - 1

			// decode directory listing from startListing until nextCommand - 1/i
			for ii := startListing + 1; ii <= i; ii++ {
				fileOrDir := createFileOrDir(currentDir, lines[ii])
				currentDir.children = append(currentDir.children, &fileOrDir)
			}
		} else if strings.HasPrefix(line, PrefixCD) {
			directoryStr := line[len(PrefixCD):]
			currentDir = changeDirectory(&rootNode, currentDir, directoryStr)
		} else {
			panic("Unknown command: " + line)
		}

		//printTree(rootNode)
		//println("----------------------")
	}

	calculateSize(&rootNode)
	printTree(rootNode)

	println("----------------")
	solvePart01(&rootNode)
	println("----------------")
	solvePart02(&rootNode)
}

func solvePart01(rootNode *Node) {
	var dirs []*Node

	stepIntoIfSmaller(rootNode, 100000, &dirs)

	totalSize := 0
	for _, dir := range dirs {
		print(dir.Name, " ; ")
		totalSize += int(dir.Size)
	}
	println()
	println("TotalSize: ", totalSize)
}

func stepIntoIfSmaller(node *Node, limit int64, dirs *[]*Node) {
	if !node.IsDirectory {
		return
	}

	if node.Size < limit {
		*dirs = append(*dirs, node)
		println("Found directory " + node.Name + " with size " + strconv.Itoa(int(node.Size)))
	}

	for _, child := range node.children {
		if child.IsDirectory {
			stepIntoIfSmaller(child, limit, dirs)
		}
	}
}

func solvePart02(rootNode *Node) {
	var currentDir *Node
	currentDir = rootNode

	minFreeSize := 30000000 - (70000000 - rootNode.Size)
	smallestDir := findSmallestDir(rootNode, minFreeSize, currentDir)

	println(smallestDir.Name, smallestDir.Size)
}

func findSmallestDir(node *Node, limit int64, currentSmallestDir *Node) *Node {
	if !node.IsDirectory {
		return currentSmallestDir
	}

	currentSmallestSize := currentSmallestDir.Size
	if node.Size < currentSmallestSize && node.Size >= limit {
		currentSmallestDir = node
	}

	for _, child := range node.children {
		temp := findSmallestDir(child, limit, currentSmallestDir)
		if temp.Size < currentSmallestDir.Size {
			currentSmallestDir = temp
		}
	}

	return currentSmallestDir
}

func printTree(rootNode Node) {
	printChildren("", &rootNode)
}

func printChildren(prefix string, node *Node) {
	if node.IsDirectory {
		fmt.Printf("%s - %s (dir, size=%d)\n", prefix, node.Name, node.Size)
		for _, child := range node.children {
			printChildren(prefix+"  ", child)
		}
	} else {
		fmt.Printf("%s - %s (file, size=%d)\n", prefix, node.Name, node.Size)
	}
}

func calculateSize(node *Node) int64 {
	totalSize := int64(0)
	for _, child := range node.children {
		if child.IsDirectory {
			totalSize += calculateSize(child)
		} else {
			totalSize += child.Size
		}
	}
	node.Size = totalSize
	return totalSize
}

func changeDirectory(rootNode *Node, currentDir *Node, directoryStr string) *Node {
	directoryStr = strings.Trim(directoryStr, " ")
	if directoryStr == "/" {
		return rootNode
	} else if directoryStr == ".." {
		if currentDir.Parent == nil {
			panic("Can't change to upper directory because curdir doesn't have a parent")
		}
		return currentDir.Parent
	}
	// must be a directory in this directory
	for _, child := range currentDir.children {
		if child.Name == directoryStr {
			return child
		}
	}

	panic("Unknown command " + directoryStr)
}

func createFileOrDir(currentDir *Node, line string) Node {
	if strings.HasPrefix(line, PrefixDir) { // directory
		return Node{
			Name:        strings.Trim(line[len(PrefixDir):], " "),
			IsDirectory: true,
			Parent:      currentDir,
			Size:        -1,
			children:    make([]*Node, 0),
		}
	} else { // file
		parts := strings.Split(line, " ")
		return Node{
			Name:        strings.Trim(parts[1], " "),
			IsDirectory: false,
			Parent:      currentDir,
			Size:        check2(strconv.Atoi(parts[0])),
			children:    nil,
		}
	}
}

func check2(a int, e error) int64 {
	check(e)
	return int64(a)
}
