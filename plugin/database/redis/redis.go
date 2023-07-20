package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	rds "github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

const (
	defaultName = "default"
	pluginType  = "database"
	pluginName  = "redis"
)

var (
	Client         *rds.Client
	DefaultFactory = &Factory{}
	clients        = make(map[string]*rds.Client)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *rds.Client {
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
		return fmt.Errorf("redis config is emtpy")
	}
	for confName, config := range redisMap {
		cli, err := NewClient(config)
		if err != nil {
			return err
		}
		mu.Lock()
		if confName == defaultName {
			Client = cli
		}
		clients[confName] = cli
		mu.Unlock()
	}
	return nil
}

func NewClient(config *Config) (*rds.Client, error) {
	opts := &rds.Options{
		Addr: config.Addr,
	}

	if config.User != "" {
		opts.Username = config.User
	}
	if config.Pwd != "" {
		opts.Password = config.Pwd
	}
	if config.ConnMaxIdleTimeSec > 0 {
		opts.ConnMaxIdleTime = time.Second * time.Duration(config.ConnMaxIdleTimeSec)
	}
	if config.ConnMaxLifeTimeSec > 0 {
		opts.ConnMaxLifetime = time.Second * time.Duration(config.ConnMaxLifeTimeSec)
	}
	if config.MaxIdleConn > 0 {
		opts.MaxIdleConns = config.MaxIdleConn
	}
	if config.MinIdleConn > 0 {
		opts.MinIdleConns = config.MinIdleConn
	}
	if config.DialTimeoutMils > 0 {
		opts.DialTimeout = time.Millisecond * time.Duration(config.DialTimeoutMils)
	}
	cli := rds.NewClient(opts)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if config.PingTimeoutMils > 0 {
		ctx, _ = context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.PingTimeoutMils))
	}
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping err:%s", err)
	}
	return cli, nil
}
