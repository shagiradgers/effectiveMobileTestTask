package handler

import (
	"effectiveMobileTestTask/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log         *logrus.Logger
	Engine      *gin.Engine
	userService service.UserActions
}

func New(userService service.UserActions) *Handler {
	h := &Handler{
		Engine:      gin.New(),
		userService: userService,
	}

	users := h.Engine.Group("users")
	{
		users.GET("", h.GetUser)
		users.DELETE("", h.DeleteUser)
		users.PUT("", h.EditUser)
		users.POST("", h.AddUser)
	}

	return h
}
