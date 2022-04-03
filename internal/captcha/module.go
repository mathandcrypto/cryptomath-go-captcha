package captcha

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	pbCaptcha "github.com/mathandcrypto/cryptomath-go-proto/captcha"
	"google.golang.org/grpc"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/controllers"
)

func Init(grpcServer *grpc.Server, db *sql.DB, rdb *redis.Client) error {
	captchaController := controllers.NewCaptchaController(db, rdb)

	pbCaptcha.RegisterCaptchaServiceServer(grpcServer, captchaController)

	return nil
}
