package es

import (
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v3"
)

const (
	defaultName = "default"
	pluginType  = "database"
	pluginName  = "elasticsearch"
)

var (
	Cli            *Client
	DefaultFactory = &Factory{}
	clients        = make(map[string]*Client)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *Client {
	mu.RLock()
	defer mu.RUnlock()
	return clients[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	redisMap := make(map[string]*Config)
	if err := node.Decode(&redisMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(redisMap) == 0 {
		return fmt.Errorf("es config is emtpy")
	}
	for confName, config := range redisMap {
		cli, err := NewClient(config)
		if err != nil {
			return err
		}
		mu.Lock()
		if confName == defaultName {
			Cli = cli
		}
		clients[confName] = cli
		mu.Unlock()
	}
	return nil
}

func NewClient(config *Config) (*Client, error) {
	cli, err := elastic.NewSimpleClient(elastic.SetURL(config.Addr...),
		elastic.SetBasicAuth(config.User, config.Pwd))
	if err != nil {
		return nil, err
	}
	return &Client{
		ESCli: cli,
	}, nil
}

type Client struct {
	ESCli *elastic.Client
}
