package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserDoesNotExists = fmt.Errorf("user does not exists")
	ErrIncorrectPassword = fmt.Errorf("incorrect password")
)

type UserSingUp struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type UserSingIn struct {
	ID          uint64 `json:"-"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user UserSingUp) error
	CheckUserByPhoneNumber(ctx context.Context, phone string) (*UserSingIn, error)
}

type TokenRepo interface {
	AddToken(token string, expired time.Duration) error
	GetToken(token string) bool
}
type AuthService struct {
	AuthRepo
	TokenRepo
	salt string
	cfg  *config.Config
}

func NewAuthSevice(postgres AuthRepo, redis TokenRepo, salt string, cfg *config.Config) *AuthService {
	return &AuthService{postgres, redis, salt, cfg}
}

func (s *AuthService) SingUp(ctx context.Context, user UserSingUp) error {
	var err error
	user.Password, err = s.GenerateHash(user.Password)
	if err != nil {
		return fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GenerateHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}
	return string(hash.Sum([]byte(s.salt))), nil
}

func (s *AuthService) SingIn(ctx context.Context, user UserSingIn) (*Token, error) {
	userDB, err := s.CheckUserByPhoneNumber(ctx, user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("check user by phone number failed: %w", err)
	}

	hash := sha1.New()
	_, err = hash.Write([]byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("write failed: %w", err)
	}

	if userDB.Password != string(hash.Sum([]byte(s.salt))) {
		return nil, ErrIncorrectPassword
	}

	params := TokenParams{
		ID:                userDB.ID,
		Type:              User,
		HS256_SECRET:      s.cfg.HS256_SECRET,
		ACCESS_TOKEN_EXP:  s.cfg.ACCESS_TOKEN_EXP,
		REFRESH_TOKEN_EXP: s.cfg.REFRESH_TOKEN_EXP,
	}

	token, err := NewToken(params)
	if err != nil {
		return nil, fmt.Errorf("new token failed: %w", err)
	}

	return token, nil
}

func (s *AuthService) Logout(userId string, token string, expired time.Duration) error {
	return s.AddToken(token, expired)
}

func (s *AuthService) CheckToken(userId string) bool {
	return s.GetToken(userId)
}
