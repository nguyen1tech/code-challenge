package dsa

// Gray Code
func grayCode(n int) []int {
	res := []int{0}
	if n == 0 {
		return res
	}
	res = append(res, 1)
	curr := 1
	for i := 2; i <= n; i++ {
		curr *= 2
		for j := len(res) - 1; j >= 0; j-- {
			res = append(res, curr+res[j])
		}
	}
	return res
}
