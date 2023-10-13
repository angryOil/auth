package model

import (
	"auth/domain"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id            int       `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Password      string    `bun:"password,notnull"`
	Role          string    `bun:"role"`
	IsDeleted     bool      `bun:"is_deleted,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`
}

func ToModel(d domain.User) User {
	return User{Id: d.Id,
		Email:         d.Email,
		Password:      d.Password,
		Role:          strings.Join(d.Role, ","),
		IsDeleted:     false,
		CreatedAt:     d.CreatedAt,
		LastUpdatedAt: time.Time{},
	}
}

func (u User) ToDomain() domain.User {
	return domain.User{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		Role:      strings.Split(u.Role, ","),
		CreatedAt: u.CreatedAt,
	}
}

func ToDomainList(list []User) []domain.User {
	result := make([]domain.User, len(list))
	for i, u := range list {
		result[i] = u.ToDomain()
	}
	return result
}
