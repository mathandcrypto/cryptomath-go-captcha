package tasksData

import (
	"fmt"

	tasksModels "github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/models"
)

func LimitTasks() *[]tasksModels.Task {
	return &[]tasksModels.Task{
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{2, 1024}, {1, 4}},
			Generate: func(args... int) []int {
				m := args[0]
				n := args[1]
				k := m * n

				return []int{m, k}
			},
			Math: func (args ...int) string {
				m := args[0]
				n := args[1]

				return fmt.Sprintf("\\lim_{x \\to 1} \\frac{\\sqrt[%d]{x} - 1}{\\sqrt[%d]{x} - 1}", m, n)
			},
			Answer: func (args ...int) int {
				m := args[0]
				n := args[1]

				return n / m
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{1, 99}, {1, 73}},
			Generate: func(args... int) []int {
				b := args[0]
				a := args[1]
				c := b + a

				return []int{c, b}
			},
			Math: func (args ...int) string {
				a := args[0]
				b := args[1]
				m := 3 * a * b

				return fmt.Sprintf("%d \\cdot \\lim_{x \\to 0} \\frac{1}{x \\sqrt{x}} \\Bigg ( \\sqrt{%d} " +
					"\\arctan \\sqrt{\\frac{x}{%d}} - \\sqrt{%d} \\arctan \\sqrt{\\frac{x}{%d}} \\; \\Bigg )",
					m, a, a, b, b)
			},
			Answer: func (args ...int) int {
				a := args[0]
				b := args[1]

				return a - b
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{1, 12}, {1, 9}},
			Generate: func(args... int) []int {
				b := args[0]
				a := args[1]
				c := b + a

				return []int{a, c}
			},
			Math: func (args ...int) string {
				b := args[0]
				a := args[1]

				return fmt.Sprintf("\\ln \\lim_{x \\to 0} \\Bigg ( \\frac{1 + \\sin x \\cos %d x}{1 + \\sin x \\cos %d x} \\Bigg )^{\\cot^3 x}", a, b)
			},
			Answer: func (args ...int) int {
				b := args[0]
				a := args[1]

				return b * b - a * a
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{3, 41}, {1, 23}, {7, 37}},
			Generate: func(args... int) []int {
				a := args[0]
				b := args[1]
				c := args[2]

				return []int{a, b, c}
			},
			Math: func (args ...int) string {
				a := args[0]
				b := args[1]
				c := args[2]

				return fmt.Sprintf("\\sqrt[3]{\\lim_{x \\to 0} \\Bigg ( \\frac{%d^x + %d^x + %d^x}{3} \\Bigg )^{\\frac{1}{x}}}", a, b, c)
			},
			Answer: func (args ...int) int {
				a := args[0]
				b := args[1]
				c := args[2]

				return a * b * c
			},
		},
	}
}
