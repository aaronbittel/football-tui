package algorithms

import (
	"maps"
	"slices"
	component "tui/internal/component"
	utils "tui/internal/term-utils"
)

func Selectionsort(columnCh chan<- component.Column, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)

	for i := 0; i < len(nums); i++ {
		minVal := nums[i]
		minIndex := i
		for j := 1 + i; j < len(nums); j++ {
			colors[minIndex] = utils.Green
			colors[j] = utils.Blue
			columnCh <- component.NewColumn(slices.Clone(nums), maps.Clone(colors), "")
			if nums[j] < minVal {
				delete(colors, minIndex)
				minVal = nums[j]
				minIndex = j
				colors[minIndex] = utils.Green
			}
			delete(colors, j)
			delete(colors, minIndex)
		}
		nums[i], nums[minIndex] = nums[minIndex], nums[i]
		colors[i] = utils.Orange
	}

	columnCh <- component.NewColumn(slices.Clone(nums), maps.Clone(colors), "")
}
