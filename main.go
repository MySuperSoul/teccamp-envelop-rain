package main

import (
	"envelop-rain/router"
	"os"
	"os/signal"
	"syscall"
)

// 程序启动之前配置参数
func init() {
	// fmt.Println("Call init")

}

func main() {
	router.APIServerRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.APIServerStop()
}
