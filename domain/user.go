package domain

import (
	"errors"
	"net/mail"
	"time"
	"unicode/utf8"
)

type User struct {
	Id        int
	Email     string
	Password  string
	Role      []string
	CreatedAt time.Time
}

func CreateUser(email string, password string, role []string) (User, error) {
	err := validateCreateUser(email, password)
	if err != nil {
		return User{}, err
	}
	return User{
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
	}, nil
}

func validateCreateUser(email string, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email 형식이 올바르지 않습니다")
	}
	if password == "" || utf8.RuneCountInString(password) < 4 {
		return errors.New("비밀번호는 4글자 이상이어야 합니다")
	}
	return nil
}
