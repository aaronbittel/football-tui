package algorithms

import (
	"maps"
	"slices"
	component "tui/internal/component"
	utils "tui/internal/term-utils"
)

func Bubblesort(columnCh chan<- component.ColumnGraphData, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)

	for i := len(nums) - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			colors[j] = utils.Green
			colors[j+1] = utils.Blue
			columnCh <- component.NewColumnGraphData(slices.Clone(nums), maps.Clone(colors), "")
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
				colors[j] = utils.Blue
				colors[j+1] = utils.Green
				columnCh <- component.NewColumnGraphData(slices.Clone(nums), maps.Clone(colors), "")
			}
			delete(colors, j)
			delete(colors, j+1)
		}
		colors[i] = utils.Orange
	}
	columnCh <- component.NewColumnGraphData(slices.Clone(nums), maps.Clone(colors), "")
}
