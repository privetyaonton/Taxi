package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/RipperAcskt/innotaxi/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestSingUp(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthRepo, user service.UserSingUp)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		user         service.UserSingUp
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "correct user",
			user: service.UserSingUp{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
				Password:    "12345",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, user service.UserSingUp) {
				s.EXPECT().CreateUser(context.Background(), user).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}

			service := service.Service{
				AuthService: service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{}),
			}

			tmpPass := tt.user.Password
			tt.user.Password, _ = service.GenerateHash(tt.user.Password)
			tt.mockBehavior(f.authRepo, tt.user)

			tt.user.Password = tmpPass
			err := service.SingUp(context.Background(), tt.user)
			assert.IsEqual(err, tt.err)
		})
	}
}

func TestSingIn(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuthRepo, phone_number string)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		user         service.UserSingIn
		mockBehavior mockBehavior
		token        string
		err          error
	}{
		{
			name: "correct password",
			user: service.UserSingIn{
				PhoneNumber: "2",
				Password:    "2",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckUserByPhoneNumber(context.Background(), phone_number).Return(&service.UserSingIn{
					ID:          9,
					PhoneNumber: "2",
					Password:    string([]byte{49, 50, 52, 106, 107, 104, 115, 100, 97, 102, 51, 52, 50, 53, 218, 75, 146, 55, 186, 204, 205, 241, 156, 7, 96, 202, 183, 174, 196, 168, 53, 144, 16, 176}),
				}, nil)
			},
			token: "",
			err:   nil,
		},
		{
			name: "incorrect password",
			user: service.UserSingIn{
				PhoneNumber: "+7455456",
				Password:    "123456",
			},
			mockBehavior: func(s *mocks.MockAuthRepo, phone_number string) {
				s.EXPECT().CheckUserByPhoneNumber(context.Background(), phone_number).Return(&service.UserSingIn{}, nil)
			},
			token: "",
			err:   fmt.Errorf("incorrect password"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			tt.mockBehavior(f.authRepo, tt.user.PhoneNumber)

			service := service.Service{
				AuthService: authService,
			}

			token, err := service.SingIn(context.Background(), tt.user)
			assert.NotEqual(t, token, tt.token)
			assert.Equal(t, err, tt.err)
		})
	}

}

func TestGenerateHash(t *testing.T) {
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name     string
		password string
		hash     string
		err      error
	}{
		{
			name:     "password",
			password: "2",
			hash:     string([]byte{49, 50, 52, 106, 107, 104, 115, 100, 97, 102, 51, 52, 50, 53, 218, 75, 146, 55, 186, 204, 205, 241, 156, 7, 96, 202, 183, 174, 196, 168, 53, 144, 16, 176}),
			err:      nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			service := service.Service{
				AuthService: authService,
			}

			hash, err := service.GenerateHash(tt.password)
			assert.Equal(t, hash, tt.hash)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestLogout(t *testing.T) {
	type mockBehavior func(s *mocks.MockTokenRepo)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "logout",
			mockBehavior: func(s *mocks.MockTokenRepo) {
				s.EXPECT().AddToken("", time.Duration(123)).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			service := service.Service{
				AuthService: authService,
			}

			tt.mockBehavior(f.tokenRepo)
			err := service.Logout("", "", time.Duration(123))
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestCheckToken(t *testing.T) {
	type mockBehavior func(s *mocks.MockTokenRepo)
	type fileds struct {
		authRepo  *mocks.MockAuthRepo
		tokenRepo *mocks.MockTokenRepo
	}
	test := []struct {
		name         string
		mockBehavior mockBehavior
		exist        bool
	}{
		{
			name: "check token",
			mockBehavior: func(s *mocks.MockTokenRepo) {
				s.EXPECT().GetToken("0").Return(false)
			},
			exist: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fileds{
				authRepo:  mocks.NewMockAuthRepo(ctrl),
				tokenRepo: mocks.NewMockTokenRepo(ctrl),
			}
			authService := service.NewAuthSevice(f.authRepo, f.tokenRepo, "124jkhsdaf3425", &config.Config{})

			service := service.Service{
				AuthService: authService,
			}

			tt.mockBehavior(f.tokenRepo)
			err := service.CheckToken("0")
			assert.Equal(t, err, tt.exist)
		})
	}
}
func TestVerify(t *testing.T) {
	cfg := &config.Config{
		HS256_SECRET: "QWERTfg53gxb2",
	}

	test := []struct {
		name   string
		token  string
		userId uint64
		err    error
	}{
		{
			name:   "verify token expired",
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY4Nzk5NDIsInR5cGUiOiJ1c2VyIiwidXNlcl9pZCI6MX0.qwiL4bupjm9O-ZnKpIcB8-erQytBJgkWlxnwPmRmv-c",
			userId: 0,
			err:    service.ErrTokenExpired,
		},
		{
			name:   "verify token ok",
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzk0NzAxNTQsInR5cGUiOiJ1c2VyIiwidXNlcl9pZCI6MX0.r5vZu9eOds5kti9UjQFXx8AYLHZC23YLtVVnr8dgx24",
			userId: 1,
			err:    nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id, err := service.Verify(tt.token, cfg)
			assert.IsEqual(err, tt.err)
			assert.Equal(t, id, tt.userId)
		})
	}
}
