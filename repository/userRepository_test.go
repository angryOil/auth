package repository

import (
	"auth/repository/infla"
	"auth/repository/model"
	"auth/repository/req"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repository IRepository
	rollback   func() error
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &UserRepositoryTestSuite{})
}

var mockCreateUserDto = req.Create{
	Email:     "joy@naver.com",
	Password:  "1234",
	Role:      []string{"USER"},
	CreatedAt: time.Now(),
}

func (s *UserRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	log.Println("시작전")
	db := infla.NewDB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	s.rollback = tx.Rollback
	s.repository = NewRepository(tx)
	// todo 시작전에 혹시 같은 email의 데이터가 존재하면 삭제 (after 에서 rollback 할거임 물론 commit 자체를 하지도 않을거임  test용 db 추가하기)
	tx.NewDelete().Model(&model.User{}).Where("email = ?", mockCreateUserDto.Email).Exec(context.Background())
}
func (s *UserRepositoryTestSuite) AfterTest(suiteName, testName string) {
	log.Println("after test rollback")
	err := s.rollback()
	if err != nil {
		panic(err)
	}
}

func (s *UserRepositoryTestSuite) TestCreate() {
	s.Run("유저를 생성성공", func() {
		err := s.repository.Create(context.Background(), mockCreateUserDto)
		assert.NoError(s.T(), err)
	})
	s.Run("유저를 생성2 실패(이미 있는 email 로 다시요청)", func() {
		err := s.repository.Create(context.Background(), mockCreateUserDto)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "duplicate")
	})
}

func (s *UserRepositoryTestSuite) TestGetUser() {
	s.Run("존재하지 않는 값을 조회", func() {
		results, err := s.repository.GetUser(context.Background(), req.GetUser{Email: mockCreateUserDto.Email})
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, len(results))
	})
	s.Run("저장되어있는 값을 조회", func() {
		err := s.repository.Create(context.Background(), mockCreateUserDto)
		assert.NoError(s.T(), err)
		results, err := s.repository.GetUser(context.Background(), req.GetUser{Email: mockCreateUserDto.Email})
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1, len(results))
		result := results[0].ToGetUser()
		assert.Equal(s.T(), result.Email, mockCreateUserDto.Email)
		assert.Equal(s.T(), result.Password, mockCreateUserDto.Password)
	})
}
