package nacos

import (
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	pluginType  = "serve"
	pluginName  = "nacos"
	defaultName = "default"
)

var (
	DefaultFactory = &Factory{}
	confVal        = &Config{}
	Cli            *Client
	Clis           = make(map[string]*Client)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *Client {
	mu.RLock()
	defer mu.RUnlock()
	return Clis[name]
}

type Factory struct {
}

// Type 插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载配置
func (f *Factory) Setup(name string, node *yaml.Node) error {
	nacosMap := make(map[string]*Config)
	if err := node.Decode(&nacosMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	for k, config := range nacosMap {
		cli, err := NewClient(config)
		if err != nil {
			return fmt.Errorf("newClient err:%s", err)
		}
		Clis[k] = cli
		if k == defaultName {
			Cli = cli
		}
	}
	return viper.ReadInConfig()
}

func NewClient(conf *Config) (*Client, error) {
	clientConfig := &constant.ClientConfig{
		Endpoint:            conf.EndPoint,
		NamespaceId:         conf.NamespaceID,
		Username:            conf.Username,
		Password:            conf.Password,
		TimeoutMs:           conf.TimeoutMs,
		NotLoadCacheAtStart: conf.NotLoadCacheAtStart,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: clientConfig,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Client{Cli: client, conf: conf}, nil
}

func (c *Client) ParseConfig(f MarshalFunc, v interface{}) error {
	content, err := c.GetConfig()
	if err != nil {
		return err
	}
	return f(content, v)
}

func (c *Client) ParseListenConfig(f MarshalFunc, v interface{}) error {
	if err := c.ParseConfig(f, v); err != nil {
		return err
	}
	return c.ListenConfig(f, v)
}

func (c *Client) GetConfig() (string, error) {
	if err := c.checkClient(); err != nil {
		return "", err
	}
	return c.Cli.GetConfig(vo.ConfigParam{
		DataId: c.conf.DataId,
		Group:  c.conf.Group,
	})
}

func (c *Client) ListenConfig(f MarshalFunc, v interface{}) error {
	if err := c.checkClient(); err != nil {
		return err
	}
	return c.Cli.ListenConfig(vo.ConfigParam{
		DataId: c.conf.DataId,
		Group:  c.conf.Group,
		OnChange: func(namespace, group, dataId, data string) {
			f(data, v)
		},
	})
}

func (c *Client) checkClient() error {
	if c == nil || c.Cli == nil || c.conf == nil {
		return fmt.Errorf("client invalid")
	}
	return nil
}
