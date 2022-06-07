package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

// Run 启动api服务
func Run(ctx context.Context) {
	// panic 自动恢复
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	start := time.Now()

	// Set HTTP server mode.
	//if conf.HttpMode() != "" {
	//	gin.SetMode(conf.HttpMode())
	//} else if conf.Debug() == false {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	// http router engine
	engine := gin.New()

	// 添加中间件
	// engine.Use(middleware.GinRecovery(),middleware.GinLogger())

	// 注册路由方法
	registerRouter(engine)

	// 实例化http server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", HttpHost, HttpPort),
		Handler: engine,
	}

	fmt.Printf("http: successfully initialized [%s]\n", time.Since(start))

	// 启动http server
	go func() {
		fmt.Printf("http: starting web server at %s", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("http: web server shutdown complete")
			} else {
				fmt.Printf("http: web server closed unexpect: %s\n", err)
			}
		}
	}()

	// 关闭http server
	<-ctx.Done()
	fmt.Printf("http: shutting down web server")
	err := server.Close()
	if err != nil {
		fmt.Printf("http: web server shutdown failed: %v", err)
	}
}

func registerRouter(router *gin.Engine) {

}
