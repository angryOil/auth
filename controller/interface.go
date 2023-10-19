package controller

import (
	"auth/controller/req"
	"context"
)

type IController interface {
	CreateUser(ctx context.Context, dto req.CreateDto) error
	Login(ctx context.Context, dto req.LoginDto) (string, error)
	Verify(ctx context.Context, token string) (bool, error)
}
