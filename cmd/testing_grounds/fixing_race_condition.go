package main

import (
	"slices"
	"time"
	"tui/internal/component"
	term_utils "tui/internal/term-utils"
)

func fixingRaceCondition() {
	nums := []int{1, 4, 3, 14, 5, 12}

	columnGraph := component.NewColumnGraph(nums).At(5, 30)
	columnCh := make(chan component.ColumnGraphData)

	time.Sleep(time.Second)
	component.Print(columnGraph)

	go component.Bubblesort(columnCh, slices.Clone(nums))

	func() {
		doneCh := make(chan struct{})
		go func() {
			defer close(doneCh)
			for col := range columnCh {
				columnGraph.Update(col)
				time.Sleep(time.Millisecond * 300)
			}
		}()

		row, col := columnGraph.Pos()
		_, width := columnGraph.Mask()
		box := component.NewBox("Hallo", "Tehllo", "Hello").WithColoredBorder(term_utils.Lightgray).WithRoundedCorners().At(row, col+width+5)
		component.Print(box)

		<-doneCh
	}()
}
