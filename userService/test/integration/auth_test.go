package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	handler "github.com/RipperAcskt/innotaxi/internal/handler/restapi"
	"github.com/RipperAcskt/innotaxi/internal/repo/mongo"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/repo/redis"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/go-playground/assert/v2"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func InitHandler() (*handler.Handler, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("config new failed: %w", err)
	}

	postgres, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres new failed: %w", err)
	}

	err = postgres.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		return nil, fmt.Errorf("migrate up failed: %w", err)
	}

	redis, err := redis.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("redis new failed: %w", err)
	}

	mongo, err := mongo.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("mongo new failed: %w", err)
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	writer := zapcore.AddSync(mongo)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	log := zap.New(core, zap.AddCaller())

	service := service.New(postgres, redis, cfg.SALT, cfg)
	return handler.New(service, cfg, log), nil
}

func TestSingUp(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		body string
		code int
		err  error
	}{
		{
			name: "new user",
			body: `{"name": "Ivan", "phone_number": "+7455456", "email": "ripper@algsdh", "password": "12345"}`,
			code: http.StatusCreated,
			err:  nil,
		},
		{
			name: "existed user",
			body: `{"name": "Ivan", "phone_number": "+7455456", "email": "ripper@algsdh", "password": "12345"}`,
			code: http.StatusBadRequest,
			err:  service.ErrUserAlreadyExists,
		},
		{
			name: "empty body",
			body: `{}`,
			code: http.StatusBadRequest,
			err:  fmt.Errorf("EOF"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.POST("/users/auth/sing-up", h.SingUp)

			req, _ := http.NewRequest("POST", "/users/auth/sing-up", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())
		})
	}
}

func TestSingIn(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name string
		body string
		code int
		err  error
	}{
		{
			name: "correct password",
			body: `{"phone_number": "+7455456", "password": "12345"}`,
			code: http.StatusOK,
			err:  nil,
		},
		{
			name: "existed user",
			body: `{"phone_number": "+7455456", "password": "12345787979797979"}`,
			code: http.StatusForbidden,
			err:  service.ErrIncorrectPassword,
		},
		{
			name: "empty body",
			body: `{}`,
			code: http.StatusBadRequest,
			err:  fmt.Errorf("EOF"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.POST("/users/auth/sing-up", h.SingIn)

			req, _ := http.NewRequest("POST", "/users/auth/sing-up", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())
		})
	}

}

func TestRefresh(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name   string
		cookie http.Cookie
		code   int
		err    error
	}{
		{
			name:   "without cookie",
			cookie: http.Cookie{},
			code:   http.StatusForbidden,
			err:    fmt.Errorf("bad refresh token"),
		},
		{
			name: "incorrect cookie",
			cookie: http.Cookie{
				Name:   "refesh_token",
				Value:  "some_token",
				MaxAge: 300,
			},
			code: http.StatusForbidden,
			err:  fmt.Errorf("token parse failed: "),
		},
		{
			name: "correct cookie",
			cookie: http.Cookie{
				Name:   "refesh_token",
				Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzg0MjcyNzAsInVzZXJfaWQiOjEwfQ.xwBqJhor6aQU8cgkREFhg5u-jLNtCBQgon96C1ppgEE",
				MaxAge: int((time.Duration(h.Cfg.REFRESH_TOKEN_EXP) * time.Hour * 24).Seconds()),
			},
			code: http.StatusForbidden,
			err:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.GET("/users/auth/refresh", h.Refresh)

			req, _ := http.NewRequest("GET", "/users/auth/refresh", nil)
			req.AddCookie(&tt.cookie)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())
		})
	}
}

func TestLogout(t *testing.T) {
	h, err := InitHandler()
	if err != nil {
		t.Errorf("init handler failed: %v", err)
	}

	test := []struct {
		name         string
		access_token string
		getAccess    func() string
		code         int
		err          error
	}{
		{
			name: "correct access token",
			getAccess: func() string {
				r := SetUpRouter()
				r.POST("/users/auth/sing-in", h.SingIn)

				req, _ := http.NewRequest("POST", "/users/auth/sing-in", bytes.NewBufferString(`{"phone_number": "+7455456", "password": "12345"}`))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				token := make(map[string]string)
				err := json.Unmarshal(w.Body.Bytes(), &token)
				if err != nil {
					log.Fatalf("json unmarshall failed: %v", err)
				}

				return token["access_token"]
			},
			code: http.StatusOK,
			err:  nil,
		},
		{
			name:         "incorrect access token",
			access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzU4MzcwNzAsInVzZXJfaWQiOjEwfQ.PKg3NU0pwSLFOu1E-gpW2zb8e-X5BDDlv3GGTxg-HmI",
			code:         http.StatusForbidden,
			err:          fmt.Errorf("wrong signature"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r := SetUpRouter()
			r.GET("/users/auth/logout/:id", h.VerifyToken(), h.Logout)

			if tt.getAccess != nil {
				tt.access_token = tt.getAccess()
			}

			req, _ := http.NewRequest("GET", "/users/auth/logout/1", nil)
			req.Header.Add("Authorization", "Bearer "+tt.access_token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.IsEqual(tt.err, w.Body.String())
		})
	}
}
