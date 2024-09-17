package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
	algo "tui/internal/algorithms"
	"tui/internal/component"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

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
		"SelectionSort",
		"HeapSort",
		"MergeSort").
		At(9, 2)

	algorithms := [][]string{
		{"Bubblesort", "Quicksort", "selectionsort", "Heapsort", "Mergesort"},
		{"Binary", "Linear"},
		{"Breadth-First", "Depth-First"},
		{"SampleStack1", "SampleStack2", "SampleStack3"},
		{"SampleArray1", "SampleArray2", "SampleArray3"},
		{"SampleHashSet1", "SampleHashSet2", "SampleHashSet3"},
	}

	nums := []int{14, 4, 12, 1, 16, 6, 13, 8, 11, 17, 7, 15, 2, 9, 18, 3, 5, 10}
	// nums := []int{4, 8, 2, 1, 6, 7, 3}
	utils.MoveCursor(40, 79)
	columnGraph := component.NewColumnGraph(component.NewColumn(nums, nil, "")).At(10, 25)

	component.Print(columnGraph)
	component.Print(box)
	component.Print(tabs)
	component.Print(list)

	reader := bufio.NewReader(os.Stdin)
	running := true

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
			// col := <-columnCh
			// columnGraph.Update(col)
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
				legend := []string{
					fmt.Sprintf("%s%s%s", utils.Green, utils.WhiteSquare+" Current", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Blue, utils.WhiteSquare+" Compare", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Orange, utils.WhiteSquare+" Locked", utils.Reset),
				}

				handleGraph(columnGraph, algo.Bubblesort, time.Millisecond*200, legend)

			} else if tabs.Selected == 0 && list.Selected == 2 {

				legend := []string{
					fmt.Sprintf("%s%s%s", utils.Green, utils.WhiteSquare+" Current Lowest", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Blue, utils.WhiteSquare+" Compare", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Orange, utils.WhiteSquare+" Locked", utils.Reset),
				}

				handleGraph(columnGraph, algo.Selectionsort, time.Millisecond*200, legend)

			} else if tabs.Selected == 0 && list.Selected == 1 {

				legend := []string{
					fmt.Sprintf("%s%s%s", utils.Green, utils.WhiteSquare+" Pivot", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Blue, utils.WhiteSquare+" Compare", utils.Reset),
					fmt.Sprintf("%s%s%s", utils.Orange, utils.WhiteSquare+" Locked", utils.Reset),
				}

				handleGraph(columnGraph, algo.Quicksort, time.Millisecond*400, legend)
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

type Algorithm func(chan<- component.Column, []int)

func handleGraph(columnGraph *component.ColumnGraph, algo Algorithm, waitTime time.Duration, legend []string) {
	columnCh := make(chan component.Column)
	nums := columnGraph.Nums()
	go algo(columnCh, nums)

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)
		for col := range columnCh {
			columnGraph.Update(col)
			time.Sleep(waitTime)
		}
	}()

	_, cWidth := columnGraph.Mask()
	cRow, cCol := columnGraph.Pos()
	legendRow, legendCol := cRow, cCol+cWidth+7

	legendBox := component.NewBox(legend...).
		WithRoundedCorners().
		WithTitle("Legend").
		WithPadding(1, 2, 1, 2).
		At(legendRow, legendCol)

	component.Print(legendBox)

	go func() {
		<-doneCh
		component.Clear(legendBox)
	}()
}
