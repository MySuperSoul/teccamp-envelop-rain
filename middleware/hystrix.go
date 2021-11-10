package middleware

import (
	"envelop-rain/constant"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ConfigHystrix(name string, maxreq int) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                int(3 * time.Second),
		MaxConcurrentRequests:  maxreq,
		RequestVolumeThreshold: 10 * maxreq,
		SleepWindow:            int(2 * time.Second),
		ErrorPercentThreshold:  20,
	})
}

type HystrixMiddleWare struct {
	Name string
}

func (h *HystrixMiddleWare) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hystrix.Do(h.Name, func() error {
			c.Next()
			return nil
		}, func(err error) error {
			logrus.Error("breaker into.")
			if h.Name == "snatch" {
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": constant.SNATCH_BUSY, "msg": constant.SNATCH_BUSY_MESSAGE, "data": gin.H{}})
			} else if h.Name == "open" {
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": constant.OPEN_BUSY, "msg": constant.OPEN_BUSY_MESSAGE, "data": gin.H{}})
			} else {
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": constant.WALLET_BUSY, "msg": constant.WALLET_BUSY_MESSAGE, "data": gin.H{}})
			}
			return nil
		})
	}
}
