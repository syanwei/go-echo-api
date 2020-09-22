package router

import (
	"go-echo-api/control"

	"github.com/labstack/echo/v4"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	api.POST(`/eg`, control.Eg)            // 简单例子
}
