package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
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
	not_started
	slower
	faster
	reset
)

var (
	debug  = term_utils.GetDebugFunc()
	logger *slog.Logger
)

func main() {
	f, err := os.Create("logging.txt")
	if err != nil {
		fmt.Println("error creating logging file")
		os.Exit(1)
	}
	logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))

	fd, oldState, err := term_utils.Start()
	if err != nil {
		fmt.Println("error initializing raw terminal", err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)
	defer term_utils.TearDown()

	title := component.NewBox("Terminal Algorithm Visualizer").
		WithRoundedCorners().
		At(3, 35)

	algoList := component.NewList("Bubble sort", "Selection sort", "Quick sort", "Merge sort", "Heap sort").
		At(8, 5)

	component.Print(title)
	component.Print(algoList)

	nums := []int{14, 4, 12, 1, 16, 6, 13, 8, 11, 17, 7, 15, 2, 9, 18, 3, 5, 10}
	// nums := []int{14, 12, 1, 8, 11, 15, 2, 3, 5}

	columnGraph := component.NewColumnGraphFrames(nums).At(8, 25)
	columnGraph.PrintIdle()

	controlCh := make(chan State)
	graphState := not_started

	reader := bufio.NewReader(os.Stdin)
	active := true

	for active {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', term_utils.CtrlC:
			active = false
		case 'j':
			algoList.Next()
		case 'x':
			logger.Debug("pressed x")
			controlCh <- stop
			graphState = not_started
			columnGraph.Clear()
			columnGraph.PrintIdle()
		case 'k':
			algoList.Prev()
		case term_utils.Enter:
			logger.Debug("pressed enter")
			if graphState != not_started {
				controlCh <- stop
			}

			algo := component.NewAlgorithm(
				component.ToAlgoName(algoList.SelectedValue()),
			)

			doneCh := make(chan struct{})
			go func() {
				columnGraph.Init(algo)
				doneCh <- struct{}{}
			}()

			<-doneCh
			graphState = running
			go handleGraph(controlCh, columnGraph, time.Millisecond*400)
		case term_utils.Space:
			logger.Debug("pressed space")
			if graphState == paused {
				graphState = running
				controlCh <- running
			} else {
				graphState = paused
				controlCh <- paused
			}
		case 'n':
			logger.Debug("pressed n")
			graphState = paused
			controlCh <- next
		case 'p':
			logger.Debug("pressed p")
			graphState = paused
			controlCh <- prev
		case 'f':
			controlCh <- faster
		case 's':
			controlCh <- slower
		}
	}
}

func handleGraph(
	controlCh <-chan State,
	columnGraph *component.ColumnGraphFrames,
	waitTime time.Duration,
) {
	state := running
	for {
		select {
		case state = <-controlCh:
			switch state {
			case stop:
				logger.Debug("got stop")
				columnGraph.ClearDescription()
				debug("CLEARING DESC")
				return
			case next:
				logger.Debug("got next")
				state = paused
				columnGraph.Next()
			case prev:
				logger.Debug("got prev")
				state = paused
				columnGraph.Prev()
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
			columnGraph.Next()
			time.Sleep(waitTime)
		}
	}
}
