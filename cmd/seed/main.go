package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dbConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/db"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/common/logger"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/models"
	"github.com/mathandcrypto/cryptomath-go-captcha/internal/providers"
	tasksServices "github.com/mathandcrypto/cryptomath-go-captcha/internal/tasks/services"
)

func main() {
	//	Init context
	ctx, cancel := context.WithCancel(context.Background())

	//	Init logger
	l := logger.CreateLogger("captcha-seed").WithContext(ctx)

	//	Init database
	dbConf, err := dbConfig.New()
	if err != nil {
		l.WithError(err).Fatal("failed to load database config")
	}

	db, err := providers.NewDBProvider(dbConf)
	if err != nil {
		l.WithError(err).Fatal("failed init database provider")
	}

	//	Removing all old captcha tasks from the database
	_, err = db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s", models.TableNames.CaptchaTasks))
	if err != nil {
		l.WithError(err).Fatal("failed to truncate captcha tasks table")
	}

	//	Subscribe to system signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		stop := <-signalChan
		l.Printf("signal: %v, waiting for all processes to stop", stop)

		cancel()
	}()

	//	Generate captcha tasks
	generator := tasksServices.NewGeneratorService(db, l)

	generator.Start(ctx)

	if err = db.Close(); err != nil {
		l.WithError(err).Fatal("failed to close database connection")
	}
}
