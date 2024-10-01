package main

import (
	"fmt"
	"math/rand/v2"
	"os"
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

	title := component.NewBox("Terminal Algorithm Visualizer").
		WithColoredBorder(term_utils.Blue).
		WithRoundedCorners().
		WithPadding(1, 3).
		At(2, 25)

	for i, line := range strings.Split(title.String(), "\n") {
		buf.Write(term_utils.MoveCur(2+i, 25))
		buf.Write(line)
	}

	list := component.NewList("Aalskdfj", "ABsfdasjfd", "lksjalkfj", "KLJLKJG").At(7, 5)
	for i, line := range strings.Split(list.String(), "\n") {
		buf.Write(term_utils.MoveCur(7+i, 5))
		buf.Write(line)
	}

	tabs := component.NewTabs("Aalskdfj", "ABsfdasjfd", "lksjalkfj", "KLJLKJG").At(12, 24)
	for i, line := range strings.Split(tabs.String(), "\n") {
		buf.Write(term_utils.MoveCur(12+i, 24))
		buf.Write(line)
	}
	nums := []int{1, 5, 8, 2, 11, 3, 12, 4, 9, 14, 13, 7, 15, 6, 10}
	columnGraph := component.NewColumnGraph(nums)
	columnGraph.At(18, 28)
	columnGraph.Init(component.NewAlgorithm(component.Bubble))

	statusbar := component.NewStatusbar(rows, cols)
	for _, line := range strings.Split(statusbar.String(), "\n") {
		buf.Write(term_utils.MoveCur(rows-1, cols))
		buf.Write(line)
	}
	buf.Write(statusbar.Set("Welcome to the Terminal Algorithm Visualizer"))

	go func() {
		ticker := time.NewTicker(time.Millisecond * 450)
		defer ticker.Stop()
		for {
			<-ticker.C
			buf.Write(columnGraph.Next())
		}
	}()

	for i, line := range strings.Split(columnGraph.Idle(), "\n") {
		buf.Write(term_utils.MoveCur(18+i, 28))
		buf.Write(line)
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		defer ticker.Stop()
		for {
			if rand.IntN(2) == 0 {
				buf.Write(list.Next())
			} else {
				buf.Write(list.Prev())
			}
			<-ticker.C
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		defer ticker.Stop()
		for {
			buf.Write(tabs.Next())
			<-ticker.C
		}
	}()

	term_utils.WaitForEnter()
}

func handleGraph(
	controlCh <-chan State,
	visualizer component.Visualizer,
	waitTime time.Duration,
) {
	state := running
	for {
		select {
		case state = <-controlCh:
			switch state {
			case stop:
				visualizer.Reset()
				return
			case next:
				state = paused
				visualizer.Next()
			case prev:
				state = paused
				visualizer.Prev()
			case reset:
			case faster:
				if waitTime-time.Millisecond*50 >= time.Millisecond*50 {
					waitTime -= time.Millisecond * 50
				}
			case slower:
				if waitTime+time.Millisecond*50 <= time.Millisecond*1500 {
					waitTime += time.Millisecond * 50
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

// func old() {
// 	// statusbar := component.NewStatusbar(width, height)
// 	// statusbar.PrintIdle()
// 	// statusbar.Set("Welcome to Terminal Algorithm Visualizer")
//
// 	algoList := component.NewList(
// 		"Bubble sort",
// 		"Selection sort",
// 		"Insertion sort",
// 		"Quick sort",
// 		"Merge sort",
// 		"Heap sort").At(8, 5)
//
// 	component.Print(algoList)
//
// 	// nums := []int{14, 4, 12, 1, 16, 6, 13, 8, 11, 17, 7, 15, 2, 9, 18, 3, 5, 10}
// 	nums := []int{14, 4, 12, 1, 6, 13, 8, 11, 7, 2, 9, 15, 3, 5, 10}
// 	var visualizer component.Visualizer
// 	visualizer = component.NewColumnGraph(nums)
// 	visualizer.At(8, 25)
// 	visualizer.PrintIdle()
//
// 	controlCh := make(chan State)
// 	defer close(controlCh)
//
// 	graphState := notStarted
// 	selected := ""
//
// 	reader := bufio.NewReader(os.Stdin)
// 	active := true
//
// 	for active {
// 		b, err := reader.ReadByte()
// 		if err != nil {
// 			fmt.Println("Error reading byte:", err)
// 			break
// 		}
//
// 		switch b {
// 		case 'q', term_utils.CtrlC:
// 			active = false
// 		case 'j':
// 			algoList.Next()
// 			if graphState != notStarted {
// 				break
// 			}
// 			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
// 				visualizer = component.NewTree(nums)
// 			} else {
// 				visualizer = component.NewColumnGraph(nums)
// 			}
// 			visualizer.At(8, 25)
//
// 			for i := range 20 {
// 				term_utils.ClearLine(8+i, 25)
// 			}
//
// 			visualizer.PrintIdle()
// 		case 'x':
// 			controlCh <- stop
// 			graphState = notStarted
// 			// statusbar.After(fmt.Sprintf("%s: Stopped", selected), time.Second*3)
// 		case 'k':
// 			algoList.Prev()
// 			if graphState != notStarted {
// 				break
// 			}
// 			if component.ToAlgoName(algoList.SelectedValue()) == component.Heap {
// 				visualizer = component.NewTree(nums)
// 			} else {
// 				visualizer = component.NewColumnGraph(nums)
// 			}
// 			visualizer.At(8, 25)
//
// 			for i := range 20 {
// 				term_utils.ClearLine(8+i, 25)
// 			}
//
// 			visualizer.PrintIdle()
// 		case term_utils.Enter:
// 			if graphState != notStarted {
// 				// statusbar.Set(component.Error("First stop the current visualization with 'x' before selecting a new one"))
// 				break
// 			}
//
// 			selected = algoList.SelectedValue()
// 			algoName := component.ToAlgoName(selected)
// 			if algoName == component.NotImplemented {
// 				// statusbar.After(component.Info("This algorithm is not implemented yet, sadge"), time.Second*3)
// 				break
// 			}
//
// 			algo := component.NewAlgorithm(algoName)
// 			// statusbar.After(fmt.Sprintf("%s: Started", selected), time.Second*3)
//
// 			if algoName == component.Heap {
// 				visualizer = component.NewTree(nums)
// 			} else {
// 				visualizer = component.NewColumnGraph(nums)
// 			}
//
// 			visualizer.At(8, 25)
// 			visualizer.PrintIdle()
//
// 			doneCh := make(chan struct{})
// 			go func() {
// 				defer close(doneCh)
// 				visualizer.Init(algo)
// 			}()
//
// 			<-doneCh
// 			graphState = running
// 			go handleGraph(controlCh, visualizer, time.Millisecond*600)
// 		case term_utils.Space:
// 			if graphState == paused {
// 				// statusbar.After(fmt.Sprintf("%s: Resumed", selected), time.Second*3)
// 				graphState = running
// 				controlCh <- running
// 			} else {
// 				// statusbar.Set(fmt.Sprintf("%s: Paused - Press [ space ] to continue", selected))
// 				graphState = paused
// 				controlCh <- paused
// 			}
// 		case 'n':
// 			graphState = paused
// 			controlCh <- next
// 		case 'p':
// 			graphState = paused
// 			controlCh <- prev
// 		case 'f':
// 			controlCh <- faster
// 			// statusbar.After(fmt.Sprintf("%s: Faster", selected), time.Second)
// 		case 's':
// 			controlCh <- slower
// 			// statusbar.After(fmt.Sprintf("%s: Slower", selected), time.Second)
// 		}
// 	}
// }
