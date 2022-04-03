package tasksServices

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/models"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/data"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/models"
)

type GeneratorService struct {
	db	*sql.DB
	l	*logrus.Entry
	chSize int
	batchSize	int
	tmpTasks []models.CaptchaTask
	generatedTasksCount int
}

func (s *GeneratorService) generateTasks(ctx context.Context, generatedTaskCh chan <- tasksModels.GeneratedTask,
	task tasksModels.Task, generatedParams []int, generateParams[][2]int) {
	if err := ctx.Err(); err != nil {
		return
	}

	newGenerateParams := make([][2]int, len(generateParams))
	copy(newGenerateParams, generateParams)

	if len(newGenerateParams) > 0 {
		generateParam := newGenerateParams[0]
		newGenerateParams = newGenerateParams[1:]
		start := generateParam[0]
		end := generateParam[1]

		for i := start; i <= end; i++ {
			newGeneratedParams := make([]int, len(generatedParams))
			copy(newGeneratedParams, generatedParams)

			newGeneratedParams = append(newGeneratedParams, i)

			if len(newGenerateParams) > 0 {
				s.generateTasks(ctx, generatedTaskCh, task, newGeneratedParams, newGenerateParams)
			} else {
				math := task.Math(newGeneratedParams...)
				answer := task.Answer(newGeneratedParams...)
				generatedTask := tasksModels.GeneratedTask{
					Difficulty: task.Difficulty,
					Math: math,
					Answer: answer,
				}

				generatedTaskCh <- generatedTask
			}
		}
	}
}

func (s *GeneratorService) generateGroupTasks(ctx context.Context, parentWg *sync.WaitGroup,
	generatedTaskCh chan <- tasksModels.GeneratedTask, tks *[]tasksModels.Task) {
	defer parentWg.Done()

	wg := sync.WaitGroup{}

	for _, task := range *tks {
		wg.Add(1)
		t := task

		go func() {
			defer wg.Done()

			s.generateTasks(ctx, generatedTaskCh, t, []int{}, t.GenerateParams)
		}()
	}

	wg.Wait()
}

func (s *GeneratorService) prepareCaptchaTasksValues() []string {
	values := make([]string, 0, len(s.tmpTasks))

	for _, task := range s.tmpTasks {
		values = append(values, fmt.Sprintf("('%s', %d, %d, '%s', %d)", task.UUID, task.Difficulty, task.Index, task.Math, task.Answer))
	}

	return values
}

func (s *GeneratorService) insertTempTasksIntoDatabase(ctx context.Context) error {
	if err := ctx.Err(); err == nil {
		_, err = s.db.ExecContext(ctx, fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES %s",
			models.TableNames.CaptchaTasks,
			strings.Join([]string{
				models.CaptchaTaskColumns.UUID,
				models.CaptchaTaskColumns.Difficulty,
				models.CaptchaTaskColumns.Index,
				models.CaptchaTaskColumns.Math,
				models.CaptchaTaskColumns.Answer,
			}, ", "),
			strings.Join(s.prepareCaptchaTasksValues(), ", "),
		))
		if err != nil {
			return fmt.Errorf("failed to insert temp tasks into database: %w", err)
		}
	}

	s.tmpTasks = s.tmpTasks[:0]

	return nil
}

func (s *GeneratorService) insertTaskIntoDatabase(ctx context.Context, index int, generatedTask tasksModels.GeneratedTask) error {
	if err := ctx.Err(); err == nil {
		if len(s.tmpTasks) < s.batchSize {
			uuidV4, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("failed to generate uuid (v4), last index: %d, error: %w", index, err)
			}

			captchaTask := models.CaptchaTask{
				UUID: uuidV4.String(),
				Difficulty: generatedTask.Difficulty,
				Index:      index,
				Math:       generatedTask.Math,
				Answer:     generatedTask.Answer,
			}

			s.tmpTasks = append(s.tmpTasks, captchaTask)
		} else {
			err := s.insertTempTasksIntoDatabase(ctx)
			if err != nil {
				return fmt.Errorf("failed to create captcha tasks, last index: %d, error: %w", index, err)
			}

			return s.insertTaskIntoDatabase(ctx, index, generatedTask)
		}
	}

	return nil
}

func (s *GeneratorService) writeTasks(ctx context.Context, parentWg *sync.WaitGroup,
	generatedTaskCh <-chan tasksModels.GeneratedTask, writeTasksErrorCh chan <- error) {
	defer parentWg.Done()

	s.generatedTasksCount = 0
	s.tmpTasks = make([]models.CaptchaTask, 0, s.batchSize)

	for task := range generatedTaskCh {
		if err := ctx.Err(); err != nil {
			continue
		}

		err := s.insertTaskIntoDatabase(ctx, s.generatedTasksCount + 1, task)
		if err != nil {
			writeTasksErrorCh <- fmt.Errorf("failed to insert captcha tasks: %w", err)

			return
		}

		if s.generatedTasksCount > 0 && s.generatedTasksCount % 1e4 == 0 {
			s.l.Infof("successfully generated and created more than %d tasks", s.generatedTasksCount)
		}

		s.generatedTasksCount++
	}

	if err := s.insertTempTasksIntoDatabase(ctx); err != nil {
		writeTasksErrorCh <- fmt.Errorf("failed to insert remaining captcha tasks batch: %w", err)

		return
	}

	writeTasksErrorCh <- nil
}

func (s *GeneratorService) startGeneration(ctx context.Context, parentWg *sync.WaitGroup,
	generatedTaskCh chan <- tasksModels.GeneratedTask) {
	defer parentWg.Done()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(){
		defer wg.Done()

		groupWg := sync.WaitGroup{}
		groupWg.Add(6)

		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.CalculusTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.IntegralTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.LimitTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.LogarithmTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.SeriesTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasksData.SummingTasks())

		groupWg.Wait()
	}()

	wg.Wait()

	close(generatedTaskCh)
}

func (s *GeneratorService) Start(ctx context.Context) {
	startCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.l.Info("starting captcha tasks generator")

	begin := time.Now()
	generatedTaskCh := make(chan tasksModels.GeneratedTask, s.chSize)
	writeTasksErrorCh := make(chan error, 1)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go s.startGeneration(startCtx, &wg, generatedTaskCh)
	go s.writeTasks(startCtx, &wg, generatedTaskCh, writeTasksErrorCh)

	err := <-writeTasksErrorCh
	if err != nil {
		s.l.WithError(err).Fatal("failed to insert captcha tasks into database")

		cancel()
	}

	wg.Wait()

	elapsed := time.Since(begin)
	elapsedSec := float64(elapsed.Milliseconds()) / 1000
	s.l.Infof("generated and added into database %d tasks for %.3f seconds", s.generatedTasksCount, elapsedSec)
}

func NewGeneratorService(db *sql.DB, l *logrus.Entry) *GeneratorService {
	return &GeneratorService{
		db: db,
		l: l,
		chSize: 5000,
		batchSize: 500,
	}
}