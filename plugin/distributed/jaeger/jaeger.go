package jaeger

import (
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"gopkg.in/yaml.v3"
)

const (
	defaultName = "default"
	pluginType  = "distributed"
	pluginName  = "jaeger"
)

var (
	DefaultFactory = &Factory{}
	mu             = sync.RWMutex{}
	Tracer         *TraceClient
	tracerMap      = make(map[string]*TraceClient)
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetTracer(name string) *TraceClient {
	mu.RLock()
	defer mu.RUnlock()
	return tracerMap[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	jaegerMap := make(map[string]*Config)
	if err := node.Decode(&jaegerMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(jaegerMap) == 0 {
		return fmt.Errorf("jaeger config is emtpy")
	}
	for confName, config := range jaegerMap {
		t, err := NewTracer(config)
		if err != nil {
			return err
		}
		if confName == defaultName {
			Tracer = t
		}
		mu.Lock()
		tracerMap[confName] = t
		mu.Unlock()
	}
	return nil
}

func NewTracer(conf *Config) (*TraceClient, error) {
	fmt.Println(str.ToJSON(conf))
	interval := conf.BufferFlushIntervalMs
	if interval == 0 {
		interval = defaultIntervalMs
	}
	cfg := jaegercfg.Configuration{
		ServiceName: conf.ServiceName, // 对其发起请求的的调用链，叫什么服务
		Sampler: &jaegercfg.SamplerConfig{
			Type:  conf.SamplerType,
			Param: conf.SamplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            conf.LogSpans,
			CollectorEndpoint:   conf.CollectorEndpoint,
			User:                conf.User,
			Password:            conf.Password,
			BufferFlushInterval: time.Millisecond * time.Duration(interval),
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}
	return &TraceClient{Tracer: tracer, Closer: closer}, nil
}
