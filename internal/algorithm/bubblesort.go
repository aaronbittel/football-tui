package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Bubblesort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)

	for i := len(nums) - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			colors[j] = term_utils.Green
			colors[j+1] = term_utils.Blue

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf("Comparing %s to %s",
					term_utils.Colorize(nums[j], term_utils.Green),
					term_utils.Colorize(nums[j+1], term_utils.Blue)))

			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
				colors[j] = term_utils.Blue
				colors[j+1] = term_utils.Green
				columnCh <- NewColumnGraphData(
					slices.Clone(nums),
					maps.Clone(colors),
					fmt.Sprintf("Swapping %s and %s because %s > %s",
						term_utils.Colorize(nums[j+1], term_utils.Green),
						term_utils.Colorize(nums[j], term_utils.Blue),
						term_utils.Colorize(nums[j+1], term_utils.Green),
						term_utils.Colorize(nums[j], term_utils.Blue)),
				)
			}
			delete(colors, j)
			delete(colors, j+1)
		}
		colors[i] = term_utils.Orange
		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf("%s is now locked",
				term_utils.Colorize(nums[i], term_utils.Orange)),
		)
	}
}
