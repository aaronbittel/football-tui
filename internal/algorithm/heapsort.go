package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Heapsort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)
	columnCh <- NewColumnGraphData(slices.Clone(nums), map[int]string{}, "Building the heap")

	n := len(nums)
	for start := (n - 2) / 2; start >= 0; start-- {

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			map[int]string{
				start: term_utils.BoldGreen,
			},
			fmt.Sprintf("Heapify starting at value %d", nums[start]))

		siftDown(columnCh, nums, start, n-1)

	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		map[int]string{},
		"Heapifying process complete")

	for end := n - 1; end > 0; end-- {

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			map[int]string{
				0:   term_utils.BoldGreen,
				end: term_utils.BoldBlue,
			},
			fmt.Sprintf("Swapping root with value %d", nums[end]))

		swap(nums, 0, end)

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			map[int]string{},
			fmt.Sprintf("Array after swap: %v", nums))

		siftDown(columnCh, nums, 0, end-1)

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			map[int]string{},
			fmt.Sprintf("Heap after siftDown: %v", nums))
	}
}

func siftDown(columnCh chan<- ColumnGraphData, nums []int, start, end int) {
	root := start
	for root*2+1 <= end {

		child := root*2 + 1

		colors := map[int]string{
			root:  term_utils.BoldGreen,
			child: term_utils.BoldBlue,
		}

		if child+1 <= end && nums[child] < nums[child+1] {
			child++
			colors[child] = term_utils.BoldBlue
		}

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Checking if current root node %d is greater than leaf nodes", nums[root]))

		if nums[root] >= nums[child] {

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf(
					"Root %d is the greatest value", nums[root]))

			return
		}

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Swap root %d and leaf %d because %d < %d",
				nums[root],
				nums[child],
				nums[root],
				nums[child]))

		swap(nums, root, child)

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Sifting down: swapping index %d (value %v) with index %d (value %v)",
				root,
				nums[root],
				child,
				nums[child]))

		root = child
	}
}

func swap(a []int, i, j int) {
	a[i], a[j] = a[j], a[i]
}
