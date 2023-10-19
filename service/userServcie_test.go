package service

import (
	"auth/domain"
	"auth/jwt"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type mockRepo struct {
}

func (mr mockRepo) Create(ctx context.Context, u domain.User) error {
	if u.Email == "error@err.or" {
		return errors.New("repo error")
	}
	return nil
}

func (mr mockRepo) GetUser(ctx context.Context, userId string) ([]domain.User, error) {
	if strings.Contains(userId, "error") {
		return []domain.User{}, errors.New("repo error")
	}
	if strings.Contains(userId, "jipmj12") {
		mockHashedPw, _ := hashPassword("1234")
		return []domain.User{{
			Id:        rand.Intn(10) + 1,
			Email:     userId,
			Password:  mockHashedPw,
			Role:      []string{"USER"},
			CreatedAt: time.Time{},
		}}, nil
	}
	return []domain.User{}, nil
}

type UserServiceTestSuite struct {
	suite.Suite
	s IUserService
}

var testProvider = jwt.NewProvider("this_is_test_secretKey")

func (s *UserServiceTestSuite) SetupSuite() {
	s.s = NewUserService(mockRepo{}, testProvider)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, &UserServiceTestSuite{})
}

func (s *UserServiceTestSuite) TestUserService_CreateUser() {
	s.Run("정상적으로 유저를 생성요청하면 error nil 을 반환한다", func() {
		err := s.s.CreateUser(context.Background(), domain.User{
			Email:    "jipmj12@naver.com",
			Password: "1234",
			Role:     nil,
		})
		assert.NoError(s.T(), err)
	})
	s.Run("비밀번호가 너무 길경우 error 를 반환한다", func() {
		err := s.s.CreateUser(context.Background(), domain.User{
			Email:    "jipmj12@naver.com",
			Password: "1234lkrl;fea90uif90-uaesoiuafhsjhadsjkadsjkfdskljashkjdfhasfjkhas;dkfjask;jfhajk;hdjkshfadkshasdk;jhaskj",
			Role:     nil,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "password length")
	})
	s.Run("이메일 형식이 아닐경우 , 비밀번호가 빈칸일경우 error를 반환한다", func() {
		// 이메일 형식이 아닐경우
		err := s.s.CreateUser(context.Background(), domain.User{
			Email:    "noEmail",
			Password: "1234",
			Role:     nil,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "email")

		// 이메일 형식이 아닐경우
		err = s.s.CreateUser(context.Background(), domain.User{
			Email:    "jipmj12@naver.com",
			Password: "",
			Role:     nil,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "password")
	})
	s.Run("repo 저장 에러가 발생할경우 error 를 반환한다", func() {
		err := s.s.CreateUser(context.Background(), domain.User{
			Email:    "error@err.or",
			Password: "1234",
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "error")
	})
}

func (s *UserServiceTestSuite) TestUserService_Login() {
	s.Run("로그인 성공할경우", func() {
		result, err := s.s.Login(context.Background(), domain.User{
			Email:    "jipmj12@naver.com",
			Password: "1234",
		})
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), result)
	})
	s.Run("비밀번호가 다를경우", func() {
		result, err := s.s.Login(context.Background(), domain.User{
			Email:    "jipmj12@naver.com",
			Password: "12345",
		})
		assert.Error(s.T(), err)
		assert.Zero(s.T(), result)
		assert.Contains(s.T(), err.Error(), "login fail")
	})
	s.Run("없는 유저일경우", func() {
		result, err := s.s.Login(context.Background(), domain.User{
			Email:    "unknwon@naver.com",
			Password: "1234",
		})
		assert.Error(s.T(), err)
		assert.Zero(s.T(), result)
		assert.Contains(s.T(), err.Error(), "login fail")
	})
	s.Run("repository 에서 error가 발생한경우 ", func() {
		result, err := s.s.Login(context.Background(), domain.User{
			Email:    "error@naver.com",
			Password: "1234",
		})
		assert.Error(s.T(), err)
		assert.Zero(s.T(), result)
	})
}

func (s *UserServiceTestSuite) TestUserService_Verify() {
	s.Run("올바른 토큰일경우 true 와 error nil을 반환한다", func() {
		token, err := testProvider.CreateToken(domain.User{
			Id:       44,
			Email:    "jipmj12@naver.com",
			Password: "1234",
			Role:     []string{"USER"},
		})
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), token)

		result, err := s.s.Verify(context.Background(), token)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})
	s.Run("올바른 토큰이 아닐경우 false 와 error 를 반환한다", func() {
		// .이 3개가 아닐경우
		result, err := s.s.Verify(context.Background(), "invalid.jwt_token")
		assert.Error(s.T(), err)
		assert.False(s.T(), result)

		// .은 3개지만 jwt 형식이 아닐경우
		result, err = s.s.Verify(context.Background(), "invalid.jwt.ken")
		assert.Error(s.T(), err)
		assert.False(s.T(), result)
	})
}
