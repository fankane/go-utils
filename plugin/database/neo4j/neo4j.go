package neo4j

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gopkg.in/yaml.v3"
)

const (
	defaultCliName = "default"
	pluginType     = "database"
	pluginName     = "neo4j"
)

var (
	Dri            *Client
	DefaultFactory = &Factory{}
	Dris           = make(map[string]*Client)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *Client {
	mu.RLock()
	defer mu.RUnlock()
	return Dris[name]
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
		cli, err := NewClient(context.Background(), config)
		if err != nil {
			return err
		}
		if confName == defaultCliName {
			Dri = cli
		}
		mu.Lock()
		Dris[confName] = cli
		mu.Unlock()
	}
	return nil
}

type Client struct {
	DriverCtx  neo4j.DriverWithContext
	SessionCtx neo4j.SessionWithContext
	conf       *Config
}

func NewClient(ctx context.Context, conf *Config) (*Client, error) {
	if conf.AccessMode != int(neo4j.AccessModeWrite) && conf.AccessMode != int(neo4j.AccessModeRead) {
		conf.AccessMode = int(neo4j.AccessModeWrite) //默认为可写
	}
	driver, err := neo4j.NewDriverWithContext(conf.Target, neo4j.BasicAuth(conf.User, conf.Pwd, conf.Realm))
	if err != nil {
		return nil, err
	}
	if err = driver.VerifyConnectivity(ctx); err != nil {
		return nil, fmt.Errorf("connect neo4j failed err:%s", err)
	}
	sess := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessMode(conf.AccessMode),
		DatabaseName: conf.DatabaseName})
	if err != nil {
		return nil, err
	}
	return &Client{DriverCtx: driver, SessionCtx: sess, conf: conf}, nil
}

func (d *Client) Run(ctx context.Context, cql string, params map[string]interface{}) ([]*neo4j.Record, error) {
	if d.SessionCtx == nil {
		d.SessionCtx = d.DriverCtx.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessMode(d.conf.AccessMode),
			DatabaseName: d.conf.DatabaseName})
	}

	resultCtx, err := d.SessionCtx.Run(ctx, cql, params)
	if err != nil {
		if neo4j.IsUsageError(err) { //session断开，置空，下次cql执行，创建新的Session
			d.SessionCtx = nil
		}
		return nil, err
	}
	return resultCtx.Collect(ctx)
}

type Driver struct {
	DriverCtx neo4j.DriverWithContext
	conf      *Config
}

func NewDriver(ctx context.Context, conf *Config) (*Driver, error) {
	if conf.AccessMode != int(neo4j.AccessModeWrite) && conf.AccessMode != int(neo4j.AccessModeRead) {
		conf.AccessMode = int(neo4j.AccessModeWrite) //默认为可写
	}
	driver, err := neo4j.NewDriverWithContext(conf.Target, neo4j.BasicAuth(conf.User, conf.Pwd, conf.Realm))
	if err != nil {
		return nil, err
	}
	if err = driver.VerifyConnectivity(ctx); err != nil {
		return nil, fmt.Errorf("connect neo4j failed err:%s", err)
	}
	return &Driver{DriverCtx: driver, conf: conf}, nil
}

func (d *Driver) Close() error {
	if d == nil || d.DriverCtx == nil {
		return errors.New("empty driver")
	}
	return d.DriverCtx.Close(context.Background())
}

func (d *Driver) NewSession(ctx context.Context) (*Session, error) {
	sess := d.DriverCtx.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessMode(d.conf.AccessMode),
		DatabaseName: d.conf.DatabaseName})
	return &Session{
		SessionCtx: sess,
	}, nil
}

type Session struct {
	SessionCtx neo4j.SessionWithContext
}

func (s *Session) Close() error {
	if s == nil || s.SessionCtx == nil {
		return nil
	}
	return s.SessionCtx.Close(context.Background())
}

func (s *Session) Run(ctx context.Context, cql string, params map[string]interface{}) ([]*neo4j.Record, error) {
	resultCtx, err := s.SessionCtx.Run(ctx, cql, params)
	if err != nil {
		return nil, err
	}
	return resultCtx.Collect(ctx)
}
