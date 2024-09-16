package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"
	algo "tui/internal/algorithms"
	"tui/internal/component"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

func toColString(nums []int) string {
	var b strings.Builder

	m := slices.Max(nums)
	spaces := len(fmt.Sprint(m))
	var char string

	for i := m; i >= 1; i-- {
		for _, n := range nums {
			if n >= i {
				char = utils.FullBlock
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range nums {
		s := spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}
	return b.String()
}

func main() {
	fd, oldState, err := utils.Start()
	if err != nil {
		fmt.Println("error initializing raw terminal", err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)
	defer utils.TearDown()

	box := component.NewBox("Terminal Algorithm Visualizer").
		WithRoundedCorners().
		At(2, 35)

	tabs := component.NewTabs("Sorting", "Searching", "Graph", "Stack", "Array", "HashSet").
		At(5, 20)

	list := component.NewList(
		"BubbleSort",
		"QuickSort",
		"InsertionSort",
		"SelectionSort",
		"HeapSort",
		"MergeSort").
		At(9, 2)

	algorithms := [][]string{
		{"BubbleSort", "QuickSort", "InsertionSort", "SelectionSort", "HeapSort", "MergeSort"},
		{"Binary", "Linear"},
		{"Breadth-First", "Depth-First"},
		{"SampleStack1", "SampleStack2", "SampleStack3"},
		{"SampleArray1", "SampleArray2", "SampleArray3"},
		{"SampleHashSet1", "SampleHashSet2", "SampleHashSet3"},
	}

	nums := []int{4, 12, 1, 6, 13, 8, 11, 7, 2, 9, 3, 5}
	columnGraph := component.NewColumnGraph(component.NewColumn(nums, nil)).At(10, 25)

	component.Print(columnGraph)
	component.Print(box)
	component.Print(tabs)
	component.Print(list)

	reader := bufio.NewReader(os.Stdin)
	running := true

	// FOR Partial Graph Update
	columnCh := make(chan component.Column)
	go algo.BubbleSort(columnCh, columnGraph.Nums())

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', utils.CtrlC:
			running = false
		case 't':
			col := <-columnCh
			columnGraph.Update(col)
		case '\t':
			tabs.Next()

			component.Clear(list)
			list = component.NewList(algorithms[tabs.Selected]...).At(9, 2)
			list.Selected = 0
			component.Print(list)

		case 'j':
			list.Next()

			if tabs.Selected > 1 {
				break
			}

		case 'k':
			list.Prev()

			if tabs.Selected > 1 {
				break
			}
		case 13:
			if tabs.Selected == 0 && list.Selected == 0 {

				columnCh := make(chan component.Column)
				go algo.BubbleSort(columnCh, columnGraph.Nums())

				doneCh := make(chan struct{})

				go func() {
					defer close(doneCh)
					for col := range columnCh {
						columnGraph.Update(col)
						time.Sleep(time.Millisecond * 800)
					}
				}()

				legend := []string{
					fmt.Sprintf("%s%s%s", utils.Green, utils.WhiteSquare+" Current", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Blue, utils.WhiteSquare+" Compare", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Orange, utils.WhiteSquare+" Locked", utils.Reset),
				}
				legendBox := component.NewBox(legend...).
					WithRoundedCorners().
					WithTitle("Legend").
					WithPadding(1, 2, 1, 2).
					// WithSpace(1).
					At(10, 70)

				component.Print(legendBox)

				go func() {
					<-doneCh
					component.Clear(legendBox)
				}()

			}

			//
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
