package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)
	r.Run()
}

func SnatchHandler(c *gin.Context) {

}

func OpenHandler(c *gin.Context) {

}

func WalletListHandler(c *gin.Context) {

}
