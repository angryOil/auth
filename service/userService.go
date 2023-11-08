package service

import (
	"auth/domain"
	"auth/jwt"
	"auth/repository"
	req2 "auth/repository/req"
	"auth/service/req"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	//todo 현재 jwt provider를 interface 로 할지 구현체로 그대로 할지 결정
	p    jwt.Provider
	repo repository.IRepository
}

func NewUserService(repo repository.IRepository, p jwt.Provider) UserService {
	return UserService{repo: repo, p: p}
}

const (
	NotMatch = "password or email not matched"
)

func (us UserService) CreateUser(ctx context.Context, c req.CreateUser) error {
	email, password := c.Email, c.Password
	createdAt := time.Now()
	roles := c.Role

	err := domain.NewUserBuilder().
		Email(email).
		Password(password).
		Role(roles).
		CreatedAt(createdAt).
		Build().ValidCreate()
	if err != nil {
		return err
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}

	err = us.repo.Create(ctx, req2.Create{
		Email:     email,
		Password:  hashed,
		Role:      roles,
		CreatedAt: createdAt,
	})
	return err
}

func (us UserService) Login(ctx context.Context, l req.Login) (string, error) {
	email, password := l.Email, l.Password
	getDomains, err := us.repo.GetUser(ctx, req2.GetUser{Email: email})
	if err != nil {
		return "", err
	}

	if len(getDomains) != 1 {
		return "", errors.New(NotMatch)
	}
	getDomain := getDomains[0]
	userVo := getDomain.ToGetUser()
	isMatched := checkPasswordHash(password, userVo.Password)
	if !isMatched {
		return "", errors.New(NotMatch)
	}

	token, err := us.p.CreateToken(jwt.ReqUser{
		UserId: userVo.Id,
		Email:  userVo.Email,
		Role:   userVo.Role,
	})
	return token, nil
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
