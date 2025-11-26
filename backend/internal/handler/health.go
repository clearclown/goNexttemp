package handler

import (
	"net/http"

	"github.com/ablaze/gonexttemp-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(gin.H{
		"status": "ok",
	}))
}
