package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
	"tui/internal/component"
	term_utils "tui/internal/term-utils"

	"golang.org/x/term"
)

type State int

const (
	running State = iota
	paused
	next
	stop
	prev
	notStarted
	slower
	faster
	reset
)

var debug = term_utils.GetDebugFunc()

const (
	MIN_UPDATE_TIME   = time.Millisecond * 100
	MAX_UPDATE_TIME   = time.Millisecond * 950
	START_UPDATE_TIME = time.Millisecond * 500
)

func main() {
	fd, oldState, err := term_utils.Start()
	if err != nil {
		fmt.Println("error initializing raw terminal", err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)
	defer term_utils.TearDown()

	rows, cols, err := term_utils.GetSize(fd)
	if err != nil {
		fmt.Println("error getting terminal size", err)
		os.Exit(1)
	}

	instrCh := make(chan string)
	buf := component.NewBuf(instrCh)
	go buf.ReadLoop()
	go buf.FlushLoop()

	titleBox := component.NewBox(instrCh, "Terminal Algorithm Visualizer").
		WithColoredBorder(term_utils.Blue).
		WithRoundedCorners().
		WithPadding(1, 3).
		At(2, 34)
	titleBox.PrintIdle()

	algoList := component.NewList(
		instrCh,
		"Bubble sort",
		"Selection sort",
		"Insertion sort",
		"Quick sort",
		"Merge sort",
		"Heap sort").At(12, 5)
	algoList.PrintIdle()

	controlBoxContent := createControlBoxContent()
	controlBox := component.NewBox(instrCh, controlBoxContent...).
		WithTitle("Controls").
		WithColoredBorder(term_utils.Lightgray).
		WithRoundedCorners().
		At(30, 23)
	controlBox.PrintIdle()

	categoryTabs := component.NewTabs(instrCh, "Sorting", "Searching", "Graphs").At(7, 36)
	categoryTabs.PrintIdle()

	statusbar := component.NewStatusbar(instrCh, rows, cols, &buf)
	statusbar.PrintIdle()
	statusbar.Set("Welcome to the Terminal Algorithm Visualizer")

	nums := []int{1, 5, 8, 2, 11, 3, 12, 4, 9, 14, 13, 7, 15, 6, 10}
	var visualizer component.Visualizer
	visualizer = component.NewColumnGraph(instrCh, nums)
	visualizer.At(11, 32)
	visualizer.PrintIdle()

	speedBox := component.NewBox(instrCh, fmt.Sprintf("%d ms", START_UPDATE_TIME.Milliseconds())).
		WithRoundedCorners().WithTitle("Speed").At(20, 85)

	var legendBox *component.Box

	controlCh := make(chan State)
	defer close(controlCh)

	graphState := notStarted
	selected := ""

	reader := bufio.NewReader(os.Stdin)
	active := true

	for active {
		reader.Discard(reader.Buffered())
		b, err := reader.ReadByte()
		if err != nil {
			return
		}

		switch b {
		case 'q':
			// for final version do something here

			// statusbar.Set(component.Info("Goodbye"))
			// <-time.After(time.Millisecond * 500)
			return

		case 'j':
			algoList.Next()
			if graphState != notStarted {
				break
			}

			var newVisualizer component.Visualizer
			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
				newVisualizer = component.NewTree(instrCh, nums)
			} else {
				newVisualizer = component.NewColumnGraph(instrCh, nums)
			}

			if reflect.TypeOf(newVisualizer) == reflect.TypeOf(visualizer) {
				break
			}

			component.Clear(visualizer)
			visualizer = newVisualizer
			visualizer.At(11, 32)
			component.Print(visualizer)

		case 'k':
			algoList.Prev()
			if graphState != notStarted {
				break
			}

			var newVisualizer component.Visualizer
			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
				newVisualizer = component.NewTree(instrCh, nums)
			} else {
				newVisualizer = component.NewColumnGraph(instrCh, nums)
			}

			if reflect.TypeOf(newVisualizer) == reflect.TypeOf(visualizer) {
				break
			}

			component.Clear(visualizer)
			visualizer = newVisualizer
			visualizer.At(11, 32)
			component.Print(visualizer)

		case term_utils.Tab:
			categoryTabs.Next()
			statusbar.After(component.Info("There is nothing implemented yet"), time.Second*3)

		case term_utils.Enter:
			if graphState != notStarted {
				statusbar.After(component.Error(
					"First stop the current visualization with 'x' before selecting a new one"),
					time.Second*3)
				break
			}

			selected = algoList.SelectedValue()
			algoName := component.ToAlgoName(selected)
			if algoName == component.NotImplemented {
				statusbar.After(component.Info(
					"This algorithm is not implemented yet, sadge"),
					time.Second*3)
				break
			}

			algo := component.NewAlgorithm(algoName)
			statusbar.After(fmt.Sprintf("%s: Started", selected), time.Second*3)
			visualizer.Init(algo)

			speedBox.PrintIdle()
			legendBox = component.NewBox(instrCh, algo.Legend...).WithRoundedCorners().WithTitle("Legend").At(13, 82)
			legendBox.PrintIdle()

			graphState = running
			go handleGraph(&buf, controlCh, visualizer, speedBox)

		case term_utils.Space:
			if graphState == notStarted {
				break
			}

			if graphState == paused {
				statusbar.After(fmt.Sprintf("%s: Resumed", selected), time.Second*3)
				graphState = running
				controlCh <- running
			} else {
				statusbar.Set(fmt.Sprintf("%s: Paused - Press [ space ] to continue", selected))
				graphState = paused
				controlCh <- paused
			}

		case 'x':
			if graphState == notStarted {
				break
			}

			controlCh <- stop
			graphState = notStarted
			statusbar.After(fmt.Sprintf("%s: Stopped", selected), time.Second*3)

			component.Clear(legendBox)
		case 'n':
			if graphState == notStarted {
				break
			}

			graphState = paused
			controlCh <- next

		case 'p':
			if graphState == notStarted {
				break
			}

			graphState = paused
			controlCh <- prev

		case 'f':
			if graphState == notStarted {
				break
			}

			controlCh <- faster
			statusbar.After(fmt.Sprintf("%s: Faster", selected), time.Second)

		case 's':
			if graphState == notStarted {
				break
			}

			controlCh <- slower
			statusbar.After(fmt.Sprintf("%s: Slower", selected), time.Second)
		}

		<-time.After(time.Millisecond * 50)
	}

	term_utils.WaitForEnter()
}

