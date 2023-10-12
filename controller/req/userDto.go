package req

import (
	"auth/domain"
)

type CreateDto struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     []string `json:"role"`
}

func (cd CreateDto) ToDomain() domain.User {
	return domain.User{
		Email:    cd.Email,
		Password: cd.Password,
		Role:     cd.Role,
	}
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ld LoginDto) ToDomain() domain.User {
	return domain.User{
		Email:    ld.Email,
		Password: ld.Password,
	}
}
