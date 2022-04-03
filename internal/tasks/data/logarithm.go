package tasksData

import (
	"fmt"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/helpers"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/models"
)

func LogarithmTasks() *[]tasksModels.Task {
	return &[]tasksModels.Task{
		{
			Difficulty: 2,
			GenerateParams: [][2]int{{3, 23}, {2, 100}},
			Generate: func(args... int) []int {
				p := args[0]
				m := args[1]

				return []int{p, m}
			},
			Math: func (args ...int) string {
				p := args[0]
				m := args[1]
				k := p * p
				n := m * k
				t := tasksHelpers.IntPow(p * m, 2)

				return fmt.Sprintf("%d \\cdot \\Bigg( \\frac{\\log_{%d} %d}{\\log_{%d} %d} - \\frac{\\log_{%d} %d}{\\log_{%d} %d} \\Bigg )",
					m + n, m, n, n, m, m ,t, k, m,
				)
			},
			Answer: func (args ...int) int {
				p := args[0]
				m := args[1]

				return m * (1 + p * p)
			},
		},
	}
}
