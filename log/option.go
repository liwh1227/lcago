package log

import (
	"path/filepath"
)

var sp = string(filepath.Separator)

type Options struct {
	// 开发模式
	Development bool
	// 日志存放路径
	LogFileDir string
	// 服务名称
	ServiceName string
	// 切分文件阈值 MB
	MaxSize int
	// 保留日志文件数量
	MaxBackups int
	// 日志文件存放时间 Day
	MaxAge int
	// 日志级别
	Level string
	// context
	CtxKey string
	// 是否写入文件
	WriteFile bool
	// 是否打印到控制台
	WriteConsole bool
	// 是否按日志级别输出至不同文件
	WriteWithLevel bool
	// error 日志文件
	ErrorFileName string
	// info 日志文件
	InfoFileName string
	// debug 日志文件
	DebugFileName string
	// warn 日志文件
	WarnFileName string
}

type LoggerOptions func(*Options)

func newOptions(opts ...LoggerOptions) *Options {
	opt := &Options{
		Development:    false,
		ServiceName:    DefaultAppName,
		MaxSize:        DefaultMaxSize,
		MaxBackups:     60,
		MaxAge:         30,
		Level:          "debug",
		CtxKey:         "hlog-ctx",
		WriteFile:      true,
		WriteConsole:   true,
		WriteWithLevel: false,
		ErrorFileName:  "error.log",
		WarnFileName:   "warn.log",
		InfoFileName:   "info.log",
		DebugFileName:  "debug.log",
	}
	var err error
	opt.LogFileDir, err = filepath.Abs(filepath.Dir(filepath.Join(".")))
	if err != nil {

	}
	opt.LogFileDir += sp + "logs" + sp
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func SetDevelopment(development bool) LoggerOptions {
	return func(options *Options) {
		options.Development = development
	}
}

func SetLogFileDir(logFileDir string) LoggerOptions {
	return func(options *Options) {
		options.LogFileDir = logFileDir
	}
}

func SetAppName(appName string) LoggerOptions {
	return func(options *Options) {
		options.ServiceName = appName
	}
}

func SetMaxSize(maxSize int) LoggerOptions {
	return func(options *Options) {
		options.MaxSize = maxSize
	}
}
func SetMaxBackups(maxBackups int) LoggerOptions {
	return func(options *Options) {
		options.MaxBackups = maxBackups
	}
}
func SetMaxAge(maxAge int) LoggerOptions {
	return func(options *Options) {
		options.MaxAge = maxAge
	}
}

func SetLevel(level string) LoggerOptions {
	return func(options *Options) {
		options.Level = level
	}
}

func SetCtxKey(ctxKey string) LoggerOptions {
	return func(options *Options) {
		options.CtxKey = ctxKey
	}
}

func SetWriteFile(writeFile bool) LoggerOptions {
	return func(options *Options) {
		options.WriteFile = writeFile
	}
}

func SetWriteConsole(writeConsole bool) LoggerOptions {
	return func(options *Options) {
		options.WriteConsole = writeConsole
	}
}
