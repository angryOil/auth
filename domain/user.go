package domain

import (
	"auth/domain/vo"
	"errors"
	"net/mail"
	"time"
	"unicode/utf8"
)

var _ User = (*user)(nil)

type User interface {
	ValidCreate() error

	ToGetUser() vo.GetUser
}

type user struct {
	id        int
	email     string
	password  string
	role      []string
	createdAt time.Time
}

func (u *user) ToGetUser() vo.GetUser {
	return vo.GetUser{
		Id:       u.id,
		Email:    u.email,
		Password: u.password,
		Role:     u.role,
	}
}

func (u *user) ValidCreate() error {
	if _, err := mail.ParseAddress(u.email); err != nil {
		return errors.New("email 형식이 올바르지 않습니다")
	}
	if utf8.RuneCountInString(u.password) < 4 {
		return errors.New("password is too short")
	} else if utf8.RuneCountInString(u.password) > 72 {
		return errors.New("password is too long")
	}
	return nil
}
