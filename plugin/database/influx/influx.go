package influx

import (
	"context"
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"gopkg.in/yaml.v3"
)

const (
	defaultInfluxName = "default"
	pluginType        = "database"
	pluginName        = "influx"
)

var (
	Cli            *Client
	DefaultFactory = &Factory{}
	clis           = make(map[string]*Client)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetClient(name string) *Client {
	mu.RLock()
	defer mu.RUnlock()
	return clis[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	influxMap := make(map[string]*Config)
	if err := node.Decode(&influxMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(influxMap) == 0 {
		return fmt.Errorf("influx config is emtpy")
	}
	for confName, config := range influxMap {
		cli, err := NewClient(config)
		if err != nil {
			return err
		}
		if confName == defaultInfluxName {
			Cli = cli
		}
		mu.Lock()
		clis[confName] = cli
		mu.Unlock()
	}
	return nil
}

type Client struct {
	InfCli   influxdb2.Client
	QueryCli api.QueryAPI

	writerBlocking api.WriteAPIBlocking
	writer         api.WriteAPI
	conf           *Config
}

func NewClient(conf *Config) (*Client, error) {
	client := influxdb2.NewClient(conf.URL, conf.Token)
	result := &Client{InfCli: client, conf: conf}
	if !conf.WriteAsync {
		result.writerBlocking = client.WriteAPIBlocking(conf.Org, conf.Bucket)
	} else {
		result.writer = client.WriteAPI(conf.Org, conf.Bucket)
	}
	return result, nil
}

func (c *Client) WritePoint(ctx context.Context, point ...*write.Point) error {
	if !c.conf.WriteAsync {
		return c.writerBlocking.WritePoint(ctx, point...)
	} else {
		for _, p := range point {
			c.writer.WritePoint(p)
		}
		return nil
	}
}

func (c *Client) WriteRecord(ctx context.Context, line ...string) error {
	if !c.conf.WriteAsync {
		return c.writerBlocking.WriteRecord(ctx, line...)
	} else {
		for _, l := range line {
			c.writer.WriteRecord(l)
		}
		return nil
	}
}

func (c *Client) Flush() {
	if !c.conf.WriteAsync {
		c.writerBlocking.Flush(context.Background())
		return
	} else {
		c.writer.Flush()
		return
	}
}

func (c *Client) QueryClient() api.QueryAPI {
	if c.QueryCli != nil {
		return c.QueryCli
	}
	queryAPI := c.InfCli.QueryAPI(c.conf.Org)
	c.QueryCli = queryAPI
	return queryAPI
}
