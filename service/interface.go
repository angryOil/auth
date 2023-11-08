package service

import (
	"auth/service/req"
	"context"
)

type IUserService interface {
	CreateUser(ctx context.Context, user req.CreateUser) error
	Login(ctx context.Context, user req.Login) (string, error)
	Verify(ctx context.Context, verify string) (bool, error)
}
