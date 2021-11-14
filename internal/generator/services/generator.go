package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	captchaModels "github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/models"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/generator/tasks"
)

type GeneratedTask struct {
	Difficulty	int16
	Math	string
	Answer int
}

type GeneratorService struct {
	ctx context.Context
	db	*gorm.DB
	l	*log.Logger
	chSize int
	batchSize	int
	tmpTasks []captchaModels.CaptchaTask
}

func (s *GeneratorService) insertDatabaseTasks(ctx context.Context) error {
	if err := ctx.Err(); err == nil {
		tx := s.db.WithContext(ctx).Create(s.tmpTasks)
		if tx.Error != nil {
			return fmt.Errorf("batch insert error: %w", tx.Error)
		}
	}

	s.tmpTasks = s.tmpTasks[:0]

	return nil
}

func (s *GeneratorService) createDatabaseTask(ctx context.Context, index int, generatedTask GeneratedTask) error {
	if ctxErr := ctx.Err(); ctxErr == nil {
		if len(s.tmpTasks) < s.batchSize {
			uuidV4, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("failed to generate uuid (v4, index: %d): %w", index, err)
			}

			captchaTask := captchaModels.CaptchaTask{
				Uuid:       uuidV4.String(),
				Difficulty: generatedTask.Difficulty,
				Index:      index,
				Math:       generatedTask.Math,
				Answer:     generatedTask.Answer,
			}

			s.tmpTasks = append(s.tmpTasks, captchaTask)
		} else {
			err := s.insertDatabaseTasks(ctx)
			if err != nil {
				return fmt.Errorf("failed to create captcha tasks (index: %d): %v", index, err)
			}

			return s.createDatabaseTask(ctx, index, generatedTask)
		}
	}

	return nil
}

func (s *GeneratorService) generateTasks(ctx context.Context, generatedTaskCh chan <- GeneratedTask,
	task tasks.Task, generatedParams []int, generateParams[][2]int) {
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
				generatedTask := GeneratedTask{
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
	generatedTaskCh chan <- GeneratedTask, tks *[]tasks.Task) {
	wg := sync.WaitGroup{}

	for _, task := range *tks {
		wg.Add(1)
		t := task
		go func() {
			s.generateTasks(ctx, generatedTaskCh, t, []int{}, t.GenerateParams)
			wg.Done()
		}()
	}

	wg.Wait()
	parentWg.Done()
}

func (s *GeneratorService) allocateGeneration(ctx context.Context, parentWg *sync.WaitGroup, generatedTaskCh chan <- GeneratedTask) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(){
		groupWg := sync.WaitGroup{}
		groupWg.Add(6)
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.CalculusTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.IntegralTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.LimitTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.LogarithmTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.SeriesTasks())
		go s.generateGroupTasks(ctx, &groupWg, generatedTaskCh, tasks.SummingTasks())

		groupWg.Wait()
		wg.Done()
	}()

	wg.Wait()
	close(generatedTaskCh)
	parentWg.Done()
}

func (s *GeneratorService) Start() {
	cancelCtx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	s.l.Printf("started captcha tasks generator")
	generatedTasksCount := 0
	begin := time.Now()

	generatedTaskCh := make(chan GeneratedTask, s.chSize)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go s.allocateGeneration(cancelCtx, &wg, generatedTaskCh)
	go func() {
		s.tmpTasks = make([]captchaModels.CaptchaTask, 0, s.batchSize)

		for task := range generatedTaskCh {
			if err := cancelCtx.Err(); err != nil {
				continue
			}

			err := s.createDatabaseTask(cancelCtx, generatedTasksCount + 1, task)
			if err != nil {
				s.l.Printf("generator error: %v", err)

				cancel()
			}

			if generatedTasksCount > 0 && generatedTasksCount % 1e4 == 0 {
				s.l.Printf("successfully generated and created more than %d tasks", generatedTasksCount)
			}

			generatedTasksCount++
		}

		if err := s.insertDatabaseTasks(cancelCtx); err != nil {
			s.l.Printf("failed to insert remaining batch: %v", err)

			cancel()
		}

		wg.Done()
	}()

	wg.Wait()

	elapsed := time.Since(begin)
	elapsedSec := float64(elapsed.Milliseconds()) / 1000
	s.l.Printf("generated and added to the database %d tasks for %.3f seconds", generatedTasksCount, elapsedSec)
}

func NewGeneratorService(ctx context.Context, db *gorm.DB, l *log.Logger) *GeneratorService {
	return &GeneratorService{
		ctx: ctx,
		db: db,
		l: l,
		chSize: 5000,
		batchSize: 500,
	}
}

