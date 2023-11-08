package controller

import (
	"auth/controller/req"
	"auth/service"
	req2 "auth/service/req"
	"context"
	"fmt"
)

type UserController struct {
	service service.IUserService
}

func NewController(serv service.IUserService) UserController {
	return UserController{service: serv}
}

func (c UserController) CreateUser(ctx context.Context, dto req.CreateDto) error {
	err := c.service.CreateUser(ctx, req2.CreateUser{
		Email:    dto.Email,
		Password: dto.Password,
		Role:     []string{"USER"},
	})
	return err
}

func (c UserController) Login(ctx context.Context, dto req.LoginDto) (string, error) {
	fmt.Println("dt", dto)
	token, err := c.service.Login(ctx, req2.Login{
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c UserController) Verify(ctx context.Context, token string) (bool, error) {
	result, err := c.service.Verify(ctx, token)
	return result, err
}
