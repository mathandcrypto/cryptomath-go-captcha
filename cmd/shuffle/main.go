package main

import (
	"context"

	dbConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/db"
	redisConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/redis"
	captchaJobs "github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/jobs"
	captchaServices "github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/services"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/common/logger"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/providers"
)

func main() {
	//	Init context
	ctx := context.Background()

	//	Init logger
	l := logger.CreateLogger("captcha-shuffle").WithContext(ctx)

	//	Init database
	dbConf, err := dbConfig.New()
	if err != nil {
		l.WithError(err).Fatal("failed to load database config")
	}

	db, err := providers.NewDBProvider(dbConf)
	if err != nil {
		l.WithError(err).Fatal("failed init database provider")
	}

	//	Init redis
	redisConf, err := redisConfig.New()
	if err != nil {
		l.WithError(err).Fatal("failed to load redis config")
	}

	rdb, err := providers.NewRedisProvider(ctx, redisConf)
	if err != nil {
		l.WithError(err).Fatal("failed to init redis provider")
	}

	//	Count captcha tasks
	err = captchaServices.CountTasks(ctx, db)
	if err != nil {
		l.WithError(err).Fatal("failed to count captcha tasks")
	}

	//	Starting a shuffling
	shuffleJob := captchaJobs.NewShuffleJob(db, rdb, l)

	if err = shuffleJob.Prepare(); err != nil {
		l.WithError(err).Fatal("failed to prepare shuffle job")
	}

	if err = db.Close(); err != nil {
		l.WithError(err).Fatal("failed to close database connection")
	}
}
