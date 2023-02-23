package postgres_test

import (
	"context"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/go-playground/assert/v2"
)

func TestCreateUser(t *testing.T) {
	test := []struct {
		name string
		user service.UserSingUp
		err  error
	}{
		{
			name: "add user",
			user: service.UserSingUp{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
				Password:    "12345",
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("sqlmock new failed: %v", err)
			}

			mock.ExpectQuery("SELECT name FROM users").WithArgs(tt.user.PhoneNumber, tt.user.Email, model.StatusCreated).WillReturnError(nil)
			mock.ExpectExec("INSERT INTO users").WithArgs(tt.user.Name, tt.user.PhoneNumber, tt.user.Email, []byte(tt.user.Password), model.StatusCreated).WillReturnResult(sqlmock.NewResult(1, 1))

			postgres := &postgres.Postgres{
				DB: db,
			}

			err = postgres.CreateUser(context.Background(), tt.user)
			assert.Equal(t, err, tt.err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestCheckUserByPhoneNumber(t *testing.T) {
	test := []struct {
		name         string
		phone_number string
		err          error
	}{
		{
			name:         "get user",
			phone_number: "+7455456",
			err:          nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("sqlmock new failed: %v", err)
			}

			rows := sqlmock.NewRows([]string{"id", "phone_number", "password"}).
				AddRow(1, "123", "123")
			mock.ExpectQuery("SELECT id, phone_number, password FROM users").WithArgs(tt.phone_number, model.StatusCreated).WillReturnRows(rows)

			postgres := &postgres.Postgres{
				DB: db,
			}

			_, err = postgres.CheckUserByPhoneNumber(context.Background(), tt.phone_number)
			assert.Equal(t, err, tt.err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestUpdateUserById(t *testing.T) {
	test := []struct {
		name string
		user model.User
		rows int64
		err  error
	}{
		{
			name: "user exists",
			user: model.User{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
			},
			rows: 1,
			err:  nil,
		},
		{
			name: "user does not exist",
			user: model.User{
				Name:        "Ivan",
				PhoneNumber: "+7455456",
				Email:       "ripper@algsdh",
			},
			rows: 0,
			err:  service.ErrUserDoesNotExists,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("sqlmock new failed: %v", err)
			}

			mock.ExpectExec("UPDATE users").WithArgs(tt.user.Name, tt.user.PhoneNumber, tt.user.Email, "0", model.StatusCreated).WillReturnError(nil).WillReturnResult(sqlmock.NewResult(tt.rows, tt.rows))

			postgres := &postgres.Postgres{
				DB: db,
			}

			err = postgres.UpdateUserById(context.Background(), "0", &tt.user)
			assert.Equal(t, err, tt.err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, err, nil)
		})
	}
}

func TestGetUserById(t *testing.T) {
	test := []struct {
		name string
		err  error
	}{
		{
			name: "update user",
			err:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("sqlmock new failed: %v", err)
			}

			rows := sqlmock.NewRows([]string{"id", "name", "phone_number", "email", "raiting"}).
				AddRow(1, "123", "123", "123", "123")
			mock.ExpectQuery("SELECT id, name, phone_number, email, raiting FROM users").WithArgs("0", model.StatusCreated).WillReturnRows(rows)

			postgres := &postgres.Postgres{
				DB: db,
			}

			_, err = postgres.GetUserById(context.Background(), "0")
			assert.Equal(t, err, tt.err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, err, nil)
		})
	}
}

func TestDeleteUserById(t *testing.T) {
	test := []struct {
		name string
		rows int64
		err  error
	}{
		{
			name: "user exists",
			rows: 1,
			err:  nil,
		},
		{
			name: "user does not exist",
			rows: 0,
			err:  service.ErrUserDoesNotExists,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("sqlmock new failed: %v", err)
			}

			mock.ExpectExec("UPDATE users").WithArgs(model.StatusDeleted, "0", model.StatusCreated).WillReturnError(nil).WillReturnResult(sqlmock.NewResult(tt.rows, tt.rows))

			postgres := &postgres.Postgres{
				DB: db,
			}

			err = postgres.DeleteUserById(context.Background(), "0")
			assert.Equal(t, err, tt.err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, err, nil)
		})
	}
}
