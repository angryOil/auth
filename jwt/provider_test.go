package jwt

import (
	"auth/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TestProviderSuite struct {
	suite.Suite
	p Provider
}

func TestProviderSuiteTest(t *testing.T) {
	suite.Run(t, &TestProviderSuite{})
}

func (s *TestProviderSuite) SetupSuite() {
	s.p = NewProvider("this_is_test_provider_key_thankYou")
}

func (s *TestProviderSuite) TestProvider_CreateToken() {
	s.Run("토큰 생성후 확인한다 jwt 는 .이 3개다", func() {
		u, err := domain.CreateUser("jipmj@naver.com", "1234", []string{"user"})
		assert.NoError(s.T(), err)
		createdToken, err := s.p.CreateToken(u)
		assert.NoError(s.T(), err)
		splits := strings.Split(createdToken, ".")
		assert.Equal(s.T(), 3, len(splits))
	})

}

func (s *TestProviderSuite) TestProvider_ValidToken() {
	s.Run("토큰 검증", func() {
		token, err := s.p.CreateToken(domain.User{
			Id:    33,
			Email: "test@email.com",
			Role:  nil,
		})
		assert.NoError(s.T(), err)
		result, err := s.p.ValidToken(token)
		assert.NoError(s.T(), err)
		assert.True(s.T(), result)
	})
}
