package dsa

import "testing"

func TestFindLength(t *testing.T) {
	nums1 := []int{1, 2, 3, 2, 1}
	nums2 := []int{3, 2, 1, 4, 7}

	expected := 3
	actual := findLength(nums1, nums2)
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	nums1 = []int{0, 0, 0, 0, 0}
	nums2 = []int{0, 0, 0, 0, 0}
	expected = 5
	actual = findLength(nums1, nums2)
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
