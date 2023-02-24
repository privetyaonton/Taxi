package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/go-playground/assert/v2"
)

func TestGetProfile(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		id   string
		user string
		code int
		err  error
	}{
		{
			name: "user does not exist",
			id:   "2",
			user: "",
			code: http.StatusBadRequest,
			err:  service.ErrUserDoesNotExists,
		},
		{
			name: "existed user",
			id:   "1",
			user: `{"name":"Ivan","phone_number":"+7455456","email":"ripper@algsdh","raiting":0}`,
			code: http.StatusOK,
			err:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.GET("/users/profile/:id", h.GetProfile)

			req, _ := http.NewRequest("GET", "/users/profile/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			if tt.err != nil {
				assert.IsEqual(tt.err, w.Body.String())
				return
			}
			assert.Equal(t, tt.user, w.Body.String())
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		id   string
		body string
		code int
		err  error
	}{
		{
			name: "user does not exist",
			id:   "2",
			body: "",
			code: http.StatusBadRequest,
			err:  service.ErrUserDoesNotExists,
		},
		{
			name: "existed user",
			id:   "1",
			body: `{"phone_number": "+77777778","email": "ripper@mail.ru"}`,
			code: http.StatusOK,
			err:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.PUT("/users/profile/:id", h.UpdateProfile)

			req, _ := http.NewRequest("PUT", "/users/profile/"+tt.id, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())

		})
	}
}

func TestDeleteUser(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		id   string
		code int
		err  error
	}{
		{
			name: "user does not exist",
			id:   "2",
			code: http.StatusBadRequest,
			err:  service.ErrUserDoesNotExists,
		},
		{
			name: "existed user",
			id:   "1",
			code: http.StatusOK,
			err:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.DELETE("/users/:id", h.DeleteUser)

			req, _ := http.NewRequest("DELETE", "/users/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())

		})
	}
}
