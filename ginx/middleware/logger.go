package middleware

import (
	"lcago/log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		traceId, exists := ctx.Get("traceId")
		if !exists {
			traceId = ctx.ClientIP()
		}
		log.AddCtx(ctx, map[string]interface{}{
			"traceId": traceId,
		})

		ctx.Next()
		cost := time.Since(start)
		log.WithContext(ctx).Info(path,
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
