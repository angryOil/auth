package jwt

import (
	"auth/domain"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
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
	fmt.Println("se", p.secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := token.SignedString([]byte(p.secretKey))
	return signedAuthToken, err
}

func (p Provider) ValidToken(token string) (bool, error) {
	fmt.Println(token)
	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(p.secretKey), nil
	}
	fmt.Println("cert", p.secretKey)
	tok, err := jwt.ParseWithClaims(token, &claims, key)
	return tok.Valid, err
}
