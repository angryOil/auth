package domain

import "time"

var _ UserBuilder = (*userBuilder)(nil)

func NewUserBuilder() UserBuilder {
	return &userBuilder{}
}

type UserBuilder interface {
	Id(id int) UserBuilder
	Email(email string) UserBuilder
	Password(password string) UserBuilder
	Role(role []string) UserBuilder
	CreatedAt(createdAt time.Time) UserBuilder

	Build() User
}

type userBuilder struct {
	id        int
	email     string
	password  string
	role      []string
	createdAt time.Time
}

func (u *userBuilder) Id(id int) UserBuilder {
	u.id = id
	return u
}

func (u *userBuilder) Email(email string) UserBuilder {
	u.email = email
	return u
}

func (u *userBuilder) Password(password string) UserBuilder {
	u.password = password
	return u
}

func (u *userBuilder) Role(role []string) UserBuilder {
	u.role = role
	return u
}

func (u *userBuilder) CreatedAt(createdAt time.Time) UserBuilder {
	u.createdAt = createdAt
	return u
}

func (u *userBuilder) Build() User {
	return &user{
		id:        u.id,
		email:     u.email,
		password:  u.password,
		role:      u.role,
		createdAt: u.createdAt,
	}
}
