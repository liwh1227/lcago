package log

import (
	"os"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	l                              *Logger
	outWrite                       zapcore.WriteSyncer
	errWs, warnWs, infoWs, debugWs zapcore.WriteSyncer
	once                           sync.Once
	debugConsole                   = zapcore.Lock(os.Stdout) // 控制台标准输出
)

type Logger struct {
	*zap.SugaredLogger
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

// loadCfg 加载zap的配置项
func (l *Logger) loadCfg() {
	// 1.开发模式
	l.loadDevModeCfg()
}

// loadDevModeDev 加载当前的开发环境
func (l *Logger) loadDevModeCfg() {
	if l.opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	}
}

// getLogInstance 封装的日志服务
func getLogInstance() *Logger {
	if l == nil {
		panic("Please initialize the log service first")
		return nil
	}
	return l
}

// initLog 初始化日志
func (l *Logger) initLog() {
	l.setSyncers()
	logger, err := l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	l.SugaredLogger = logger.Sugar()
	defer l.SugaredLogger.Sync()
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

// setSyncers 同步日志设置
func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		logf, _ := rotatelogs.New(l.opts.LogFileDir+l.opts.ServiceName+"-"+fN+".%Y-%m-%d-%H"+".log",
			rotatelogs.WithLinkName(l.opts.LogFileDir+l.opts.ServiceName+"-"+fN),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(time.Minute),
		)
		return zapcore.AddSync(logf)
	}

	errWs = f(l.opts.ErrorFileName)
	warnWs = f(l.opts.WarnFileName)
	infoWs = f(l.opts.InfoFileName)
	debugWs = f(l.opts.DebugFileName)
	outWrite = f("app")

	return
}

// cores zap日志切割
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
	// 是否将日志按不同的级别输出至不同日志
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})

	// 是否开启按日志级别进行切割
	if l.opts.WriteWithLevel {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(encoder, errWs, errPriority),
			zapcore.NewCore(encoder, warnWs, warnPriority),
			zapcore.NewCore(encoder, infoWs, infoPriority),
			zapcore.NewCore(encoder, debugWs, debugPriority),
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
