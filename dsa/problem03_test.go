package dsa

import (
	"reflect"
	"testing"
)

func TestSumOfDistancesInTree(t *testing.T) {
	n := 6
	edges := [][]int{{0, 1}, {0, 2}, {2, 3}, {2, 4}, {2, 5}}

	expected := []int{8, 12, 6, 10, 10, 10}
	actual := sumOfDistancesInTree(n, edges)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Actual %v", expected, actual)
	}

	reset()
	n = 1
	edges = [][]int{}

	expected = []int{0}
	actual = sumOfDistancesInTree(n, edges)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Actual %v", expected, actual)
	}

	reset()
	n = 2
	edges = [][]int{{1, 0}}

	expected = []int{1, 1}
	actual = sumOfDistancesInTree(n, edges)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Actual %v", expected, actual)
	}
}
