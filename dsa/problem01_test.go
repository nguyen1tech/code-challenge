package dsa

import (
	"reflect"
	"testing"
)

func TestGrayCode(t *testing.T) {
	expected := []int{0, 1, 3, 2}
	actual := grayCode(2)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	expected = []int{0, 1}
	actual = grayCode(1)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
