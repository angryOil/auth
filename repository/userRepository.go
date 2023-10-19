package repository

import (
	"auth/domain"
	"auth/repository/model"
	"context"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db bun.IDB
}

func NewRepository(db bun.IDB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, u domain.User) error {
	m := model.ToModel(u)
	_, err := r.db.NewInsert().Model(&m).Exec(ctx)
	return err
}

func (r UserRepository) GetUser(ctx context.Context, userId string) ([]domain.User, error) {
	var result []model.User
	err := r.db.NewSelect().Model(&result).Where("email=?", userId).Scan(ctx)
	if err != nil {
		return []domain.User{}, err
	}
	return model.ToDomainList(result), nil
}
