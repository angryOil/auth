package repository

import (
	"auth/domain"
	"auth/repository/model"
	"auth/repository/req"
	"context"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db bun.IDB
}

func NewRepository(db bun.IDB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, c req.Create) error {
	m := model.ToCreateModel(c)
	_, err := r.db.NewInsert().Model(&m).Exec(ctx)
	return err
}

func (r UserRepository) GetUser(ctx context.Context, u req.GetUser) ([]domain.User, error) {
	var result []model.User
	err := r.db.NewSelect().Model(&result).Where("email=?", u.Email).Scan(ctx)
	if err != nil {
		return []domain.User{}, err
	}
	return model.ToDomainList(result), nil
}
