package service

import (
	"auth/domain"
	"auth/jwt"
	"auth/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.IRepository
}

func NewUserService(repo repository.IRepository) UserService {
	return UserService{repo: repo}
}

func (us UserService) CreateUser(ctx context.Context, u domain.User) error {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return err
	}

	createDomain, err := domain.CreateUser(u.Email, hashed, u.Role)
	if err != nil {
		return err
	}

	err = us.repo.Create(ctx, createDomain)
	return err
}

// todo token 이 실제 어디서 사용되야할지 고민하기

var p = jwt.NewProvider("hello warmOil world this is secret key thank you")

func (us UserService) Login(ctx context.Context, reqDomain domain.User) (string, error) {
	getDomains, err := us.repo.GetUser(ctx, reqDomain.Email)
	if err != nil {
		return "", err
	}
	getDomain := getDomains[0]
	if getDomain.Id == 0 {
		return "", errors.New("user not found")
	}
	getDomain, _ = domain.CreateUser(getDomain.Email, getDomain.Password, getDomain.Role)
	fmt.Println("getDo", getDomain)
	isMatched := checkPasswordHash(reqDomain.Password, getDomain.Password)
	fmt.Println("matched getDomain", getDomain)
	if !isMatched {
		return "", errors.New("password not matched")
	}
	token, err := p.CreateToken(getDomain)
	return token, nil
}

func (us UserService) Verify(ctx context.Context, token string) (bool, error) {
	result, err := p.ValidToken(token)
	return result, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
