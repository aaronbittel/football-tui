package algorithms

import (
	"fmt"
	"maps"
	"slices"
	component "tui/internal/component"
	utils "tui/internal/term-utils"
)

func partition(
	columnCh chan<- component.Column,
	nums []int,
	low, high int,
	locked map[int]string,
) int {
	pivot := nums[high]
	i := low - 1

	locked[high] = utils.Green

	for j := low; j < high; j++ {
		locked[j] = utils.Blue
		columnCh <- component.NewColumn(
			slices.Clone(nums),
			maps.Clone(locked),
			fmt.Sprintf(
				"Comparing %s to Pivot %s",
				utils.Colorize(fmt.Sprintf("%d", nums[j]), utils.Blue),
				utils.Colorize(fmt.Sprintf("%d", high), utils.Green),
			),
		)
		if nums[j] < pivot {
			i++
			delete(locked, j)
			nums[i], nums[j] = nums[j], nums[i]
			locked[i] = utils.Blue
			columnCh <- component.NewColumn(
				slices.Clone(nums),
				maps.Clone(locked),
				"[Blue] is smaller than Pivot [Green] → Swap",
			)
			delete(locked, i)
		}
		delete(locked, j)
	}

	delete(locked, high)
	nums[i+1], nums[high] = nums[high], nums[i+1]
	locked[i+1] = utils.Green
	columnCh <- component.NewColumn(slices.Clone(nums), maps.Clone(locked), "Swap Pivot to correct position")
	delete(locked, i+1)
	return i + 1
}

func quicksortHelper(
	columnCh chan<- component.Column,
	nums []int,
	low, high int,
	locked map[int]string,
) {
	if low >= high {
		for i := high; i <= low; i++ {
			locked[i] = utils.Orange
		}
		columnCh <- component.NewColumn(slices.Clone(nums), locked, "Mark numbers in correct position as locked")
		return
	}

	pi := partition(columnCh, nums, low, high, locked)

	quicksortHelper(columnCh, nums, low, pi-1, locked)
	quicksortHelper(columnCh, nums, pi+1, high, locked)
}

func Quicksort(columnCh chan<- component.Column, nums []int) {
	defer close(columnCh)
	quicksortHelper(columnCh, nums, 0, len(nums)-1, map[int]string{})
}
