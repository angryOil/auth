package controller

import (
	"auth/controller/req"
	"auth/service"
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
	err := c.service.CreateUser(ctx, dto.ToDomain())
	return err
}

func (c UserController) Login(ctx context.Context, dto req.LoginDto) (string, error) {
	fmt.Println("dt", dto)
	token, err := c.service.Login(ctx, dto.ToDomain())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c UserController) Verify(ctx context.Context, token string) (bool, error) {
	result, err := c.service.Verify(ctx, token)
	return result, err
}
