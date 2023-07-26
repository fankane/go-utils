package conf

import (
	"fmt"
	"github.com/fankane/go-utils/utime"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"strings"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/slice"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	pluginType = "serve"
	pluginName = "conf"
)

var (
	DefaultFactory = &Factory{}
	confVal        = &Config{}
)

var validType = []string{
	"yaml",
	"toml",
}

type Config struct {
	ConfFile    string `yaml:"conf_file"`    // 配置文件
	WatchChange bool   `yaml:"watch_change"` // 监听文件更新, 默认false
	ChangeCron  string `yaml:"change_cron"`  // 更新配置的执行频率, 当 watch_change = true时 生效
}

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	conf := &Config{}
	if err := node.Decode(conf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	confVal = conf

	filePath, fileName := filepath.Dir(conf.ConfFile), filepath.Base(conf.ConfFile)
	configInfos := strings.Split(fileName, ".")
	if len(configInfos) != 2 {
		return fmt.Errorf("confFile:[%s] invalid", fileName)
	}
	if !slice.InStrings(configInfos[1], validType) {
		return fmt.Errorf("invalied config type, only support:%v", validType)
	}
	viper.AddConfigPath(filePath)
	viper.SetConfigName(configInfos[0])
	viper.SetConfigType(configInfos[1])
	viper.ReadInConfig()
	return nil
}

func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if confVal.WatchChange {
		viper.OnConfigChange(func(in fsnotify.Event) {
			if !in.Has(fsnotify.Write) {
				return
			}
			if confVal.ChangeCron != "" {
				utime.CronDo(confVal.ChangeCron, func() {
					viper.Unmarshal(rawVal, opts...)
				})
			}
		})
		viper.WatchConfig()
	}
	return viper.Unmarshal(rawVal, opts...)
}
