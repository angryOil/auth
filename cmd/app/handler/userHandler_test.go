package handler

import (
	"auth/controller/req"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockController struct {
}

func (mc mockController) CreateUser(ctx context.Context, dto req.CreateDto) error {
	if dto.Email == "" {
		return errors.New("email is empty")
	}
	if dto.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (mc mockController) Login(ctx context.Context, dto req.LoginDto) (string, error) {
	// 단순히 controller 에서 예상치 못한 error가 발생할경우의 수를 예상한것
	if dto.Email == "error" {
		return "", errors.New("server error")
	}
	if dto.Email == "jipmj12@naver.com" && dto.Password == "1234" {
		return "success", nil
	}
	return "", errors.New("login fail id or password is wrong")
}

func (mc mockController) Verify(ctx context.Context, token string) (bool, error) {
	if token == "jipmj12@naver.com" {
		return true, nil
	}
	return false, errors.New("invalid token")
}

type UserHandlerTestSuite struct {
	suite.Suite
	h http.Handler
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}

// 전체 테스트 실행전 1번만 실해됨
func (s *UserHandlerTestSuite) SetupSuite() {
	s.h = NewHandler(mockController{})
}

func (s *UserHandlerTestSuite) TestUserHandler_createUser() {
	s.Run("정상적으로 유저를 생성 요청했을때", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()
		reqData := `{"email":"joy@naver.com","password":"1234"}`
		re, err := http.NewRequest("POST", ts.URL+"/users", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)

		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	})
	s.Run("email 값을 공백으로 주었을경우", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()
		reqData := `{"email":"","password":"1234"}`
		re, err := http.NewRequest("POST", ts.URL+"/users", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)

		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("잘못된 json 데이터로 요청했을경우 badRequest를 리턴한다  ", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()
		reqData := `{"wrongEmail":"joy@naver.com","wrongPssword":"1234`
		re, err := http.NewRequest("POST", ts.URL+"/users", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)

		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("데이터가 누락인 경우에도 badRequest를 반환한다.", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()
		reqData := `{}`
		re, err := http.NewRequest("POST", ts.URL+"/users", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)

		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *UserHandlerTestSuite) TestUserHandler_login() {
	s.Run("정상적인 요청시 토큰을 발급한다", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		reqData := `{"email":"jipmj12@naver.com","password":"1234"}`
		re, err := http.NewRequest("POST", ts.URL+"/users/login", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		bodyByteArr, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
		assert.NotZero(s.T(), string(bodyByteArr))
	})
	s.Run("존재하지 않는 id , 비밀번호를 요청하면 badRequest를 반환한다.", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		// 없는 아이디일경우
		reqData := `{"email":"jipmj123@naver.com","password":"1234"}`
		re, err := http.NewRequest("POST", ts.URL+"/users/login", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		bodyByteArr, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
		assert.Contains(s.T(), string(bodyByteArr), "login fail")

		// 비밀번호가 틀렸을경우
		reqData = `{"email":"jipmj12@naver.com","password":"12345"}`
		re, err = http.NewRequest("POST", ts.URL+"/users/login", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err = http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		bodyByteArr, err = io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
		assert.Contains(s.T(), string(bodyByteArr), "login fail")
	})
	s.Run("깨진 json을 보낼경우 badRequest를 반환한다", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		reqData := `{"email":"jipmj12@naver.com","password":"1234}`
		re, err := http.NewRequest("POST", ts.URL+"/users/login", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		bodyByteArr, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "unexpected EOF", string(bodyByteArr))
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("controller 안에 에러가 있을경우", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		reqData := `{"email":"error"}`
		re, err := http.NewRequest("POST", ts.URL+"/users/login", strings.NewReader(reqData))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		bodyByteArr, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(bodyByteArr), "server error")
		assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)
	})
}

func (s *UserHandlerTestSuite) TestUserHandler_verifyToken() {
	s.Run("토큰값이 정상적일경우 statusOk를 반환한다", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		re, err := http.NewRequest("POST", ts.URL+"/users/verify", strings.NewReader("jipmj12@naver.com"))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()

		readBody, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
		assert.Equal(s.T(), string(readBody), "true")
	})
	s.Run("유효하지 않은 토큰 값을 요청할경우 unauthorized 를 반환한다", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		re, err := http.NewRequest("POST", ts.URL+"/users/verify", strings.NewReader("invalid token"))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()

		readBody, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
		assert.Contains(s.T(), string(readBody), "invalid token")
	})
	s.Run("토큰이 빈값일경우 unauthorized 를 반환한다", func() {
		ts := httptest.NewServer(s.h)
		defer ts.Close()

		re, err := http.NewRequest("POST", ts.URL+"/users/verify", nil)
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()

		readBody, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "no token")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
}
