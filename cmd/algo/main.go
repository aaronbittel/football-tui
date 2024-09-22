package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
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

	// nums := []int{14, 4, 12, 1, 16, 6, 13, 8, 11, 17, 7, 15, 2, 9, 18, 3, 5, 10}
	nums := []int{14, 12, 1, 8, 11, 15, 2, 3, 5}

	columnGraph := component.NewColumnGraph(slices.Clone(nums)).At(10, 25)

	component.Print(columnGraph)
	// component.Print(title)
	// component.Print(tabs)
	// component.Print(list)
	// component.Print(controlBox)

	controlCh := make(chan string)
	defer close(controlCh)
	// visualizationRunning := false

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
		case utils.Enter:
			legend := []string{
				utils.Colorize(fmt.Sprintf("%s Current", utils.WhiteSquare), utils.Green),
				utils.Colorize(fmt.Sprintf("%s Compare", utils.WhiteSquare), utils.Blue),
				utils.Colorize(fmt.Sprintf("%s Locked", utils.WhiteSquare), utils.Orange),
			}
			handleGraph(controlCh, columnGraph, algo.Bubblesort, time.Millisecond*100, legend)
		}
	}
}

func getRandomColor() string {
	const colorCode = "\033[%dm"
	r := rand.IntN(8) + 30
	return fmt.Sprintf(colorCode, r)
}

func stopVisualization(controlCh chan<- string) {
	controlCh <- "STOP"
}

type Algorithm func(chan<- component.ColumnGraphData, []int)

func handleGraph(
	controlCh <-chan string,
	columnGraph *component.ColumnGraph,
	algo Algorithm,
	waitTime time.Duration,
	legend []string,
) {

	columnCh := make(chan component.ColumnGraphData)

	go algo(columnCh, slices.Clone(columnGraph.Nums()))

	go func() {
		for col := range columnCh {
			columnGraph.Update(col)
			time.Sleep(waitTime)
		}
	}()

	_, cWidth := columnGraph.Mask()
	cRow, cCol := columnGraph.Pos()
	legendRow, legendCol := cRow, cCol+cWidth+20

	legendBox := component.NewBox(legend...).
		WithRoundedCorners().
		WithTitle("Legend").
		WithPadding(1, 2, 1, 2).
		At(legendRow, legendCol)

	component.Print(legendBox)

}
