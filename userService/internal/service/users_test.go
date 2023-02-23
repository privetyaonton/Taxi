package service_test

import (
	"context"
	"testing"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/RipperAcskt/innotaxi/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestGetProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "get user",
			mockBehavior: func(s *mocks.MockUserRepo) {
				s.EXPECT().GetUserById(context.Background(), "").Return(&model.User{
					Name:        "2",
					PhoneNumber: "2",
					Email:       "2",
					Raiting:     0,
				}, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo)

			service := service.Service{
				UserService: userService,
			}

			_, err := service.GetProfile(context.Background(), "")
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo, user model.User)

	test := []struct {
		name         string
		user         model.User
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "update user",
			user: model.User{
				PhoneNumber: "+77777778",
				Email:       "ripper@mail.ru",
			},
			mockBehavior: func(s *mocks.MockUserRepo, user model.User) {
				s.EXPECT().UpdateUserById(context.Background(), "", &user).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo, tt.user)

			service := service.Service{
				UserService: userService,
			}

			err := service.UpdateProfile(context.Background(), "", &tt.user)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserRepo)

	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "delete user",
			mockBehavior: func(s *mocks.MockUserRepo) {
				s.EXPECT().DeleteUserById(context.Background(), "").Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepo(ctrl)
			userService := service.NewUserService(userRepo)

			tt.mockBehavior(userRepo)

			service := service.Service{
				UserService: userService,
			}

			err := service.DeleteUser(context.Background(), "")
			assert.Equal(t, err, tt.err)
		})
	}
}
