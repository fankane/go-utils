package log

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const (
	pluginType     = "log"
	pluginName     = "zap_sugar"
	defaultLogName = "default"
)

var (
	Logger         *zap.SugaredLogger
	Logger2        *Log
	DefaultFactory = &Factory{}
	mu             sync.RWMutex
	loggers        = make(map[string]*zap.SugaredLogger)
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetLogger(name string) *zap.SugaredLogger {
	mu.RLock()
	defer mu.RUnlock()
	return loggers[name]
}

func GetLogger2(name string) *Log {
	mu.RLock()
	defer mu.RUnlock()
	return &Log{
		log: loggers[name],
	}
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	return newLogger(node)
}

func newLogger(node *yaml.Node) error {
	mu.Lock()
	defer mu.Unlock()
	logMap := make(map[string]*Config)
	if err := node.Decode(&logMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}

	if len(logMap) == 0 {
		return fmt.Errorf("log is empty")
	}
	for logName, conf := range logMap {
		core := zapcore.NewCore(getEncoder(conf), getLogWriter(conf), getLevel(conf.Level))
		logger := zap.New(core, zap.AddCaller())
		loggers[logName] = logger.Sugar()
		if logName == defaultLogName {
			Logger = logger.Sugar()
			Logger2 = &Log{
				log: Logger,
			}
		}
	}
	return nil
}

func getEncoder(conf *Config) zapcore.Encoder {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	if conf.EnableColor {
		encoderConf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
	if strings.ToLower(conf.Format) == FormatJSON {
		return zapcore.NewJSONEncoder(encoderConf)
	}
	return zapcore.NewConsoleEncoder(encoderConf)
}

func getLogWriter(conf *Config) zapcore.WriteSyncer {
	if conf.Filename == "" {
		conf.Filename = "./log.log"
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
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
