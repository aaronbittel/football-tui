package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Insertionsort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)

	colors := make(map[int]string)
	grays := make(map[int]string)
	for i := 1; i < len(nums); i++ {
		grays[i] = term_utils.Lightgray
	}

	for i := 1; i < len(nums); i++ {
		colors[i] = term_utils.Green
		delete(grays, i)

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf("Current %s", term_utils.Colorize(nums[i], term_utils.Green)))

		grays[i] = term_utils.Green
		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(grays),
			fmt.Sprintf(
				"Insert %s into current array",
				term_utils.Colorize(nums[i], term_utils.Green)))
		j := i
		for j = i; j > 0; j-- {
			if nums[j] >= nums[j-1] {
				break
			}
			nums[j], nums[j-1] = nums[j-1], nums[j]
		}

		delete(grays, i)
		grays[j] = term_utils.Green

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(grays),
			fmt.Sprintf(
				"%s Inserted",
				term_utils.Colorize(nums[j], term_utils.Green)))

		delete(grays, j)
		delete(colors, i)
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		map[int]string{},
		term_utils.Colorize("Insertionsort complete!", term_utils.Green))
}
