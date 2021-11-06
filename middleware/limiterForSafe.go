/*
 * @Author: your name
 * @Date: 2021-11-06 21:41:01
 * @LastEditTime: 2021-11-07 00:23:40
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /teccamp-envelop-rain/middleware/limiterForSafe.go
 */
package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type RateKeyFunc func(ctx *gin.Context) (string, error)

type RateLimiterMiddlewareForSafe struct {
	fillInterval time.Duration
	capacity     int64
	ratekeygen   RateKeyFunc
	limiters     map[string]*ratelimit.Bucket
}

func (r *RateLimiterMiddlewareForSafe) get(ctx *gin.Context) (*ratelimit.Bucket, error) {
	key, err := r.ratekeygen(ctx)

	if err != nil {
		return nil, err
	}

	if limiter, existed := r.limiters[key]; existed {
		return limiter, nil
	}

	limiter := ratelimit.NewBucketWithQuantum(r.fillInterval, r.capacity, r.capacity)
	r.limiters[key] = limiter
	return limiter, nil
}

func (r *RateLimiterMiddlewareForSafe) MiddlewareForSafe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter, err := r.get(ctx)
		if err != nil || limiter.TakeAvailable(1) == 0 {
			if err == nil {
				err = errors.New("too many requests")
			}
			ctx.AbortWithError(http.StatusTooManyRequests, err)
		} else {
			ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
			ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
			ctx.Next()
		}
	}
}

func NewRateLimiterForSafe(interval time.Duration, capacity int64, keyGen RateKeyFunc) *RateLimiterMiddlewareForSafe {
	limiters := make(map[string]*ratelimit.Bucket)
	return &RateLimiterMiddlewareForSafe{
		interval,
		capacity,
		keyGen,
		limiters,
	}
}

func NewRateLimiterForIP(interval time.Duration, capacity int64) *RateLimiterMiddlewareForSafe {
	keyGen := func(ctx *gin.Context) (string, error) {
		key := ctx.ClientIP()
		if key != "" {
			return key, nil
		}
		return "", errors.New("unkown IP")
	}
	return NewRateLimiterForSafe(interval, capacity, keyGen)
}
