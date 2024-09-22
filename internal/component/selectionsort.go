package component

import (
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Selectionsort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)

	for i := 0; i < len(nums); i++ {
		minVal := nums[i]
		minIndex := i
		for j := 1 + i; j < len(nums); j++ {
			colors[minIndex] = term_utils.Green
			colors[j] = term_utils.Blue
			columnCh <- NewColumnGraphData(slices.Clone(nums), maps.Clone(colors), "")
			if nums[j] < minVal {
				delete(colors, minIndex)
				minVal = nums[j]
				minIndex = j
				colors[minIndex] = term_utils.Green
			}
			delete(colors, j)
			delete(colors, minIndex)
		}
		nums[i], nums[minIndex] = nums[minIndex], nums[i]
		colors[i] = term_utils.Orange
	}

	columnCh <- NewColumnGraphData(slices.Clone(nums), maps.Clone(colors), "")
}
