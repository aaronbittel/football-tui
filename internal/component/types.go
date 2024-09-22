package component

import (
	term_utils "tui/internal/term-utils"
)

type AlgorithmFn func(chan<- ColumnGraphData, []int)

type AlgoName int

const (
	Bubble AlgoName = iota
	Selection
	Quick
)

type Algorithm struct {
	AlgorithmFn AlgorithmFn
	legend      []string
}

func NewAlgorithm(name AlgoName) *Algorithm {
	var (
		algo   AlgorithmFn
		legend []string
	)

	switch name {
	case Bubble:
		algo = Bubblesort
		legend = bubbleLegend
	case Selection:
		algo = Selectionsort
		legend = selectionLegend
	case Quick:
		algo = Quicksort
		legend = quickLegend
	}

	return &Algorithm{
		AlgorithmFn: algo,
		legend:      legend,
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
	case "Quick sort":
		return Quick
	default:
		panic("unknown algorithm")
	}
}

var (
	bubbleLegend = []string{
		term_utils.Colorize("▣ Current", term_utils.Green),
		term_utils.Colorize("▣ Compare", term_utils.Blue),
		term_utils.Colorize("▣ Locked", term_utils.Orange),
	}
	selectionLegend = []string{
		term_utils.Colorize("▣ Lowest", term_utils.Green),
		term_utils.Colorize("▣ Compare", term_utils.Blue),
		term_utils.Colorize("▣ Locked", term_utils.Orange),
	}
	quickLegend = []string{
		term_utils.Colorize("▣ Current", term_utils.Green),
		term_utils.Colorize("▣ Compare", term_utils.Blue),
		term_utils.Colorize("▣ Locked", term_utils.Orange),
		term_utils.Colorize("▣ !interesting", term_utils.Lightgray),
	}
)
