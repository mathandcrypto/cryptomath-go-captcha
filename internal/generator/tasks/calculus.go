package tasks

import "fmt"

func CalculusTasks() *[]Task {
	return &[]Task{
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{2, 783}, {3, 97}},
			Generate: func(args... int) []int {
				a := args[0]
				b := args[1]
				c := b + a

				return []int{a, c}
			},
			Math: func (args ...int) string {
				a := args[0]
				b := args[1]

				return fmt.Sprintf("%d \\cdot \\Bigg ( \\frac{\\sqrt{%d}}{\\sqrt{%d} + " +
					"\\sqrt{%d}} + \\frac{\\sqrt{%d}}{\\sqrt{%d} - \\sqrt{%d}} \\Bigg )",
					a - b, a, a, b, b, a, b,
				)
			},
			Answer: func (args ...int) int {
				a := args[0]
				b := args[1]

				return a + b
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{2, 1000}},
			Generate: func(args... int) []int {
				n := args[0]

				return []int{n}
			},
			Math: func (args ...int) string {
				n := args[0]

				return fmt.Sprintf("\\big ( (%d! + %d!) \\cdot %d) / (%d! - %d!) \\big )",
					n + 1, n, n, n + 1, n,
				)
			},
			Answer: func (args ...int) int {
				n := args[0]

				return n + 2
			},
		},
	}
}