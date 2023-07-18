package plugin

import (
	"fmt"
	"io/ioutil"
	"sync"

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

func Load(configFile string) error {
	if configFile == "" {
		configFile = "system_plugin.yaml"
	}
	// 默认读取system.yaml 文件，来加载 log 配置
	res, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read file err:%s", err)
	}
	pluginConf := &GlobalConfig{}
	if err = yaml.Unmarshal(res, &pluginConf); err != nil {
		return fmt.Errorf("yaml unmarshal err:%s", err)
	}
	if pluginConf == nil || len(pluginConf.Plugins) == 0 {
		return fmt.Errorf("plugin is empty")
	}
	if err = pluginConf.Plugins.Setup(); err != nil {
		return fmt.Errorf("setup err:%s", err)
	}
	return nil
}

func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	factories, ok := plugins[factory.Type()]
	if !ok {
		fmt.Println(fmt.Sprintf("%s not exists,create", factory.Type()))
		factories = make(map[string]Factory)
		plugins[factory.Type()] = factories
	}
	factories[name] = factory
	fmt.Println("register:", fmt.Sprintf("%+v", factory))
}

func (c Config) Setup() error {
	for typ, factories := range c {
		for pluginName, conf := range factories {
			f := Get(typ, pluginName)
			if f == nil {
				fmt.Println(fmt.Sprintf("global:%+v", plugins))
				return fmt.Errorf("[%s - %s] not register", typ, pluginName)
			}
			fmt.Println("type:", typ, ", name:", pluginName, ", node:", fmt.Sprintf("%+v", conf))
			if err := f.Setup(pluginName, &conf); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get 根据插件类型，插件名字获取插件工厂。
func Get(typ string, name string) Factory {
	return plugins[typ][name]
}
