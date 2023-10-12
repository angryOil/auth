package jwt

import (
	"auth/domain"
	"encoding/json"
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

func TestProvider_GetPayLoad(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJlbWFpbCI6ImppcG1qMTJAbmF2ZXIuY29tIiwicm9sZSI6WyJ1c2VyIl0sImV4cCI6MTY5NzEzMjA2Mi44NzIxMTd9.tAss76ospDXuH5PEm-1_Z4st95YUEhpDlG8sBJIb1wQ"
	result, err := p.GetPayLoad(token)
	assert.NoError(t, err)
	u := domain.User{}
	json.Unmarshal([]byte(result), &u)
	fmt.Println(u)
}
