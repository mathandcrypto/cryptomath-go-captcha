package tasksModels

type Task struct {
	Difficulty	int16
	GenerateParams [][2]int
	Generate	func(...int) []int
	Math	func(...int) string
	Answer	func(...int) int
}

type GeneratedTask struct {
	Difficulty	int16
	Math	string
	Answer int
}
