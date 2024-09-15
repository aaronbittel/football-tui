package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"syscall"
	component "tui/internal/component"

	"golang.org/x/term"
)

func main() {
	component.Start()
	fd := int(syscall.Stdin)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)
	defer component.TearDown()

	// if err != nil {
	// 	fmt.Println("Error getting terminal size", err)
	// 	return
	// }

	box := component.NewBox("Terminal Algorithm Visualizer").
		WithRoundedCorners().
		At(2, 35)

	tabs := component.NewTabs("Sorting", "Searching", "Graph", "Stack", "Array", "HashSet").
		At(5, 20)

	list := component.NewList("BubbleSort", "QuickSort", "InsertionSort", "SelectionSort", "HeapSort", "MergeSort").
		At(9, 2)

	algorithms := [][]string{
		{"BubbleSort", "QuickSort", "InsertionSort", "SelectionSort", "HeapSort", "MergeSort"},
		{"Binary", "Linear"},
		{"Breadth-First", "Depth-First"},
		{"SampleStack1", "SampleStack2", "SampleStack3"},
		{"SampleArray1", "SampleArray2", "SampleArray3"},
		{"SampleHashSet1", "SampleHashSet2", "SampleHashSet3"},
	}

	nums := []int{4, 2, 8, 1, 27, 12, 10, 7, 5, 6, 23}
	graph := component.NewGraph(nums).At(9, 30)

	component.Print(box)
	component.Print(tabs)
	component.Print(list)
	component.Print(graph)

	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', component.CtrlC:
			running = false
		case '\t':
			tabs.Next()
			component.Clear(tabs)
			component.Print(tabs)

			component.Clear(list)
			list = component.NewList(algorithms[tabs.Selected]...).At(9, 2)
			list.Selected = 0
			component.Print(list)

			if tabs.Selected > 1 {
				component.Clear(graph)
			} else {
				component.Clear(graph)
				component.Print(graph)
			}

		case 'j':
			list.Next()
			component.Clear(list)
			component.Print(list)

			if tabs.Selected > 1 {
				break
			}

			component.Clear(graph)
			newNums := make([]int, len(nums))
			for i := range len(nums) {
				newNums[i] = rand.IntN(30) + 1
			}
			graph = component.NewGraph(newNums).At(9, 30)
			component.Print(graph)

		case 'k':
			list.Prev()
			component.Clear(list)
			component.Print(list)

			if tabs.Selected > 1 {
				break
			}

			component.Clear(graph)
			newNums := make([]int, len(nums))
			for i := range len(nums) {
				newNums[i] = rand.IntN(30) + 1
			}
			graph = component.NewGraph(newNums).At(9, 30)
			component.Print(graph)

			//
			// 		case 'p':
			// 			Print(box)
			// 		case 'u':
			// 			box.WithColoredBorder(getRandomColor())
			// 			Print(box)
			// 		case 'c':
			// 			Clear(box)
			// 		case 27:
			// 			next, _, err := reader.ReadRune() // Read the next rune
			// 			if err != nil {
			// 				panic(err)
			// 			}
			//
			// 			if next == '[' { // '[' character after escape
			// 				arrow, _, err := reader.ReadRune() // Read the final rune
			// 				if err != nil {
			// 					panic(err)
			// 				}
			//
			// 				switch arrow {
			// 				case 'A': // arrow up
			// 					count++
			// 					counter.Update(fmt.Sprint(count))
			// 					Clear(counter)
			// 					counter.At(NewPos(counter.pos.row-1, counter.pos.col))
			// 					Print(counter)
			// 				case 'B': // arrow down
			// 					count--
			// 					counter.Update(fmt.Sprint(count))
			// 					Clear(counter)
			// 					counter.At(NewPos(counter.pos.row+1, counter.pos.col))
			// 					Print(counter)
			// 				case 'C': //arrow right
			// 					Clear(counter)
			// 					counter.At(NewPos(counter.pos.row, counter.pos.col+1))
			// 					Print(counter)
			// 				case 'D': //arrow left
			// 					Clear(counter)
			// 					counter.At(NewPos(counter.pos.row, counter.pos.col-1))
			// 					Print(counter)
			// 				}
			// 			}
		}
	}
}

func getRandomColor() string {
	const colorCode = "\033[%dm"
	r := rand.IntN(8) + 30
	return fmt.Sprintf(colorCode, r)
}
