package plugin

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fankane/go-utils/goroutine"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	Plugins Config
}

// Config 插件统一配置 plugin type => { plugin name => plugin config } 。
type Config map[string]map[string]yaml.Node

type Factory interface {
	// Type 插件的类型 如 selector log config tracing
	Type() string
	// Setup 根据配置项节点装载插件，需要用户自己先定义好具体插件的配置数据结构
	Setup(name string, node *yaml.Node) error
}

var (
	mu      = sync.RWMutex{}
	plugins = make(map[string]map[string]Factory) // plugin type => { plugin name => plugin factory }
)

type LoadParam struct {
	ConfigFile string
	IgnoreErr  bool //忽略插件加载失败场景，默认false
}

type Option func(param *LoadParam)

func ConfigFile(file string) Option {
	return func(param *LoadParam) {
		param.ConfigFile = file
	}
}

func IgnoreErr(ignore bool) Option {
	return func(param *LoadParam) {
		param.IgnoreErr = ignore
	}
}

func Load(opts ...Option) error {
	params := &LoadParam{
		ConfigFile: "system_plugin.yaml", //默认文件路径
	}
	for _, opt := range opts {
		opt(params)
	}
	// 默认读取 system_plugin.yaml 文件，来加载配置
	res, err := os.ReadFile(params.ConfigFile)
	if err != nil {
		return fmt.Errorf("read plugin config file err:%s, filepath:%s", err, params.ConfigFile)
	}
	pluginConf := &GlobalConfig{}
	if err = yaml.Unmarshal(res, &pluginConf); err != nil {
		return fmt.Errorf("yaml unmarshal err:%s", err)
	}
	if pluginConf == nil || len(pluginConf.Plugins) == 0 {
		return fmt.Errorf("plugin is empty")
	}
	//if err = pluginConf.Plugins.Setup(params.IgnoreErr); err != nil {
	if err = InitPlugins(pluginConf.Plugins, params.IgnoreErr); err != nil {
		return fmt.Errorf("setup err:%s", err)
	}
	return nil
}

func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	factories, ok := plugins[factory.Type()]
	if !ok {
		factories = make(map[string]Factory)
		plugins[factory.Type()] = factories
	}
	factories[name] = factory
}

func (c Config) Setup(ignoreErr bool) error {
	fs := make([]func() error, 0)
	for typT, factories := range c {
		for pluginNameT, confT := range factories {
			typ, pluginName, conf := typT, pluginNameT, confT
			if strings.Contains(pluginName, "-") {
				return fmt.Errorf("pluginName:%s contain forbbiden char \"-\"", pluginName)
			}
			fs = append(fs, func() error {
				f := Get(typ, pluginName)
				if f == nil {
					return fmt.Errorf("[%s - %s] not register", typ, pluginName)
				}
				err := f.Setup(pluginName, &conf)
				if err != nil {
					if ignoreErr {
						log.Println(fmt.Sprintf("%s setup failed, err:%s", pluginName, err))
						return nil
					}
					return err
				}
				log.Println(fmt.Sprintf("%s:%s installed ", typ, pluginName))
				return nil
			})
		}
	}
	return goroutine.Exec(fs, goroutine.WithReturnWhenError(true))
}

// Get 根据插件类型，插件名字获取插件工厂。
func Get(typ string, name string) Factory {
	fMap, ok := plugins[typ]
	if !ok {
		return nil
	}
	return fMap[name]
}
