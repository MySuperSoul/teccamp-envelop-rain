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
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type RateLimiterMiddleware struct {
	limiter *ratelimit.Bucket
}

func (r *RateLimiterMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware enter ...")
		limiter := r.limiter
		if limiter.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}

func NewRateLimiter(interval time.Duration, capacity int64) *RateLimiterMiddleware {
	limiter := ratelimit.NewBucketWithQuantum(interval, capacity, capacity)
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}
