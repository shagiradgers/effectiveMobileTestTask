package handler

import (
	"effectiveMobileTestTask/internal/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetUser(ctx *gin.Context) {
	var userParams store.UserParamsToFilter

	id, ok, err := GetFromQueryInt(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "invalid id param",
		})
		return
	}
	if ok {
		userParams.Id = id
	}

	age, ok, err := GetFromQueryInt(ctx, "age")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "invalid age param",
		})
		return
	}
	if ok {
		userParams.Age = age
	}

	page, ok, err := GetFromQueryUint(ctx, "page")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "invalid page param",
		})
		return
	}
	if ok {
		userParams.Page = page
	}

	limit, ok, err := GetFromQueryUint(ctx, "limit")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "invalid page limit",
		})
		return
	}
	if ok {
		userParams.Limit = limit
	}

	name, ok := ctx.GetQuery("name")
	if ok {
		userParams.Name = name
	}

	surname, ok := ctx.GetQuery("surname")
	if ok {
		userParams.Surname = surname
	}

	patronymic, ok := ctx.GetQuery("patronymic")
	if ok {
		userParams.Patronymic = patronymic
	}

	nationality, ok := ctx.GetQuery("nationality")
	if ok {
		userParams.Nationality = nationality
	}

	sex, ok := ctx.GetQuery("sex")
	if ok {
		userParams.Sex = sex
	}

	users, err := h.UserService.GetUsers(ctx, userParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "err",
			"msg":    "server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "all is ok",
		"res":    users,
	})
}

func (h *Handler) DeleteUser(ctx *gin.Context) {

}

func (h *Handler) EditUser(ctx *gin.Context) {

}

func (h *Handler) AddUser(ctx *gin.Context) {

}
