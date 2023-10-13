package req

import (
	"auth/domain"
)

type CreateDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cd CreateDto) ToDomain() domain.User {
	return domain.User{
		Email:    cd.Email,
		Password: cd.Password,
		Role:     []string{"USER"},
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
