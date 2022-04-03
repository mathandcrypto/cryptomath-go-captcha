package tasksHelpers

func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}

	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}

	return result
}
