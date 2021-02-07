package middleWare

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// https://juejin.im/post/6844903833173229581
// 阻止缓存响应,固定写法
// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}
