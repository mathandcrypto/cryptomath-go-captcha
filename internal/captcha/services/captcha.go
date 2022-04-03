package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/go-redis/redis/v8"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/models"
)

type CaptchaService struct {
	db *sql.DB
	rdb	 *redis.Client
}

var tasksCount int64

func GetTasksCount() int64 {
	return tasksCount
}

func CountTasks(ctx context.Context, db *sql.DB) error {
	count, err := models.CaptchaTasks().Count(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to count captcha tasks: %w", err)
	}

	tasksCount = count

	return nil
}

func (s *CaptchaService) GetTask(ctx context.Context, uuid string, retry bool) (*models.CaptchaTask, error) {
	task, err := models.FindCaptchaTask(ctx, s.db, uuid)
	if err == sql.ErrNoRows {
		if !retry {
			removedUuid, redisErr := s.rdb.Get(ctx, uuid).Result()
			if redisErr != nil {
				return s.GetTask(ctx, removedUuid, false)
			}
		}

		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to find captcha task (uuid: %s), error: %w", uuid, err)
	}

	return task, nil
}

func (s *CaptchaService) GetRandomTask(ctx context.Context) (*models.CaptchaTask, error) {
	if tasksCount == 0 {
		return nil, nil
	}

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(tasksCount))
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}

	randomIndex := randomNumber.Int64() + 1
	task, err := models.CaptchaTasks(models.CaptchaTaskWhere.Index.EQ(randomIndex)).One(ctx, s.db)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get random captcha task (index: %d), error: %w", randomIndex, err)
	}

	return task, nil
}

func NewCaptchaService(db *sql.DB, rdb *redis.Client) *CaptchaService {
	return &CaptchaService{
		db: db,
		rdb: rdb,
	}
}