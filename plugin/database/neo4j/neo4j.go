package neo4j

import (
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"gopkg.in/yaml.v3"
)

const (
	defaultCliName = "default"
	pluginType     = "database"
	pluginName     = "neo4j"
)

var (
	Cli            *Client
	DefaultFactory = &Factory{}
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

// Type 日志插件类型
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
		return fmt.Errorf("neo4j config is emtpy")
	}
	for confName, config := range confMap {
		cli, err := NewSession(config)
		if err != nil {
			return err
		}
		if confName == defaultCliName {
			Cli = cli
		}
		mu.Lock()
		Clis[confName] = cli
		mu.Unlock()
	}
	return nil
}

type Client struct {
	Session neo4j.Session
}

func NewSession(conf *Config) (*Client, error) {
	if conf.AccessMode != int(neo4j.AccessModeWrite) && conf.AccessMode != int(neo4j.AccessModeRead) {
		conf.AccessMode = int(neo4j.AccessModeWrite) //默认为可写
	}
	driver, err := neo4j.NewDriver(conf.Target, neo4j.BasicAuth(conf.User, conf.Pwd, conf.Realm))
	if err != nil {
		return nil, err
	}
	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("connect neo4j failed err:%s", err)
	}
	sess := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessMode(conf.AccessMode), DatabaseName: conf.DatabaseName})
	return &Client{
		Session: sess,
	}, nil
}
