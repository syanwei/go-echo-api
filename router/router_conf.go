package router

import (
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	utils "go-echo-api/util"

	"github.com/labstack/echo/v4"
	_ "go-echo-api/docs"
)

// HTTPErrorHandler 全局错误捕捉
func HTTPErrorHandler(err error, ctx echo.Context) {
	if !ctx.Response().Committed {
		if he, ok := err.(*echo.HTTPError); ok {
			ctx.JSON(utils.NewErrSvr("系统错误", he.Message))
		}
	} else {
		ctx.JSON(utils.NewErrSvr("系统错误", err.Error()))
	}
}

// RegDocs 注册文档
func RegDocs(engine *echo.Echo) {
	docUrl := echoSwagger.URL("/swagger/doc.json")
	engine.GET("/swagger/*", echoSwagger.EchoWrapHandler(docUrl))
}

// 跨越配置
var crosConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
}