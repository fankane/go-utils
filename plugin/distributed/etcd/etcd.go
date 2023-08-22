package etcd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

const (
	defaultEtcdName  = "default"
	pluginType       = "distributed"
	pluginName       = "etcd"
	defaultTimeoutMs = 5000
)

var (
	DefaultFactory = &Factory{}
	mu             = sync.RWMutex{}
	Op             Operate
	opMap          = make(map[string]Operate)
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetOperate(name string) Operate {
	mu.RLock()
	defer mu.RUnlock()
	return opMap[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	etcdMap := make(map[string]*Config)
	if err := node.Decode(&etcdMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(etcdMap) == 0 {
		return fmt.Errorf("etcd config is emtpy")
	}
	for confName, config := range etcdMap {
		op, err := NewEtcdCli(confName, config)
		if err != nil {
			return err
		}
		if confName == defaultEtcdName {
			Op = op
		}
		mu.Lock()
		opMap[confName] = op
		mu.Unlock()
	}
	return nil
}

type Operate interface {
	GetClient() *clientv3.Client //获取原生etcd client
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (map[string]*EValue, error)
	Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (int64, error)
	Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (int64, error)
	Watch(ctx context.Context, key string, h HandleFunc, opts ...clientv3.OpOption)
	RegisterServer() error
	UnRegisterServer() error
	GetServers() map[string]*ServerInfo
}

type etcd struct {
	cli      *clientv3.Client
	conf     *Config
	confName string
}

type EValue struct {
	Val     []byte
	Version int64
	Lease   int64
}

type HandleFunc interface {
	PutHandle(key, value []byte, version int64)
	DelHandle(key []byte)
}

func NewEtcdCli(confName string, conf *Config) (Operate, error) {
	if conf.DialTimeoutMs <= 0 {
		conf.DialTimeoutMs = defaultTimeoutMs
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: time.Millisecond * time.Duration(conf.DialTimeoutMs),
		Username:    conf.Username,
		Password:    conf.Password,
	})
	if err != nil {
		return nil, err
	}
	result := &etcd{cli: cli, conf: conf, confName: confName}
	if conf.OpenDiscovery {
		if err = result.RegisterServer(); err != nil {
			return nil, fmt.Errorf("register server failed err:%s", err)
		}
	}
	return result, nil
}

func (e *etcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (map[string]*EValue, error) {
	resp, err := e.cli.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*EValue)
	for _, kv := range resp.Kvs {
		result[str.FromBytes(kv.Key)] = &EValue{
			Val:     kv.Value,
			Version: kv.Version,
			Lease:   kv.Lease,
		}
	}
	return result, nil
}

func (e *etcd) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (int64, error) {
	res, err := e.cli.Put(ctx, key, value, opts...)
	if err != nil {
		return 0, err
	}
	return res.Header.Revision, nil
}

func (e *etcd) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (int64, error) {
	resp, err := e.cli.Delete(ctx, key, opts...)
	if err != nil {
		return 0, err
	}
	return resp.Deleted, nil
}

func (e *etcd) Watch(ctx context.Context, key string, h HandleFunc, opts ...clientv3.OpOption) {
	watchChan := e.cli.Watch(ctx, key, opts...)
	for {
		c := <-watchChan
		if c.Canceled {
			return
		}
		for _, event := range c.Events {
			if event.Kv == nil {
				continue
			}
			switch event.Type {
			case clientv3.EventTypePut:
				h.PutHandle(event.Kv.Key, event.Kv.Value, event.Kv.Version)
			case clientv3.EventTypeDelete:
				h.DelHandle(event.Kv.Key)
			default:
				log.Printf("unknown event type:%s", event.Type)
			}
		}
	}
}

func (e *etcd) GetClient() *clientv3.Client {
	return e.cli
}
