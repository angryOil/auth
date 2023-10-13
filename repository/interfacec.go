package repository

import (
	"auth/domain"
	"context"
)

type IRepository interface {
	Create(ctx context.Context, u domain.User) error
	GetUser(ctx context.Context, userId string) ([]domain.User, error)
}
