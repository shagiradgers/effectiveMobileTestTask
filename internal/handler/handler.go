package handler

import (
	"effectiveMobileTestTask/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Engine      *gin.Engine
	UserService service.UserActions
}

func New(userService service.UserActions) *Handler {
	h := &Handler{
		Engine:      gin.New(),
		UserService: userService,
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
