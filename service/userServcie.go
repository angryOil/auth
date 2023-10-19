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
	if len(getDomains) == 0 {
		return "", errors.New("login fail user not found")
	}
	getDomain := getDomains[0]
	//getDomain, _ = domain.CreateUser(getDomain.Email, getDomain.Password, getDomain.Role)
	isMatched := checkPasswordHash(reqDomain.Password, getDomain.Password)
	if !isMatched {
		return "", errors.New("login fail password not matched")
	}

	//
	resDomain := toResponseDomain(getDomain)
	token, err := p.CreateToken(resDomain)
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
