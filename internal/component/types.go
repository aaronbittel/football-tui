package component

import (
	"errors"
	"fmt"
	"tui/internal/algorithm"
	term_utils "tui/internal/term-utils"
)

type Algorithm interface {
	GetFrames(columnCh chan algorithm.ColumnGraphData, nums []int) []algorithm.ColumnGraphData
	fmt.Stringer
	Legend() []string
}

type Visualizer interface {
	Next()
	Prev()
	Init(Algorithm)
	At(row, col int)
	Printer
	Clearer
}

var (
	BubbleSort    = NewSortAlgorithm("Bubble sort", algorithm.Bubblesort, bubbleLegend)
	InsertionSort = NewSortAlgorithm("Insertion sort", algorithm.Insertionsort, insertionLegend)
	SelectionSort = NewSortAlgorithm("Selection sort", algorithm.Selectionsort, selectionLegend)
	QuickSort     = NewSortAlgorithm("Quick sort", algorithm.Quicksort, quickLegend)
	MergeSort     = NewSortAlgorithm("Merge sort", algorithm.Mergesort, mergeLegend)
	HeapSort      = NewSortAlgorithm("Heap sort", algorithm.Heapsort, heapLegend)

	LinearSearch = createSearchAlgorithm("Linear search", algorithm.LinearSearch, linearLegend)
	BinarySearch = createSearchAlgorithm("Binary search", algorithm.BinarySearch, binaryLegend)
	JumpSearch   = createSearchAlgorithm("Jump search", algorithm.JumpSearch, jumpLegend)

	NotImplementedErr = errors.New("This algorithm is not yet implemented")
)

// FIX: Do I need the err?
func GetSortAlgoByName(name string) (Algorithm, error) {
	switch name {
	case "Bubble sort":
		return BubbleSort, nil
	case "Insertion sort":
		return InsertionSort, nil
	case "Selection sort":
		return SelectionSort, nil
	case "Quick sort":
		return QuickSort, nil
	case "Merge sort":
		return MergeSort, nil
	case "Heap sort":
		return HeapSort, nil
	default:
		return nil, NotImplementedErr
	}
}

func GetSearchAlgoByName(name string, target int) (Algorithm, error) {
	switch name {
	case "Linear search":
		return NewSearchAlgorithm(LinearSearch, target), nil
	case "Binary search":
		return NewSearchAlgorithm(BinarySearch, target), nil
	case "Jump search":
		return NewSearchAlgorithm(JumpSearch, target), nil
	default:
		return nil, NotImplementedErr
	}
}

type SortAlgorithm struct {
	name     string
	sortFunc SortFunc
	legend   []string
}

func NewSortAlgorithm(name string, fn SortFunc, legend []string) *SortAlgorithm {
	return &SortAlgorithm{
		name:     name,
		sortFunc: fn,
		legend:   legend,
	}
}

func (sa SortAlgorithm) Legend() []string {
	return sa.legend
}

func (sa SortAlgorithm) String() string {
	return sa.name
}

func (sa SortAlgorithm) GetFrames(columnCh chan algorithm.ColumnGraphData, nums []int) []algorithm.ColumnGraphData {

	go sa.sortFunc(columnCh, nums)
	frames := make([]algorithm.ColumnGraphData, 0, 300)
	for col := range columnCh {
		frames = append(frames, col)
	}
	return frames
}

type SearchAlgorithm struct {
	name     string
	sortFunc SearchFunc
	target   int
	legend   []string
}

func NewSearchAlgorithm(sa *SearchAlgorithm, target int) *SearchAlgorithm {
	sa.SetTarget(target)
	return sa
}

func createSearchAlgorithm(name string, fn SearchFunc, legend []string) *SearchAlgorithm {
	return &SearchAlgorithm{
		name:     name,
		sortFunc: fn,
		legend:   legend,
	}
}

func (sa SearchAlgorithm) Legend() []string {
	return sa.legend
}

func (sa *SearchAlgorithm) SetTarget(target int) {
	sa.target = target
}

func (sa SearchAlgorithm) String() string {
	return sa.name
}

func (sa SearchAlgorithm) GetFrames(columnCh chan algorithm.ColumnGraphData, nums []int) []algorithm.ColumnGraphData {

	go sa.sortFunc(columnCh, nums, sa.target)
	frames := make([]algorithm.ColumnGraphData, 0, 300)
	for col := range columnCh {
		frames = append(frames, col)
	}

	return frames
}

type SortFunc func(columnCh chan<- algorithm.ColumnGraphData, nums []int)
type SearchFunc func(columnCh chan<- algorithm.ColumnGraphData, nums []int, target int)

var (
	bubbleLegend = []string{
		term_utils.Colorize("▣  Current", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange)}

	selectionLegend = []string{
		term_utils.Colorize("▣  Lowest", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange)}

	quickLegend = []string{
		term_utils.Colorize("▣  Current", term_utils.Green),
		term_utils.Colorize("▣  Compare", term_utils.Blue),
		term_utils.Colorize("▣  Locked", term_utils.Orange),
		term_utils.Colorize("▣  !interesting", term_utils.Lightgray)}

	insertionLegend = []string{
		term_utils.Colorize("▣  Current Value", term_utils.Green),
		term_utils.Colorize("▣  Current Array", term_utils.Lightgray)}

	mergeLegend = []string{
		term_utils.Colorize("▣  Left side", term_utils.Green),
		term_utils.Colorize("▣  Right side", term_utils.Blue),
		term_utils.Colorize("▣  !interesting", term_utils.Lightgray)}

	heapLegend = []string{
		term_utils.Colorize("▣  Implement Me", term_utils.BoldRed)}

	linearLegend = []string{
		term_utils.Colorize("▣  Implement Me", term_utils.BoldRed)}

	binaryLegend = []string{
		term_utils.Colorize("▣  Implement Me", term_utils.BoldRed)}

	jumpLegend = []string{
		term_utils.Colorize("▣  Implement Me", term_utils.BoldRed)}
)
