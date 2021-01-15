package middleWare

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/log"
	"time"
)

type UserRequest struct {
	ClientIP   string
	Host       string
	Url        string
	StatusCode int
	Method     string
	Agent      []string
}

func AccessLog(notLogPath ...string) gin.HandlerFunc {
	var skipLogUrl map[string]bool
	if length := len(notLogPath); length > 0 {
		skipLogUrl = make(map[string]bool, length)
		for _, url := range notLogPath {
			skipLogUrl[url] = true
		}
	}

	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		req := ctx.Request
		ctx.Next()

		_, ok := skipLogUrl[req.URL.Path]
		if !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			raw := req.URL.RawPath
			url := req.URL.Path

			if raw != "" {
				url = url + "?" + raw
			}

			req := &UserRequest{
				ClientIP:   ctx.ClientIP(),
				Host:       req.Host,
				Url:        url,
				StatusCode: ctx.Writer.Status(), // http code
				Method:     req.Method,
				Agent:      req.Header.Values("User-Agent"),
			}
			logger := log.GetLogger()
			// - 在右侧而非左侧填充空格（左对齐该区域）
			logger.Infof("ip: %-15s dest_addr: %-15s latency: %-10v code: %-5d method: %-8s path: %-8s agent: %-10v \n", req.ClientIP, req.Host, latency, req.StatusCode, req.Method, req.Url, req.Agent)
		}

	}
}
