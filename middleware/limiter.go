/*
 * @Author: your name
 * @Date: 2021-11-06 21:38:50
 * @LastEditTime: 2021-11-06 21:38:50
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /teccamp-envelop-rain/middleware/limiter.go
 */
package middleware

import (
	"net/http"
	"time"

	"envelop-rain/constant"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type RateLimiterMiddleware struct {
	limiter *ratelimit.Bucket
	request int64 // 1 for snatch, 2 for open, 3 for get_wallet_list, defined on constant/constant.go
}

func (r *RateLimiterMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := r.limiter
		if limiter.TakeAvailable(1) < 1 {
			if constant.REQUEST_SNATCH == r.request {
				c.JSON(http.StatusOK, gin.H{"code": constant.SNATCH_BUSY, "msg": constant.SNATCH_BUSY_MESSAGE, "data": gin.H{}})
			} else if constant.REQUEST_OPEN == r.request {
				c.JSON(http.StatusOK, gin.H{"code": constant.OPEN_BUSY, "msg": constant.OPEN_BUSY_MESSAGE, "data": gin.H{}})
			} else {
				c.JSON(http.StatusOK, gin.H{"code": constant.WALLET_BUSY, "msg": constant.WALLET_BUSY_MESSAGE, "data": gin.H{}})
			}
			return
		}
		c.Next()
	}
}

func NewRateLimiter(interval time.Duration, capacity int64, request int64) *RateLimiterMiddleware {
	limiter := ratelimit.NewBucketWithQuantum(interval, capacity, capacity)
	return &RateLimiterMiddleware{
		limiter: limiter,
		request: request,
	}
}
