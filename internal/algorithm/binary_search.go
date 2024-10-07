package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func BinarySearch(columnCh chan<- ColumnGraphData, nums []int, target int) {
	defer close(columnCh)
	sortedNums := slices.Clone(nums)
	slices.Sort(sortedNums)

	columnCh <- NewColumnGraphData(
		slices.Clone(sortedNums),
		map[int]string{},
		"Array needs to be sorted for binary search")

	binarySearchHelper(columnCh, sortedNums, target)
}

func binarySearchHelper(columnCh chan<- ColumnGraphData, nums []int, target int) {
	var (
		mid    int
		low    = 0
		high   = len(nums) - 1
		colors = map[int]string{}
	)

	for low <= high {

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"This is the current search space, searching for %s",
				term_utils.Colorize(target, term_utils.Orange)))

		mid = (high-low)/2 + low

		colors[mid] = term_utils.Green

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf(
				"Comparing %s to target %s",
				term_utils.Colorize(nums[mid], term_utils.Green),
				term_utils.Colorize(target, term_utils.Orange)))

		if nums[mid] < target {
			for i := low; i < mid+1; i++ {
				colors[i] = term_utils.Lightgray
			}
			low = mid + 1
			if low >= len(nums) {
				continue
			}

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf(
					"%s < %s -> limit search space from %s to %s",
					term_utils.Colorize(nums[mid], term_utils.Green),
					term_utils.Colorize(target, term_utils.Orange),
					term_utils.Colorize(nums[low], term_utils.Blue),
					term_utils.Colorize(nums[high], term_utils.Blue)))

		} else if nums[mid] > target {
			for i := high; i > mid-1; i-- {
				colors[i] = term_utils.Lightgray
			}

			high = mid - 1
			if high < 0 {
				continue
			}

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				maps.Clone(colors),
				fmt.Sprintf(
					"%s > %s -> limit search space from %s to %s",
					term_utils.Colorize(nums[mid], term_utils.Green),
					term_utils.Colorize(target, term_utils.Orange),
					term_utils.Colorize(nums[low], term_utils.Blue),
					term_utils.Colorize(nums[high], term_utils.Blue)))

		} else {

			columnCh <- NewColumnGraphData(
				slices.Clone(nums),
				// maps.Clone(colors),
				map[int]string{mid: term_utils.Orange},
				fmt.Sprintf(
					"Found target %s!",
					term_utils.Colorize(nums[mid], term_utils.Orange)))

		}
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		maps.Clone(colors),
		fmt.Sprintf(
			"Target %s is not in the array!",
			term_utils.Colorize(target, term_utils.Orange)))
}
