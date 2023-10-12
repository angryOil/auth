package controller

import (
	"auth/controller/req"
	"auth/service"
	"context"
	"encoding/json"
)

type UserController struct {
	uService service.IUserService[string]
}

func NewController() UserController {
	return UserController{}
}

func (c UserController) CreateUser(ctx context.Context, dto req.CreateDto) (string, error) {
	data, err := json.Marshal(dto)

	return string(data), err
}
