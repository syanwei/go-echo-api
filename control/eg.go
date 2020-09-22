package control

import (
	"context"
	utils "go-echo-api/util"
	"time"

	"github.com/labstack/echo/v4"
)

// eg doc
// @Tags eg
// @Summary eg
// @Accept json
// @Param name formData string true "用戶名" default('')
// @Param email formData string true "邮箱" default('')
// @Router /eg [post]
func Eg(ctx echo.Context) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}{}

	err := ctx.Bind(&user)

	if err != nil {
		return ctx.JSON(utils.ErrIpt("参数有误", err.Error()))
	}

	return ctx.JSON(utils.Succ(`ok`))
}