package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxi/internal/model"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary get user profile
// @Tags user
// @Param id path int true "user's id"
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/profile/{id} [GET]
// @Security Bearer
func (h *Handler) GetProfile(c *gin.Context) {
	logger := getLogger(c)

	user, err := h.s.GetProfile(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		logger.Error("/users/profile/{id}", zap.Error(fmt.Errorf("get profile failed: %w", err)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get profile failed: %w", err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary update user profile
// @Tags user
// @Param input body model.User false "rows to update"
// @Param id path int true "user's id"
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/profile/{id} [PUT]
// @Security Bearer
func (h *Handler) UpdateProfile(c *gin.Context) {
	logger := getLogger(c)

	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.s.UpdateProfile(c.Request.Context(), c.Param("id"), &user)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		logger.Error("/users/profile/{id}", zap.Error(fmt.Errorf("update profile failed: %w", err)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

// @Summary delete user
// @Tags user
// @Param id path int false "user's id to delete"
// @Accept json
// @Success 200
// @Failure 401 {object} error "error: err"
// @Failure 403 {object} error "error: err"
// @Failure 500 {object} error "error: err"
// @Router /users/{id} [DELETE]
// @Security Bearer
func (h *Handler) DeleteUser(c *gin.Context) {
	logger := getLogger(c)

	err := h.s.DeleteUser(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		logger.Error("/users/{id}", zap.Error(fmt.Errorf("delete user failed: %w", err)))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(200)
}
