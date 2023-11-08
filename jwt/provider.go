package jwt

import (
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

type ReqUser struct {
	UserId int      `json:"user_id"`
	Email  string   `json:"email"`
	Role   []string `json:"role"`
}

type AuthTokenClaims struct {
	ReqUser
	jwt.StandardClaims
}

func (p Provider) CreateToken(u ReqUser) (string, error) {
	at := AuthTokenClaims{
		ReqUser: u,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 30)),
		},
	}
	fmt.Println("se", p.secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := token.SignedString([]byte(p.secretKey))
	return signedAuthToken, err
}

func (p Provider) ValidToken(token string) (bool, error) {
	fmt.Println(token)
	if !tokenCheck(token) {
		return false, errors.New("invalid token")
	}
	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(p.secretKey), nil
	}
	tok, err := jwt.ParseWithClaims(token, &claims, key)
	return tok.Valid, err
}

func tokenCheck(token string) bool {
	if token == "" {
		return false
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	return true
}
