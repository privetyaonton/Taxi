package service

import (
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/golang-jwt/jwt"
)

const (
	User   = "user"
	Driver = "driver"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
	ErrUnknownType  = fmt.Errorf("unknown type")
)

type Token struct {
	Access           string `json:"access_token"`
	RT               string `json:"refresh_token"`
	AccessExpiration time.Time
	RTExpiration     time.Time
}

type TokenParams struct {
	ID                any
	Type              string
	HS256_SECRET      string
	ACCESS_TOKEN_EXP  int
	REFRESH_TOKEN_EXP int
}

func NewToken(params TokenParams) (*Token, error) {
	if params.Type != User && params.Type != Driver {
		return nil, ErrUnknownType
	}

	accessExp := time.Now().Add(time.Duration(params.ACCESS_TOKEN_EXP) * time.Minute)

	access, err := newJwt(accessExp, params)
	if err != nil {
		return nil, fmt.Errorf("new jwt failed: %w", err)
	}

	rtExp := time.Now().Add(time.Duration(params.REFRESH_TOKEN_EXP) * 24 * time.Hour)

	rt, err := newJwt(rtExp, params)
	if err != nil {
		return nil, fmt.Errorf("new rt failed: %w", err)
	}

	return &Token{access, rt, accessExp, rtExp}, nil
}

func newJwt(jwtExp time.Time, p TokenParams) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = p.ID
	claims["type"] = p.Type
	claims["exp"] = jwtExp.UTC().Unix()

	secret := []byte(p.HS256_SECRET)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signed string failed: %w", err)
	}

	return tokenString, nil
}

func Verify(token string, cfg *config.Config) (uint64, error) {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.HS256_SECRET), nil
		},
	)

	if err != nil {
		return 0, fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("jwt map claims failed")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return 0, ErrTokenExpired
	}
	if string(claims["type"].(string)) != User {
		return 0, ErrUnknownType
	}
	return uint64(claims["user_id"].(float64)), nil
}
