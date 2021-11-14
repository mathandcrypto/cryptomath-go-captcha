package mathHelpers

func Factorial(n int) int {
	val := 1
	for i := 2; i <= n; i++ {
		val *= i
	}

	return val
}

func Binomial(n, k int) int {
	numerator := Factorial(n)
	denominator := Factorial(n - k) * Factorial(k)

	return numerator / denominator
}
