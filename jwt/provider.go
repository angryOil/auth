package jwt

import (
	"auth/domain"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"strings"
	"time"
)

type Provider struct {
	secretKey string
}

func NewProvider(secretKey string) Provider {
	return Provider{secretKey: secretKey}
}

type AuthTokenClaims struct {
	UserId int      `json:"user_id"`
	Email  string   `json:"email"`
	Role   []string `json:"role"`
	jwt.StandardClaims
}

func (p Provider) CreateToken(u domain.User) (string, error) {
	at := AuthTokenClaims{
		UserId: u.Id,
		Email:  u.Email,
		Role:   u.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := token.SignedString([]byte(p.secretKey))
	return signedAuthToken, err
}

func (p Provider) ValidToken(token string) (bool, error) {
	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	}

	tok, err := jwt.ParseWithClaims(token, &claims, key)
	return tok.Valid, err
}

func (p Provider) GetPayLoad(token string) (string, error) {
	strs := strings.Split(token, ".")
	fmt.Println(strs)
	if len(strs) != 3 {
		fmt.Println("num:", len(strs))
		return "", errors.New("is not token")
	}

	payload := strs[1]
	result, err := base64.StdEncoding.DecodeString(payload)
	return string(result), err
}
