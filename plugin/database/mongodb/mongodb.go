package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

const (
	defaultMogoName = "default"
	pluginType      = "database"
	pluginName      = "mongo"
)

var (
	Cli            *MongoCli
	DefaultFactory = &Factory{}
	dbs            = make(map[string]*MongoCli)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *MongoCli {
	mu.RLock()
	defer mu.RUnlock()
	return dbs[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	mongoMap := make(map[string]*Config)
	if err := node.Decode(&mongoMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(mongoMap) == 0 {
		return fmt.Errorf("mongo config is emtpy")
	}
	for confName, config := range mongoMap {
		cli, err := NewClient(config)
		if err != nil {
			return err
		}
		if confName == defaultMogoName {
			Cli = cli
		}
		mu.Lock()
		dbs[confName] = cli
		mu.Unlock()
	}
	return nil
}

func NewClient(config *Config) (*MongoCli, error) {
	if config.ConnectTimeoutMs <= 0 {
		config.ConnectTimeoutMs = DefaultConnTimeoutMs
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI(config)),
		options.Client().SetConnectTimeout(time.Duration(config.ConnectTimeoutMs)*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return &MongoCli{
		Cli: client,
	}, nil
}

func mongoURI(config *Config) string {
	if config.User != "" || config.Pwd != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Pwd, config.Host, config.Port)
	}
	return fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)
}
