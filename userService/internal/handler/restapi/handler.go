package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/RipperAcskt/innotaxi/config"
	_ "github.com/RipperAcskt/innotaxi/docs"
	"github.com/RipperAcskt/innotaxi/internal/service"
)

type Handler struct {
	s   *service.Service
	Cfg *config.Config
	log *zap.Logger
}

func New(s *service.Service, cfg *config.Config, log *zap.Logger) *Handler {
	return &Handler{s, cfg, log}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	users := router.Group("/users")
	users.Use(h.Log())

	auth := users.Group("/auth")
	auth.POST("sing-up", h.SingUp)
	auth.POST("sing-in", h.SingIn)
	auth.GET("refresh", h.Refresh)
	auth.GET("logout", h.VerifyToken(), h.Logout)

	users.GET("/profile/:id", h.VerifyToken(), h.GetProfile)
	users.PUT("/profile/:id", h.VerifyToken(), h.UpdateProfile)
	users.DELETE("/:id", h.VerifyToken(), h.DeleteUser)

	return router
}
