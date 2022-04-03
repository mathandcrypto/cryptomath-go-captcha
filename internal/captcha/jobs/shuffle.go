package jobs

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/services"
)

type ShuffleJob struct {
	db	*sql.DB
	rdb *redis.Client
	l	*logrus.Entry
	shufflesCount int64
}

func (j *ShuffleJob) Prepare() error {
	tasksCount := services.GetTasksCount()
	if tasksCount == 0 {
		return errors.New("there are no captcha tasks in the database")
	}

	maxShufflesCount := tasksCount / 3
	shufflesCount, err := rand.Int(rand.Reader, big.NewInt(tasksCount / 3))
	if err != nil {
		return fmt.Errorf("failed to generate shuffle count: %w", err)
	}

	j.shufflesCount = shufflesCount.Int64()

	if j.shufflesCount > maxShufflesCount {
		j.shufflesCount = maxShufflesCount
	}

	return nil
}

func NewShuffleJob(db *sql.DB, rdb *redis.Client, l *logrus.Entry) *ShuffleJob {
	return &ShuffleJob{
		db: db,
		rdb: rdb,
		l: l,
	}
}