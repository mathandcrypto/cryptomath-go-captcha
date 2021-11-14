package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	databaseConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/database"
	captchaModels "github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/models"
	generatorService "github.com/mathandcrypto/cryptomath-go-captcha/internal/generator/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	l := log.New(os.Stdout, "[captcha-seed]: ", log.LstdFlags)
	gl := logger.New(l, logger.Config{
		LogLevel:	logger.Warn,
		SlowThreshold: time.Second,
	})

	//	Init database
	dbConfig, err := databaseConfig.New()
	if err != nil {
		l.Fatalf("failed to load database config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dbConfig.DSN()), &gorm.Config{
		Logger: gl,
	})
	if err != nil {
		l.Fatalf("failed to connect database: %v", err)
	}

	db = db.WithContext(ctx)

	//	Check database
	migrator := db.Migrator()
	tableName := captchaModels.CaptchaTask{}.TableName()
	if !migrator.HasTable(&captchaModels.CaptchaTask{}) {
		l.Fatalf("there is no database table '%s' for captcha tasks", tableName)
	}

	l.Printf("using existing table '%s'", tableName)

	tx := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	if tx.Error != nil {
		l.Fatalf("failed to truncate table '%s', error: %v", tableName, tx.Error)
	}

	l.Printf("truncated table '%s' (affected %d rows)", tableName, tx.RowsAffected)

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

	//	Generate tasks
	generator := generatorService.NewGeneratorService(ctx, db, l)

	generator.Start()

	l.Print("seed finished")
}