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
	colors map[int]string,
) int {
	pivot := nums[high]
	i := low - 1

	for i := 0; i < len(nums); i++ {
		if i >= low && i <= high {
			continue
		}

		if colors[i] == utils.Orange {
			continue
		}
		colors[i] = utils.Lightgray
	}

	colors[high] = utils.Green

	for j := low; j < high; j++ {
		colors[j] = utils.Blue
		columnCh <- component.NewColumn(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Comparing %s to Pivot %s",
				utils.Colorize(fmt.Sprintf("%d", nums[j]), utils.Blue),
				utils.Colorize(fmt.Sprintf("%d", pivot), utils.Green),
			),
		)
		if nums[j] < pivot {
			i++
			delete(colors, j)
			nums[i], nums[j] = nums[j], nums[i]
			colors[i] = utils.Blue
			columnCh <- component.NewColumn(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf(
					"Swap %s and %s because %s is smaller than Pivot %s",
					utils.Colorize(fmt.Sprintf("%d", nums[i]), utils.Blue),
					utils.Colorize(fmt.Sprintf("%d", nums[j]), utils.White),
					utils.Colorize(fmt.Sprintf("%d", nums[i]), utils.Blue),
					utils.Colorize(fmt.Sprintf("%d", pivot), utils.Green),
				),
			)
			delete(colors, i)
		}
		delete(colors, j)
	}

	delete(colors, high)
	nums[i+1], nums[high] = nums[high], nums[i+1]
	colors[i+1] = utils.Green
	columnCh <- component.NewColumn(
		slices.Clone(nums),
		maps.Clone(colors),
		fmt.Sprintf("Swap Pivot %s to correct position", utils.Colorize(
			fmt.Sprintf("%d", nums[i+1]),
			utils.Green,
		)),
	)

	for k, v := range colors {
		if v != utils.Orange {
			delete(colors, k)
		}
	}

	return i + 1
}

func quicksortHelper(
	columnCh chan<- component.Column,
	nums []int,
	low, high int,
	colors map[int]string,
) {
	if low >= high {
		//FIX: Implement Correct highlighting of Locked values

		// utils.Debug("lowIdx", low, "highIdx", high, "lowVal", nums[low], "highVal", nums[high])
		// colors[high] = utils.Orange
		// colors[low] = utils.Orange
		//
		// var msg string
		// if low != high {
		// 	msg = fmt.Sprintf("Mark numbers %s, %s in correct position as locked",
		// 		utils.Colorize(fmt.Sprintf("%d", nums[high]), utils.Orange),
		// 		utils.Colorize(fmt.Sprintf("%d", nums[low]), utils.Orange),
		// 	)
		// } else {
		// 	msg = fmt.Sprintf("Mark numbers %s in correct position as locked",
		// 		utils.Colorize(fmt.Sprintf("%d", nums[low]), utils.Orange),
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

func Quicksort(columnCh chan<- component.Column, nums []int) {
	defer close(columnCh)
	quicksortHelper(columnCh, nums, 0, len(nums)-1, map[int]string{})
}
