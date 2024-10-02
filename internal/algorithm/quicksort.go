package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func partition(
	columnCh chan<- ColumnGraphData,
	nums []int,
	low, high int,
	colors map[int]string,
) int {
	pivot := nums[high]
	i := low - 1

	for i := 0; i < len(nums); i++ {
		if i >= low && i <= high {
			continue
		}

		if colors[i] == term_utils.Orange {
			continue
		}
		colors[i] = term_utils.Lightgray
	}

	colors[high] = term_utils.Green

	for j := low; j < high; j++ {
		colors[j] = term_utils.Blue
		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Comparing %s to Pivot %s",
				term_utils.Colorize(fmt.Sprintf("%d", nums[j]), term_utils.Blue),
				term_utils.Colorize(fmt.Sprintf("%d", pivot), term_utils.Green),
			),
		)
		if nums[j] < pivot {
			i++
			delete(colors, j)
			nums[i], nums[j] = nums[j], nums[i]
			colors[i] = term_utils.Blue
			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf(
					"Swap %s and %s because %s is smaller than Pivot %s",
					term_utils.Colorize(fmt.Sprintf("%d", nums[i]), term_utils.Blue),
					term_utils.Colorize(fmt.Sprintf("%d", nums[j]), term_utils.White),
					term_utils.Colorize(fmt.Sprintf("%d", nums[i]), term_utils.Blue),
					term_utils.Colorize(fmt.Sprintf("%d", pivot), term_utils.Green),
				),
			)
			delete(colors, i)
		}
		delete(colors, j)
	}

	delete(colors, high)
	nums[i+1], nums[high] = nums[high], nums[i+1]
	colors[i+1] = term_utils.Green
	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		maps.Clone(colors),
		fmt.Sprintf("Swap Pivot %s to correct position", term_utils.Colorize(
			fmt.Sprintf("%d", nums[i+1]),
			term_utils.Green,
		)),
	)

	for k, v := range colors {
		if v != term_utils.Orange {
			delete(colors, k)
		}
	}

	return i + 1
}

func quicksortHelper(
	columnCh chan<- ColumnGraphData,
	nums []int,
	low, high int,
	colors map[int]string,
) {
	if low >= high {
		//FIX: Implement Correct highlighting of Locked values

		// term_utils.Debug("lowIdx", low, "highIdx", high, "lowVal", nums[low], "highVal", nums[high])
		// colors[high] = term_utils.Orange
		// colors[low] = term_utils.Orange
		//
		// var msg string
		// if low != high {
		// 	msg = fmt.Sprintf("Mark numbers %s, %s in correct position as locked",
		// 		term_utils.Colorize(fmt.Sprintf("%d", nums[high]), term_utils.Orange),
		// 		term_utils.Colorize(fmt.Sprintf("%d", nums[low]), term_utils.Orange),
		// 	)
		// } else {
		// 	msg = fmt.Sprintf("Mark numbers %s in correct position as locked",
		// 		term_utils.Colorize(fmt.Sprintf("%d", nums[low]), term_utils.Orange),
		// 	)
		// }
		//
		// columnCh <- component.NewColumn(
		// 	slices.Clone(nums),
		// 	maps.Clone(colors),
		// 	msg,
		// )
		return
	}

	pi := partition(columnCh, nums, low, high, colors)

	quicksortHelper(columnCh, nums, low, pi-1, colors)
	quicksortHelper(columnCh, nums, pi+1, high, colors)
}

func Quicksort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)
	quicksortHelper(columnCh, nums, 0, len(nums)-1, map[int]string{})
}
