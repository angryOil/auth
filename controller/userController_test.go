package controller

import (
	"auth/controller/req"
	req2 "auth/service/req"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// 현재 컨트롤러 부분에서는 단순히 값만 전달후 응답된결과를 다시 전달할 뿐이어서 test 할부분이 딱히없음

type mockService struct {
}

func (ms mockService) CreateUser(ctx context.Context, user req2.CreateUser) error {
	if user.Email == "" {
		return errors.New("email is empty")
	}
	if user.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (ms mockService) Login(ctx context.Context, user req2.Login) (string, error) {
	if user.Email != "jipmj12@naver.com" {
		return "", errors.New("login fail email or password is wrong")
	}
	if user.Password != "1234" {
		return "", errors.New("login fail email or password is wrong")
	}
	return "this.is.mock_token", nil
}

func (ms mockService) Verify(ctx context.Context, verify string) (bool, error) {
	if verify == "token" {
		return true, nil
	}
	return false, errors.New("invalid token")
}

type UserControllerTestSuite struct {
	suite.Suite
	c UserController
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}

func (s *UserControllerTestSuite) SetupSuite() {
	s.c = NewController(mockService{})
}

func (s *UserControllerTestSuite) TestCreateUser() {
	s.Run("유저 생성시 정상적인 요청일경우 error nil 을 반환한다", func() {
		reqUser := req.CreateDto{
			Email:    "jipmj12@naver.com",
			Password: "1234",
		}
		err := s.c.CreateUser(context.Background(), reqUser)
		assert.NoError(s.T(), err)
	})
	s.Run("정상적이지 않은 요청일경우 ERROR 를 반환한다", func() {
		reqUser := req.CreateDto{
			Email:    "",
			Password: "1234",
		}
		err := s.c.CreateUser(context.Background(), reqUser)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "empty")
	})
}

func (s *UserControllerTestSuite) TestUserController_Login() {
	s.Run("정상적인 로그인일경우 token 과 error nil 을 반환한다", func() {
		reqDto := req.LoginDto{
			Email:    "jipmj12@naver.com",
			Password: "1234",
		}
		token, err := s.c.Login(context.Background(), reqDto)
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), token)
	})
	s.Run("아이디 혹은 비밀번호가 틀렸을경우 error 를 반환한다", func() {
		reqDto := req.LoginDto{
			Email:    "wrongEmail@naver.com",
			Password: "1234",
		}
		token, err := s.c.Login(context.Background(), reqDto)
		assert.Error(s.T(), err)
		assert.Zero(s.T(), token)
	})
}

func (s *UserControllerTestSuite) TestUserController_Verify() {
	s.Run("옳바른 토큰을 주었을경우 true와 error nil 값을 반환한다", func() {
		result, err := s.c.Verify(context.Background(), "token")
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})
	s.Run("옳바르지 않은 토큰이 주어졌을경우 false 와 invalid token error를 반환한다.", func() {
		result, err := s.c.Verify(context.Background(), "wrongToken")
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid token")
		assert.False(s.T(), result)
	})
}
