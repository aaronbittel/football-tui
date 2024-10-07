package algorithm

import (
	"fmt"
	"slices"
	term_utils "tui/internal/term-utils"
)

// TODO: Dont really need to send nums in the columnCh because they never change
// in a serach algorithm
func LinearSearch(columnCh chan<- ColumnGraphData, nums []int, target int) {
	defer close(columnCh)

	for i, n := range nums {

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			map[int]string{i: term_utils.Green},
			fmt.Sprintf(
				"Is %s our target (%s)?",
				term_utils.Colorize(n, term_utils.Green),
				term_utils.Colorize(target, term_utils.Orange)))

		if n == target {

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				map[int]string{i: term_utils.Green},
				fmt.Sprintf(
					"Found target %s at index %s!",
					term_utils.Colorize(target, term_utils.Orange),
					term_utils.Colorize(i, term_utils.Green)))
			return
		}
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		map[int]string{},
		fmt.Sprintf(
			"Target %s is not in the list!",
			term_utils.Colorize(target, term_utils.Orange)))

}
