package handler

import (
	"effectiveMobileTestTask/internal/store"
	"effectiveMobileTestTask/pkg/age_api"
	"effectiveMobileTestTask/pkg/gender_api"
	"effectiveMobileTestTask/pkg/nationality_api"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	users, err := h.userService.GetUsers(ctx, userParams)
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

type AddUserBodyRequest struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic" binding:"-"`
}

func (h *Handler) AddUser(ctx *gin.Context) {
	var requestBody AddUserBodyRequest

	if err := ctx.BindJSON(&requestBody); err != nil {
		logrus.Debugf("(handler)[Add user] while bind body error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    " invalid body",
		})
	}

	age, err := age_api.GetAge(requestBody.Name)
	if err != nil {
		if errors.Is(err, &age_api.AgeNotFound{}) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"msg":    "age not found",
			})
			return
		}

		logrus.Debugf("(handler)[Add user] while get age error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "err",
			"msg":    "server error",
		})
	}

	nationality, err := nationality_api.GetNationality(requestBody.Name)
	if err != nil {
		if errors.Is(err, &nationality_api.NationalityNotFound{}) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"msg":    "nationality not found",
			})
			return
		}

		logrus.Debugf("(handler)[Add user] while get nationality error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "err",
			"msg":    "server error",
		})
	}

	gender, err := gender_api.GetGender(requestBody.Name)
	if err != nil {
		if errors.Is(err, &gender_api.GenderNotFound{}) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"msg":    "gender not found",
			})
			return
		}

		logrus.Debugf("(handler)[Add user] while get gender error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "err",
			"msg":    "server error",
		})
	}

	err = h.userService.AddUser(ctx, store.UserParamsToAdd{
		Name:        requestBody.Name,
		Surname:     requestBody.Surname,
		Patronymic:  requestBody.Patronymic,
		Sex:         gender,
		Nationality: nationality,
		Age:         int16(age),
	})
}
