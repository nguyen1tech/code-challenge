package dsa

var (
	tree  []map[int]bool
	res   []int
	count []int
	N     int
)

// Sum of Distances in Tree
func sumOfDistancesInTree(n int, edges [][]int) []int {
	N = n
	tree = make([]map[int]bool, N)
	res = make([]int, N)
	count = make([]int, N)
	for i := 0; i < N; i++ {
		tree[i] = make(map[int]bool)
	}
	for _, e := range edges {
		tree[e[0]][e[1]] = true
		tree[e[1]][e[0]] = true
	}

	dfs(0, -1)
	dfs2(0, -1)
	return res
}

func dfs(root int, pre int) {
	for i := range tree[root] {
		if i != pre {
			dfs(i, root)
			count[root] += count[i]
			res[root] += res[i] + count[i]
		}
	}
	count[root]++
}

func dfs2(root int, pre int) {
	for i := range tree[root] {
		if i != pre {
			res[i] = res[root] - count[i] + N - count[i]
			dfs2(i, root)
		}
	}
}

func reset() {
	tree, res, count = nil, nil, nil
	N = 0
}
