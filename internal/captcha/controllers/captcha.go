package controllers

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	pbCaptcha "github.com/mathandcrypto/cryptomath-go-proto/captcha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mathandcrypto/cryptomath-go-captcha/internal/captcha/services"
)

type CaptchaController struct {
	pbCaptcha.UnimplementedCaptchaServiceServer
	captchaService *services.CaptchaService
}

func (c *CaptchaController) GenerateTask(ctx context.Context, req *emptypb.Empty) (*pbCaptcha.GenerateTaskResponse, error) {
	task, err := c.captchaService.GetRandomTask(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate captcha task: %v", err)
	}

	if task == nil {
		return nil, status.Errorf(codes.NotFound, "empty captcha task")
	}

	return &pbCaptcha.GenerateTaskResponse{
		Math: task.Math,
		TaskPayload: &pbCaptcha.TaskPayload{
			Uuid: task.UUID,
			Difficulty: int32(task.Difficulty),
		},
	}, nil
}

func (c *CaptchaController) ValidateTask(ctx context.Context, req *pbCaptcha.ValidateTaskRequest) (*pbCaptcha.ValidateTaskResponse, error) {
	task, err := c.captchaService.GetTask(ctx, req.Uuid, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find captcha task: %v", err)
	}

	if task == nil {
		return nil, status.Errorf(codes.NotFound, "captcha task not found")
	}

	isAnswerCorrect := int32(task.Answer) == req.Answer
	return &pbCaptcha.ValidateTaskResponse{
		IsAnswerCorrect: isAnswerCorrect,
	}, nil
}

func NewCaptchaController(db *sql.DB, rdb *redis.Client) *CaptchaController {
	return &CaptchaController{
		captchaService: services.NewCaptchaService(db, rdb),
	}
}