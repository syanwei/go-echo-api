package main

import (
	"go-echo-api/router"
	"os"
	"os/signal"
	"syscall"

	logs "github.com/sirupsen/logrus"
)

// @Title go-echo-api Api文档
// @Version 1.0
// @Host 127.0.0.1:8888/api
func main() {
	logs.Info("app initializing")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	go router.RunApp()
	<-quit
	logs.Info("app quited")
}
