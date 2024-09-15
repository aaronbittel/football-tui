package algorithms

import (
	"slices"
)

func BubbleSort(numsCh chan<- []int, nums []int) {
	defer close(numsCh)
	for i := len(nums) - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			numsCh <- slices.Clone(nums)
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
}
