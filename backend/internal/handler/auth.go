package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/ablaze/gonexttemp-backend/internal/middleware"
	"github.com/ablaze/gonexttemp-backend/internal/service"
	"github.com/ablaze/gonexttemp-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	RefreshTokenCookie = "refresh_token"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(
			response.CodeValidationError,
			"Invalid request body",
		))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(
			response.CodeValidationError,
			err.Error(),
		))
		return
	}

	authRes, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, response.Error(
				response.CodeConflict,
				"User with this email already exists",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(
			response.CodeInternalError,
			"Failed to register user",
		))
		return
	}

	h.setRefreshTokenCookie(c, authRes.RefreshToken)
	c.JSON(http.StatusCreated, response.Success(gin.H{
		"user":         authRes.User,
		"access_token": authRes.AccessToken,
		"expires_in":   authRes.ExpiresIn,
	}))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(
			response.CodeValidationError,
			"Invalid request body",
		))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(
			response.CodeValidationError,
			err.Error(),
		))
		return
	}

	authRes, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, response.Error(
				response.CodeInvalidCredentials,
				"Invalid email or password",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(
			response.CodeInternalError,
			"Failed to login",
		))
		return
	}

	h.setRefreshTokenCookie(c, authRes.RefreshToken)
	c.JSON(http.StatusOK, response.Success(gin.H{
		"user":         authRes.User,
		"access_token": authRes.AccessToken,
		"expires_in":   authRes.ExpiresIn,
	}))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(
			response.CodeUnauthorized,
			"Refresh token not found",
		))
		return
	}

	authRes, err := h.authService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			h.clearRefreshTokenCookie(c)
			c.JSON(http.StatusUnauthorized, response.Error(
				response.CodeTokenInvalid,
				"Invalid or expired refresh token",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(
			response.CodeInternalError,
			"Failed to refresh token",
		))
		return
	}

	h.setRefreshTokenCookie(c, authRes.RefreshToken)
	c.JSON(http.StatusOK, response.Success(gin.H{
		"user":         authRes.User,
		"access_token": authRes.AccessToken,
		"expires_in":   authRes.ExpiresIn,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenCookie)
	if err == nil {
		_ = h.authService.Logout(c.Request.Context(), refreshToken)
	}

	h.clearRefreshTokenCookie(c)
	c.JSON(http.StatusOK, response.Success(gin.H{
		"message": "Logged out successfully",
	}))
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(
			response.CodeUnauthorized,
			"User not authenticated",
		))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(
			response.CodeValidationError,
			"Invalid user ID",
		))
		return
	}

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error(
				response.CodeNotFound,
				"User not found",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(
			response.CodeInternalError,
			"Failed to get user",
		))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		RefreshTokenCookie,
		token,
		int(7*24*time.Hour.Seconds()), // 7 days
		"/",
		"",
		false, // secure: set to true in production with HTTPS
		true,  // httpOnly
	)
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		RefreshTokenCookie,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}
