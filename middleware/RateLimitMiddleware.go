package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"share.ac.cn/response"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Fail(c, nil, "请求太频繁了，休息一下吧~")
			c.Abort()
			return
		}
		c.Next()
	}
}
