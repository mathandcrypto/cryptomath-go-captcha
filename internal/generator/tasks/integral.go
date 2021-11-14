package tasks

import (
	"fmt"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/generator/helpers/math"
)

func IntegralTasks() *[]Task {
	return &[]Task{
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{1, 15}},
			Generate: func(args... int) []int {
				m := args[0]

				return []int{m}
			},
			Math: func (args ...int) string {
				m := args[0]
				p := m * 3

				return fmt.Sprintf("\\lim_{\\alpha \\to 0} \\int\\limits_{0}^{%d} x^2 \\cos \\alpha x \\; d{x}", p)
			},
			Answer: func (args ...int) int {
				m := args[0]

				return 9 * mathHelpers.IntPow(m, 3)
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{2, 12}, {3, 23}},
			Generate: func(args... int) []int {
				m := args[0]
				n := args[1]

				return []int{m, n}
			},
			Math: func (args ...int) string {
				m := args[0]
				n := args[1]
				a := 2 * m
				b := 2 * n

				return fmt.Sprintf("exp \\Bigg ( {\\frac{1}{\\pi} \\int\\limits_{0}^{\\pi/2} " +
					"\\ln (%d^2 \\sin^2 x \\; + \\; %d^2 \\cos^2 x) \\; {d}x} \\Bigg )",
					a, b)
			},
			Answer: func (args ...int) int {
				m := args[0]
				n := args[1]

				return m + n
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{1, 11}, {2, 83}},
			Generate: func(args... int) []int {
				a := args[0]
				b := args[1]
				c := a * b

				return []int{a, c}
			},
			Math: func (args ...int) string {
				a := args[0]
				b := args[1]

				return fmt.Sprintf("\\exp \\Bigg ( \\int\\limits_{0}^{\\infty} \\frac{\\exp (-%d x) - \\exp (-%d x)}{x} \\; {d}x \\Bigg )", a, b)
			},
			Answer: func (args ...int) int {
				a := args[0]
				b := args[1]

				return b / a
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{1, 54321}},
			Generate: func(args... int) []int {
				a := args[0]

				return []int{a}
			},
			Math: func (args ...int) string {
				a := args[0]

				return fmt.Sprintf("\\ln \\Bigg ( \\frac{2}{\\pi} \\int\\limits_{0}^{+\\infty} \\frac{\\cos %dx}{1 + x^2} \\; dx \\Bigg )", a)
			},
			Answer: func (args ...int) int {
				a := args[0]

				return (-1) * a
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{1, 17}},
			Generate: func(args... int) []int {
				a := args[0]

				return []int{a}
			},
			Math: func (args ...int) string {
				a := args[0]

				return fmt.Sprintf("\\frac{32}{\\pi} \\int\\limits_{0}^{%d} x^2 \\sqrt{%d - x^2} \\; dx", a, a * a)
			},
			Answer: func (args ...int) int {
				a := args[0]

				return 2 * mathHelpers.IntPow(a, 4)
			},
		},
	}
}
