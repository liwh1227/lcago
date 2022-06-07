package log

import (
	"context"
	//"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	l *Logger
	// IO输出
	outWrite zapcore.WriteSyncer
	// 控制台标准输出
	debugConsole = zapcore.Lock(os.Stdout)
	once         sync.Once
)

type Logger struct {
	*zap.Logger
	opts      *Options
	zapConfig zap.Config
}

func NewLogger(opts ...LoggerOptions) {
	once.Do(func() {
		l = &Logger{
			opts: newOptions(opts...),
		}
		l.loadCfg()
		l.initLog()
		l.Info("[initLogger] zap plugin initializing completed")
	})
}

// Log 封装的日志服务
func Log() *Logger {
	if l == nil {
		panic("Please initialize the log service first")
		return nil
	}
	return l
}

func (l *Logger) GetCtx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

// AddCtx 添加上下文
func (l *Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	log := l.With(field...)
	ctx = context.WithValue(ctx, l.opts.CtxKey, log)
	return ctx, log
}

// initLog 初始化日志
func (l *Logger) initLog() {
	l.setSyncers()
	var err error
	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

// GetLevel 获取日志级别
func (l *Logger) GetLevel() (level zapcore.Level) {
	switch strings.ToLower(l.opts.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel //默认为调试模式
	}
}

func (l *Logger) loadCfg() {
	if l.opts.Development {
		//l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		//l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
}

// setSyncers 同步日志设置
func (l *Logger) setSyncers() {
	outWrite = zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.opts.LogFileDir + "/" + l.opts.ServiceName + ".log",
		MaxSize:    l.opts.MaxSize,
		MaxBackups: l.opts.MaxBackups,
		MaxAge:     l.opts.MaxAge,
		Compress:   true,
		LocalTime:  true,
	})
	return
}

func (l *Logger) cores() zap.Option {
	encoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)

	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l.GetLevel()
	})
	var cores []zapcore.Core
	if l.opts.WriteFile {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(encoder, outWrite, priority),
		}...)
	}
	if l.opts.WriteConsole {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(encoder, debugConsole, priority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

// timeEncoder 可自定义时间
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// timeUnixNano unix时间格式
func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}
