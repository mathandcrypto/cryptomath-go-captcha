package tasks

import (
	"fmt"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/generator/helpers/math"
)

func SeriesTasks() *[]Task {
	return &[]Task{
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{1, 10}, {2, 7}, {1, 6}},
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

				return fmt.Sprintf("\\sum\\limits_{k = 0}^{%d} {%d \\choose k} + \\sum\\limits_{k = 1}^{%d} " +
					"k {%d \\choose k} + \\sum\\limits_{k = 0}^{%d} {%d \\choose k}^2",
					a, a, b, b, c, c)
			},
			Answer: func (args ...int) int {
				a := args[0]
				b := args[1]
				c := args[2]

				return mathHelpers.IntPow(2, a) + b * mathHelpers.IntPow(2, b - 1) + mathHelpers.Binomial(2 * c, c)
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{3, 9999}},
			Generate: func(args... int) []int {
				n := args[0]

				return []int{n}
			},
			Math: func (args ...int) string {
				n := args[0]

				return fmt.Sprintf("\\sum\\limits_{k = 0}^{%d} (-1)^{%d - k} \\, 2^{2 k} {%d + k + 1 \\choose 2k + 1}", n, n, n)
			},
			Answer: func (args ...int) int {
				n := args[0]

				return n + 1
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{1, 11}, {2, 37}},
			Generate: func(args... int) []int {
				a := args[0]
				b := args[1]
				k := 2 * a
				n := k + b

				return []int{n, k}
			},
			Math: func (args ...int) string {
				n := args[0]
				k := args[1]

				return fmt.Sprintf("\\sum\\limits_{i = 0}^{%d} (-1)^i {%d \\choose %d - i} {%d \\choose i}", k, n, k, n)
			},
			Answer: func (args ...int) int {
				n := args[0]
				k := args[1]

				if k % 2 == 0 {
					half := k / 2

					return mathHelpers.IntPow(-1, half) * mathHelpers.Binomial(n, half)
				}

				return 0
			},
		},
		{
			Difficulty: 1,
			GenerateParams: [][2]int{{1, 5}, {2, 13}},
			Generate: func(args... int) []int {
				k := args[0]
				l := args[1]

				return []int{k, l}
			},
			Math: func (args ...int) string {
				k := args[0]
				l := args[1]
				a := mathHelpers.IntPow(2, k)
				b := a + mathHelpers.IntPow(2, l)

				return fmt.Sprintf("\\log_2 \\Bigg( \\sum\\limits_{n=1}^{\\infty} \\Big ( \\frac{%d}{%d} \\Big )^2 \\Bigg )", a, b)
			},
			Answer: func (args ...int) int {
				k := args[0]
				l := args[1]

				return k - l
			},
		},
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{2, 100}, {1, 3}},
			Generate: func(args... int) []int {
				q := args[0]
				p := args[1]

				return []int{q, p}
			},
			Math: func (args ...int) string {
				q := args[0]
				p := args[1]
				a := mathHelpers.IntPow(q - 1, 2)
				b := q * q

				return fmt.Sprintf("\\frac{%d}{%d} \\cdot \\sum\\limits_{k=1}^{\\infty} \\frac{k}{%d^{k - %d}}", a, b, q, p)
			},
			Answer: func (args ...int) int {
				q := args[0]
				p := args[1]

				return mathHelpers.IntPow(q, p - 1)
			},
		},
	}
}