func handleGraph(
	buf *component.Buf,
	controlCh <-chan State,
	visualizer component.Visualizer,
	speedBox *component.Box,
) {
	state := running
	waitTime := START_UPDATE_TIME
	for {
		select {
		case state = <-controlCh:
			switch state {
			case stop:
				component.Update(visualizer)
				component.Clear(speedBox)
				return
			case next:
				state = paused
				visualizer.Next()
			case prev:
				state = paused
				visualizer.Prev()
			case reset:
			case faster:
				if waitTime-time.Millisecond*50 >= MIN_UPDATE_TIME {
					waitTime -= time.Millisecond * 50
					speedBox.Update(int(waitTime.Milliseconds()))
				}
			case slower:
				if waitTime+time.Millisecond*50 <= MAX_UPDATE_TIME {
					waitTime += time.Millisecond * 50
					speedBox.Update(int(waitTime.Milliseconds()))
				}
			}
		default:
			if state == paused {
				break
			}
			visualizer.Next()
			time.Sleep(waitTime)
		}
	}
}

func createControlBoxContent() []string {
	content := []string{
		"  [ q ]    to quit the program",
		"  [ j ]    to select next item in list",
		"  [ k ]    to select prev item in list",
		"[ enter ]  to start algorithm visualization",
		"  [ n ]    to goto next step in visualization manually",
		"  [ p ]    to goto previous step in visualization manually",
		"  [ f ]    to speed up visualization",
		"  [ s ]    to slow down visualization",
		"[ space ]  to pause / resume visualization",
		"  [ x ]    to stop visualization",
		" [ tab ]   to switch algorithm category pages",
	}

	for i, s := range content {
		parts := strings.SplitAfter(s, "]")
		key, desc := parts[0], parts[1]
		content[i] = fmt.Sprintf(
			"%s%s",
			term_utils.Colorize(key, term_utils.Gray),
			term_utils.Colorize(desc, term_utils.Lightgray))
	}
	return content
}
