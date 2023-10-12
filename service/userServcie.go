package service

import (
	"auth/domain"
	"auth/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func (us UserService) CreateUser(ctx context.Context, u domain.User) error {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return err
	}

	encDomain, err := domain.CreateUser(u.Email, hashed, u.Role)
	if err != nil {
		return err
	}

	err = us.repo.Create(encDomain)
	return err
}

func (us UserService) Login(ctx context.Context, reqDomain domain.User) (string, error) {
	getDomain := us.repo.GetUser(reqDomain.Email)

	if getDomain.Id == 0 {
		return "", errors.New("user not found")
	}

	isMatched := checkPasswordHash(reqDomain.Password, getDomain.Password)
	fmt.Println(isMatched)
	return "", nil
}

func (us UserService) Verify(ctx context.Context, verify string) error {
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
