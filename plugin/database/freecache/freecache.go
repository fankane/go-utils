package freecache

import (
	"fmt"
	"sync"

	"github.com/coocood/freecache"
	"github.com/fankane/go-utils/plugin"
	"gopkg.in/yaml.v3"
)

const (
	defaultName = "default"
	pluginType  = "database"
	pluginName  = "freecache"
)

var (
	Cache          *freecache.Cache
	DefaultFactory = &Factory{}
	caches         = make(map[string]*freecache.Cache)
	mu             = sync.RWMutex{}
)

type Config struct {
	CacheSize int `yaml:"cache_size"`
}

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetCache(name string) *freecache.Cache {
	mu.RLock()
	defer mu.RUnlock()
	return caches[name]
}

type Factory struct {
}

// Type 插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	confMap := make(map[string]*Config)
	if err := node.Decode(&confMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(confMap) == 0 {
		return fmt.Errorf("cache config is emtpy")
	}
	for confName, config := range confMap {
		cache := freecache.NewCache(config.CacheSize)
		if confName == defaultName {
			Cache = cache
		}
		mu.Lock()
		caches[confName] = cache
		mu.Unlock()
	}
	return nil
}
