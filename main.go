package main

import (
	"fmt"
	"lcago/ginx"
	"lcago/log"
	"lcago/utils"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

func main() {
	// 必要环境检查
	err := preStartService()
	if err != nil {
		os.Exit(1)
	}

	// 初始化配置

	// 初始化日志
	log.NewLogger()

	cctx, cancel := context.WithCancel(context.Background())
	// 启动http服务
	go ginx.Run(cctx)

	// 若接收到推出信号,主动退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	// 发送退出信号
	cancel()
}

// preStartService 服务启动前必要环境检查
func preStartService() error {
	exists, err := utils.PathExists("logs")
	if err != nil {
		fmt.Println("should create logs dir")
		return err
	}

	if !exists {
		return os.Mkdir("logs", 0666)
	}

	return nil
}
