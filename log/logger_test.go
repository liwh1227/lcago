package log

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"testing"
)

func TestLog(t *testing.T) {
	NewLogger()
	Log().Info("start init logger")

	ctx, log := Log().AddCtx(context.Background(), zap.String("traceId", uuid.New().String()))
	log.Debug("TestGetLogger", zap.Any("t", "t"))

	FA(ctx)
}

func FA(ctx context.Context) {
	hlog := Log().GetCtx(ctx)
	hlog.Info("FA", zap.Any("a", "a"))
}
