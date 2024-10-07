package algorithm

import (
	"fmt"
	"maps"
	"math"
	"slices"
	term_utils "tui/internal/term-utils"
)

func JumpSearch(columnCh chan<- ColumnGraphData, nums []int, target int) {
	defer close(columnCh)

	sortedNums := slices.Clone(nums)
	slices.Sort(sortedNums)

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		map[int]string{},
		"Array needs to be sorted for jump search")

	colors := map[int]string{}

	step := int(math.Sqrt(float64(len(nums))))

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		map[int]string{},
		fmt.Sprintf("Jump amount is %s",
			term_utils.Colorize(step, term_utils.Blue)))

	i := step

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		map[int]string{i: term_utils.Green},
		fmt.Sprintf("Jump %s to %s",
			term_utils.Colorize(step, term_utils.Blue),
			term_utils.Colorize(sortedNums[i], term_utils.Green)))

	for ; i < len(nums); i += step {

		colors[i] = term_utils.Green
		columnCh <- NewColumnGraphData(
			slices.Clone(sortedNums),
			maps.Clone(colors),
			fmt.Sprintf("Check if %s > %s",
				term_utils.Colorize(sortedNums[i], term_utils.Green),
				term_utils.Colorize(target, term_utils.Orange)))

		if sortedNums[i] > target {

			columnCh <- NewColumnGraphData(
				slices.Clone(sortedNums),
				maps.Clone(colors),
				fmt.Sprintf("%s > %s -> Jump back and walk from last good position",
					term_utils.Colorize(sortedNums[i], term_utils.Green),
					term_utils.Colorize(target, term_utils.Orange)))

			break
		}

		for j := i - step; j <= i; j++ {
			colors[j] = term_utils.Lightgray
		}

		columnCh <- NewColumnGraphData(
			slices.Clone(sortedNums),
			maps.Clone(colors),
			fmt.Sprintf("%s <= %s -> Jump %s again",
				term_utils.Colorize(sortedNums[i], term_utils.Green),
				term_utils.Colorize(target, term_utils.Orange),
				term_utils.Colorize(step, term_utils.Blue)))

	}

	i -= step

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		maps.Clone(colors),
		fmt.Sprintf("Now walk linearly from here %s",
			term_utils.Colorize(sortedNums[i], term_utils.Green)))

	for j := 0; j < step && i < len(nums); j++ {

		colors[i] = term_utils.Green
		columnCh <- NewColumnGraphData(
			slices.Clone(sortedNums),
			maps.Clone(colors),
			fmt.Sprintf("Is %s == %s",
				term_utils.Colorize(sortedNums[i], term_utils.Green),
				term_utils.Colorize(target, term_utils.Orange)))

		if sortedNums[i] == target {

			columnCh <- NewColumnGraphData(
				slices.Clone(sortedNums),
				maps.Clone(colors),
				fmt.Sprintf("Found %s at index %s",
					term_utils.Colorize(target, term_utils.Orange),
					term_utils.Colorize(i, term_utils.Green)))
			return

		}
		colors[i] = term_utils.Lightgray
		i++
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		maps.Clone(colors),
		fmt.Sprintf("%s is not in array!",
			term_utils.Colorize(target, term_utils.Orange)))

}
