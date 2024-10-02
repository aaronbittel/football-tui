package component

import (
	"tui/internal/algorithm"
	term_utils "tui/internal/term-utils"
)

type Visualizer interface {
	Next() string
	Prev() string
	Init(algo *Algorithm)
	At(row, col int)
	Printer
	Clearer
}

type AlgorithmFn func(chan<- algorithm.ColumnGraphData, []int)

type AlgoName int

func (a AlgoName) String() string {
	switch a {
	case Bubble:
		return "Bubblesort"
	case Selection:
		return "Selectionsort"
	case Insertion:
		return "Insertionsort"
	case Quick:
		return "Quicksort"
	case Merge:
		return "Mergesort"
	case Heap:
		return "Heapsort"
	default:
		panic("not implemented")
	}
}

const (
	Bubble AlgoName = iota
	Selection
	Insertion
	Quick
	Merge
	Heap
	NotImplemented
)

type Algorithm struct {
	Name        AlgoName
	AlgorithmFn AlgorithmFn
	Legend      []string
}

func NewAlgorithm(name AlgoName) *Algorithm {
	var (
		algo   AlgorithmFn
		legend []string
	)

	switch name {
	case Bubble:
		algo = algorithm.Bubblesort
		legend = BubbleLegend
	case Selection:
		algo = algorithm.Selectionsort
		legend = SelectionLegend
	case Insertion:
		algo = algorithm.Insertionsort
		legend = InsertionLegend
	case Quick:
		algo = algorithm.Quicksort
		legend = QuickLegend
	case Merge:
		algo = algorithm.Mergesort
		legend = MergeLegend
	case Heap:
		algo = algorithm.Heapsort
		legend = HeapLegend
	}

	return &Algorithm{
		Name:        name,
		AlgorithmFn: algo,
		Legend:      legend,
	}
}

type columnParams struct {
	maxVal int
	spaces int
}

func ToAlgoName(s string) AlgoName {
	switch s {
	case "Bubble sort":
		return Bubble
	case "Selection sort":
		return Selection
	case "Insertion sort":
		return Insertion
	case "Quick sort":
		return Quick
	case "Merge sort":
		return Merge
	case "Heap sort":
		return Heap
	default:
		return NotImplemented
	}
}

var (
	BubbleLegend = []string{
		term_utils.Colorize("▣  Current", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange),
	}
	SelectionLegend = []string{
		term_utils.Colorize("▣  Lowest", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange),
	}
	QuickLegend = []string{
		term_utils.Colorize("▣  Current", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange),
		term_utils.Colorize("▣  !interesting", term_utils.Lightgray),
	}
	InsertionLegend = []string{
		term_utils.Colorize("▣  Current Value", term_utils.Green),
		term_utils.Colorize("▣  Current Array", term_utils.Lightgray),
	}
	MergeLegend = []string{
		term_utils.Colorize("▣  Left side", term_utils.Green),
		term_utils.Colorize("▣  Right side", term_utils.Blue),
		term_utils.Colorize("▣  !interesting", term_utils.Lightgray),
	}
	HeapLegend = []string{
		term_utils.Colorize("▣  Implement Me", term_utils.BoldRed),
	}
)
