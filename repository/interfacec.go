package repository

import (
	"auth/domain"
	req2 "auth/repository/req"
	"context"
)

type IRepository interface {
	Create(context.Context, req2.Create) error
	GetUser(context.Context, req2.GetUser) ([]domain.User, error)
}
