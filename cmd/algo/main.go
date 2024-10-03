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

	buf := component.NewBuf()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 30)
		for {
			buf.Flush()
			<-ticker.C
		}
	}()

	titleBox := component.NewBox("Terminal Algorithm Visualizer").
		WithColoredBorder(term_utils.Blue).
		WithRoundedCorners().
		WithPadding(1, 3).
		At(2, 34)
	buf.Write(component.Print(titleBox))

	algoList := component.NewList(
		"Bubble sort",
		"Selection sort",
		"Insertion sort",
		"Quick sort",
		"Merge sort",
		"Heap sort").At(12, 5)
	buf.Write(component.Print(algoList))

	controlBoxContent := createControlBoxContent()
	controlBox := component.NewBox(controlBoxContent...).
		WithTitle("Controls").
		WithColoredBorder(term_utils.Lightgray).
		WithRoundedCorners().
		At(30, 23)
	buf.Write(component.Print(controlBox))

	categoryTabs := component.NewTabs("Sorting", "Searching", "Graphs").At(7, 36)
	buf.Write(component.Print(categoryTabs))

	statusbar := component.NewStatusbar(rows, cols, &buf)
	buf.Write(component.Print(statusbar))
	buf.Write(statusbar.Set("Welcome to the Terminal Algorithm Visualizer"))

	nums := []int{1, 5, 8, 2, 11, 3, 12, 4, 9, 14, 13, 7, 15, 6, 10}
	var visualizer component.Visualizer
	visualizer = component.NewColumnGraph(nums)
	visualizer.At(11, 32)
	buf.Write(component.Print(visualizer))

	speedBox := component.NewBox(fmt.Sprintf("%d ms", START_UPDATE_TIME.Milliseconds())).
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

			// buf.Write(statusbar.Set(component.Info("Goodbye")))
			// <-time.After(time.Millisecond * 500)
			return

		case 'j':
			buf.Write(algoList.Next())
			if graphState != notStarted {
				break
			}

			var newVisualizer component.Visualizer
			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
				newVisualizer = component.NewTree(nums)
			} else {
				newVisualizer = component.NewColumnGraph(nums)
			}

			if reflect.TypeOf(newVisualizer) == reflect.TypeOf(visualizer) {
				break
			}

			buf.Write(component.Clear(visualizer))
			visualizer = newVisualizer
			visualizer.At(11, 32)
			buf.Write(component.Print(visualizer))

		case 'k':
			buf.Write(algoList.Prev())
			if graphState != notStarted {
				break
			}

			var newVisualizer component.Visualizer
			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
				newVisualizer = component.NewTree(nums)
			} else {
				newVisualizer = component.NewColumnGraph(nums)
			}

			if reflect.TypeOf(newVisualizer) == reflect.TypeOf(visualizer) {
				break
			}

			buf.Write(component.Clear(visualizer))
			visualizer = newVisualizer
			visualizer.At(11, 32)
			buf.Write(component.Print(visualizer))

		case term_utils.Tab:
			buf.Write(categoryTabs.Next())
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

			buf.Write(component.Print(speedBox))
			legendBox = component.NewBox(algo.Legend...).WithRoundedCorners().WithTitle("Legend").At(13, 82)
			buf.Write(component.Print(legendBox))

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
				buf.Write(statusbar.Set(fmt.Sprintf(
					"%s: Paused - Press [ space ] to continue", selected)))
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

			buf.Write(component.Clear(legendBox))
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
				buf.Write(component.Update(visualizer))
				buf.Write(component.Clear(speedBox))
				return
			case next:
				state = paused
				buf.Write(visualizer.Next())
			case prev:
				state = paused
				buf.Write(visualizer.Prev())
			case reset:
			case faster:
				if waitTime-time.Millisecond*50 >= MIN_UPDATE_TIME {
					waitTime -= time.Millisecond * 50
					buf.Write(speedBox.Update(int(waitTime.Milliseconds())))
				}
			case slower:
				if waitTime+time.Millisecond*50 <= MAX_UPDATE_TIME {
					waitTime += time.Millisecond * 50
					buf.Write(speedBox.Update(int(waitTime.Milliseconds())))
				}
			}
		default:
			if state == paused {
				break
			}
			buf.Write(visualizer.Next())
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
