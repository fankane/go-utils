package log

import (
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

func init() {
	// 默认读取system.yaml 文件，来加载 log 配置
	res, err := ioutil.ReadFile("system.yaml")
	if err != nil {
		return
	}
	conf := &ConfContent{}
	if err = yaml.Unmarshal(res, conf); err != nil {
		return
	}
	newLogger(conf.Log)
}

var Logger *zap.SugaredLogger

func newLogger(conf *Config) {
	core := zapcore.NewCore(getEncoder(), getLogWriter(conf), getLevel(conf.Level))
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderConf)
}

func getLogWriter(conf *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.file.Filename,
		MaxSize:    conf.file.MaxSize,
		MaxBackups: conf.file.MaxBackups,
		MaxAge:     conf.file.MaxAge,
		Compress:   conf.file.Compress,
	}
	if conf.EnableStdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	}
	return zapcore.InfoLevel
}
