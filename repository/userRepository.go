package repository

import "auth/domain"

type UserRepository struct {
	name string
}

func (r UserRepository) Create(u domain.User) error {
	return nil
}

func (r UserRepository) GetUser(userId string) domain.User {
	return domain.User{}
}
