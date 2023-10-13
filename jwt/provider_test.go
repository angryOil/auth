package jwt

import (
	"auth/domain"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var p = Provider{secretKey: "hello_world_this_is_secretKey"}

func TestProvider_CreateToken(t *testing.T) {
	u, err := domain.CreateUser("jipmj@naver.com", "1234", []string{"user"})
	assert.NoError(t, err)
	createdToken, err := p.CreateToken(u)
	assert.NoError(t, err)
	fmt.Println(createdToken)
}

func TestProvider_ValidToken(t *testing.T) {
	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJlbWFpbCI6ImppcG1qMTIzQG5hdmVyLmNvbSIsInJvbGUiOlsiVVNFUiJdLCJleHAiOjE2OTcxNjMwMzQuMzk2OTYyMn0.x7Aov7WFUaUbqS1K4SdrWpNJ5yT0KKx-7v7e5nsRft0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJlbWFpbCI6ImppcG1qQG5hdmVyLmNvbSIsInJvbGUiOlsidXNlciJdLCJleHAiOjE2OTcxNjMyMTAuMTYxODUyOH0.X0NdGy94SpCRi0Of8hHsSU7lAqcfEDfwI5xc9q1qCHY"
	result, err := p.ValidToken(token)
	assert.NoError(t, err)
	assert.True(t, result)
}
