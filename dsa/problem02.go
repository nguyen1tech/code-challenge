package dsa

// Maximum Length of Repeated Subarray
func findLength(nums1 []int, nums2 []int) int {
	res := 0
	for i := 0; i < len(nums1)+len(nums2)-1; i++ {
		start1 := max(0, len(nums1)-1-i)
		start2 := max(0, i-(len(nums1)-1))

		count := 0
		idx1, idx2 := start1, start2
		for idx1 < len(nums1) && idx2 < len(nums2) {
			if nums1[idx1] == nums2[idx2] {
				count++
			} else {
				count = 0
			}
			res = max(res, count)

			idx1 += 1
			idx2 += 1
		}
	}
	return res
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
