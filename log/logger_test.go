package log

import (
	"testing"

	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func TestLog(t *testing.T) {
	NewLogger()
	Info("start init logger")
	//Error("test: error log.")
	Warn("test: warn log.")
	Debug("test: debug log.")
	m := map[string]string{
		"traceId": "dsadasdasdasdasda",
	}

	ctx, _ := AddCtx(context.Background(), m)
	Debug("TestGetLogger", zap.Any("t", "t"))

	WithContext(ctx).Info("[info]: test zap log info.")
	Info("hello world")

	FA(ctx)
}

func FA(ctx context.Context) {
	GetCtx(ctx).Warn("[warn]: test get context func")
	Warn("[warn]: no traceId log out put")
}

// TestWrapper
func TestWrapper(t *testing.T) {
	NewLogger()
	str := "locahost"
	port := 8080
	Infof("get something test %s ", str, port)
	getLogInstance().Infof("get %s", str)
}
