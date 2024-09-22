package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"tui/internal/component"
	term_utils "tui/internal/term-utils"
)

var debug = term_utils.GetDebugFunc()

func Hello() {
	fmt.Println("Hello")
}

func TiminingVisualization() {
	list := component.NewList("Bubble sort", "Selection sort", "Insertion sort").At(5, 5)
	component.Print(list)

	controlCh := make(chan string)
	defer close(controlCh)

	nums := []int{1, 4, 2, 8, 3, 14, 9, 6, 5, 12}
	columnGraph := component.NewColumnGraph(nums).At(5, 30)
	columnCh := make(chan component.ColumnGraphData)

	component.Print(columnGraph)

	reader := bufio.NewReader(os.Stdin)
	running := true
	algoRunning := true

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', term_utils.CtrlC:
			running = false
		case term_utils.Enter:
			debug(fmt.Sprintf("%s task is running", "Bubblesort"))
			go component.Bubblesort(columnCh, nums)
			go HandleAlgo(controlCh, columnGraph, columnCh)
		// case 'x':
		// 	task = stopTask(controlCh, task)
		case term_utils.Space:
			if algoRunning {
				controlCh <- "PAUSED"
				debug(fmt.Sprintf("%s task is paused", "Bubblesort"))
				algoRunning = false
			} else {
				controlCh <- "RUNNING"
				debug(fmt.Sprintf("%s task is running", "Bubblesort"))
				algoRunning = true
			}
		case 'n':
			controlCh <- "PAUSED"
			controlCh <- "NEXT"
			algoRunning = false
		case 'p':
			controlCh <- "PAUSED"
			controlCh <- "PREV"
			algoRunning = false
		case 'j':
			list.Next()
		case 'k':
			list.Prev()
			// case 'f':
			// 	sendCtrl(controlCh, task, "FASTER")
			// case 's':
			// 	sendCtrl(controlCh, task, "SLOWER")
		}
	}
}

func HandleAlgo(
	controlCh <-chan string,
	columnGraph *component.ColumnGraph,
	columnCh <-chan component.ColumnGraphData,
) {
	state := "RUNNING"
	waitTime := time.Millisecond * 500
	graphStates := make([]component.ColumnGraphData, 0)
	cursor := 0
	var col component.ColumnGraphData

outer:
	for {
		select {
		case state = <-controlCh:
			debug("len states", len(graphStates), "cursor", cursor)
			switch state {
			case "STOP":
				break outer
			case "NEXT":
				state = "PAUSED"
				cursor++
				if cursor < len(graphStates) {
					col = graphStates[cursor]
				} else {
					col = <-columnCh
					graphStates = append(graphStates, col)
				}
				columnGraph.Update(col)
				time.Sleep(time.Millisecond * 5)
			case "PREV":
				state = "PAUSED"
				if cursor <= 0 {
					break
				}
				cursor--
				columnGraph.Update(graphStates[cursor])
				time.Sleep(time.Millisecond * 5)
			case "FASTER":
				debug("FASTER", waitTime)
				if waitTime-time.Millisecond*50 >= time.Millisecond*50 {
					waitTime -= time.Millisecond * 50
				}
			case "SLOWER":
				debug("SLOWER", waitTime)
				if waitTime+time.Millisecond*50 <= time.Millisecond*1500 {
					waitTime += time.Millisecond * 50
				}
			}
		default:
			if state == "PAUSED" {
				break
			}
			cursor++
			if cursor < len(graphStates) {
				col = graphStates[cursor]
			} else {
				col = <-columnCh
				graphStates = append(graphStates, col)
			}
			columnGraph.Update(col)
			time.Sleep(waitTime)

		}
	}

	debug("FINISHED")
}
