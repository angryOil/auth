package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateCreateUser(t *testing.T) {
	err := validateCreateUser("ads", "")

	assert.Contains(t, err.Error(), "email")

	err = validateCreateUser("jn@j", "111")
	assert.Contains(t, err.Error(), "비밀")

	err = validateCreateUser("jipmj12@naver.com", "this is pass word 다")
	assert.NoError(t, err)
}
