package algorithm

import (
	"fmt"
	"maps"
	"slices"
	term_utils "tui/internal/term-utils"
)

func Mergesort(columnCh chan<- ColumnGraphData, nums []int) {
	defer close(columnCh)
	colors := make(map[int]string)
	mergesortHelper(columnCh, colors, nums, 0, len(nums)-1)
}

func mergesortHelper(
	columnCh chan<- ColumnGraphData,
	colors map[int]string,
	nums []int,
	left, right int,
) {
	if left < right {
		for i := 0; i < len(nums); i++ {
			if i < left || i > right {
				colors[i] = term_utils.Lightgray
			} else {
				colors[i] = term_utils.White
			}
		}

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			"Divide the array into two equal sized arrays",
		)

		mid := (right + left) / 2
		for i := left; i <= mid; i++ {
			colors[i] = term_utils.Green
		}
		for i := mid + 1; i <= right; i++ {
			colors[i] = term_utils.Blue
		}

		columnCh <- NewColumnGraphData(
			slices.Clone(nums),
			maps.Clone(colors),
			fmt.Sprintf("Array divided into %s and %s",
				term_utils.Colorize("left side", term_utils.Green),
				term_utils.Colorize("right side", term_utils.Blue),
			),
		)

		mergesortHelper(columnCh, colors, nums, left, mid)
		mergesortHelper(columnCh, colors, nums, mid+1, right)
		merge(columnCh, colors, nums, left, mid, right)
	}
}

func merge(columnCh chan<- ColumnGraphData, colors map[int]string, nums []int, left, mid, right int) {
	for i := left; i <= mid; i++ {
		colors[i] = term_utils.Green
	}

	for i := mid + 1; i <= right; i++ {
		colors[i] = term_utils.Blue
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		maps.Clone(colors),
		fmt.Sprintf("Merge these two arrays"),
	)

	l := make([]int, mid-left+1)
	r := make([]int, right-mid)

	//TODO: Add a way to show this in the visualization
	for i, n := range nums[left : mid+1] {
		l[i] = n
	}

	for i, n := range nums[mid+1 : right+1] {
		r[i] = n
	}

	i, j, k := 0, 0, left
	for i < len(l) && j < len(r) {
		if l[i] <= r[j] {
			nums[k] = l[i]
			i++
		} else {
			nums[k] = r[j]
			j++
		}
		k++
	}

	for i < len(l) {
		nums[k] = l[i]
		i++
		k++
	}

	for j < len(r) {
		nums[k] = r[j]
		j++
		k++
	}

	for i := mid + 1; i <= right; i++ {
		colors[i] = term_utils.Green
	}

	columnCh <- NewColumnGraphData(
		slices.Clone(nums),
		maps.Clone(colors),
		fmt.Sprintf("Array merged"),
	)
}
