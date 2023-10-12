package model

import (
	"auth/domain"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id            int       `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Password      int       `bun:"password,notnull"`
	Role          string    `bun:"role"`
	IsDeleted     bool      `bun:"is_deleted,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		Id:        u.Id,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
