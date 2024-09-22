package component

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Selectionsort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)

	for i := 0; i < len(nums); i++ {
		minVal := nums[i]
		minIdx := i
		for j := 1 + i; j < len(nums); j++ {
			colors[minIdx] = term_utils.Green
			colors[j] = term_utils.Blue
			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf("Comparing current min value %s to %s",
					term_utils.Colorize(minVal, term_utils.Green),
					term_utils.Colorize(nums[j], term_utils.Blue)))
			if nums[j] < minVal {
				oldMinVal := minVal
				oldMinIdx := minIdx
				colors[oldMinIdx] = term_utils.Blue
				minVal = nums[j]
				minIdx = j
				colors[minIdx] = term_utils.Green
				columnCh <- NewColumnGraphData(
					slices.Clone(nums),
					maps.Clone(colors),
					fmt.Sprintf("New min value is %s, because %s < %s",
						term_utils.Colorize(minVal, term_utils.Green),
						term_utils.Colorize(minVal, term_utils.Green),
						term_utils.Colorize(oldMinVal, term_utils.Blue)))
				delete(colors, oldMinIdx)
				delete(colors, minIdx)
			}
			delete(colors, j)
			delete(colors, minIdx)
		}
		nums[i], nums[minIdx] = nums[minIdx], nums[i]
		colors[i] = term_utils.Orange

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"%s is now locked",
				term_utils.Colorize(nums[i], term_utils.Orange)))
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		maps.Clone(colors),
		"Selection sort completed")
}
