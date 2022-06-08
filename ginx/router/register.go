package router

import (
	"lcago/ginx/middleware"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options = []Option{}

// Register 注册app的路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.New()
	//gin.Default()
	// 添加中间件,日志和recover
	r.Use(middleware.GinRecovery(true), middleware.GinLogger())
	// 实例化http server
	for _, opt := range options {
		opt(r)
	}
	return r
}
