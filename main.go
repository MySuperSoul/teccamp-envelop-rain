package main

import (
	"envelop-rain/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router.APIServerRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.APIServerStop()
}
