package tasksData

import (
	"fmt"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/models"
)

func SummingTasks() *[]tasksModels.Task {
	return &[]tasksModels.Task{
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{10, 49497}},
			Generate: func(args... int) []int {
				n := args[0]

				return []int{n}
			},
			Math: func (args ...int) string {
				n := args[0]

				return fmt.Sprintf("6 \\cdot \\frac{1 \\cdot 2 + 2 \\cdot 3 + 3 \\cdot 4 + \\ldots +  %d \\cdot %d}{1 + 2 + 3 + \\ldots + %d}",
					n, n + 1, n)
			},
			Answer: func (args ...int) int {
				n := args[0]

				return 4 * (n + 2)
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{6, 121}},
			Generate: func(args... int) []int {
				n := args[0]

				return []int{n}
			},
			Math: func (args ...int) string {
				n := args[0]

				return fmt.Sprintf("\\frac{1}{%d!} \\cdot \\Bigg ( 1 - \\Bigg ( \\frac{1}{2!} + " +
					"\\frac{2}{3!} + \\frac{3}{4!} + \\ldots + \\frac{%d}{%d!} \\Bigg ) \\Bigg )^{-1}",
					n - 3, n - 1, n)
			},
			Answer: func (args ...int) int {
				n := args[0]

				return n * (n - 1) * (n - 2)
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{4, 976}},
			Generate: func(args... int) []int {
				k := args[0]

				return []int{k}
			},
			Math: func (args ...int) string {
				k := args[0]
				n := 2 * k + 1
				double := 2 * n

				return fmt.Sprintf("\\sin^2 \\frac{\\pi}{%d} + \\sin^2 \\frac{2 \\pi}{%d} + " +
					"\\sin^2 \\frac{3 \\pi}{%d} + \\ldots + \\sin^2 \\frac{%d \\pi}{%d}",
					double, double, double, n, double)
			},
			Answer: func (args ...int) int {
				k := args[0]

				return k + 1
			},
		},
	}
}
