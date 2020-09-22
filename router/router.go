package router

import (
	"github.com/labstack/echo/v4/middleware"
	"log"

	"github.com/labstack/echo/v4"
)

// RunApp 入口
func RunApp() {
	engine := echo.New()
	engine.Use(middleware.CORSWithConfig(crosConfig)) // 跨域设置
	//engine.HTTPErrorHandler = HTTPErrorHandler        // 自定义错误处理
	RegDocs(engine)                                   // 注册文档

	api := engine.Group("/api")         	  // api/
	apiRouter(api)                      			  // 注册分组路由
	// todo move to config
	err := engine.Start(":8888")
	if err != nil {
		log.Fatalln("run error :", err)
	}
}
