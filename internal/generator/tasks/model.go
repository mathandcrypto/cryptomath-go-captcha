package tasks

type Task struct {
	Difficulty	int16
	GenerateParams [][2]int
	Generate	func(...int) []int
	Math	func(...int) string
	Answer	func(...int) int
}
