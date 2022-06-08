package ginx

import (
	"fmt"
	"lcago/ginx/handler/lcago"
	"lcago/ginx/router"
	"lcago/log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// Run 启动api服务
func Run(ctx context.Context) {
	// panic 自动恢复
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			return
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
	router.Register(lcago.Router)
	engine := router.Init()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", HttpHost, HttpPort),
		Handler: engine,
	}

	log.Infof("http: successfully initialized [%s]", time.Since(start))

	// 启动http server
	go func() {
		log.Infof("http: starting web server at %s", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Error("http: web server shutdown complete")
			} else {
				log.Errorf("http: web server closed unexpect: %s\n", err)
			}
		}
	}()

	// 关闭http server
	<-ctx.Done()
	log.Info("http: shutting down web server")
	err := server.Close()
	if err != nil {
		log.Errorf("http: web server shutdown failed: %v", err)
	}
}
