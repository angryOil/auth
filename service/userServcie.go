package service

import (
	"auth/domain"
	"auth/jwt"
	"auth/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	//todo 현재 jwt provider를 interface 로 할지 구현체로 그대로 할지 결정
	p    jwt.Provider
	repo repository.IRepository
}

func NewUserService(repo repository.IRepository, p jwt.Provider) UserService {
	return UserService{repo: repo, p: p}
}

func (us UserService) CreateUser(ctx context.Context, u domain.User) error {
	if len(u.Password) < 4 {
		return errors.New("password is too short")
	}
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

func (us UserService) Login(ctx context.Context, reqDomain domain.User) (string, error) {
	getDomains, err := us.repo.GetUser(ctx, reqDomain.Email)
	if err != nil {
		return "", err
	}
	if len(getDomains) == 0 {
		return "", errors.New("login fail user not found")
	}
	getDomain := getDomains[0]
	isMatched := checkPasswordHash(reqDomain.Password, getDomain.Password)
	if !isMatched {
		return "", errors.New("login fail password not matched")
	}

	resDomain := toResponseDomain(getDomain)
	token, err := us.p.CreateToken(resDomain)
	return token, nil
}

func toResponseDomain(u domain.User) domain.User {
	return domain.User{
		Id:    u.Id,
		Email: u.Email,
		Role:  u.Role,
	}
}

func (us UserService) Verify(ctx context.Context, token string) (bool, error) {
	result, err := us.p.ValidToken(token)
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
