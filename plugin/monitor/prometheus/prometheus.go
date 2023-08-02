package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

const (
	pluginType = "monitor"
	pluginName = "prometheus"
)

var (
	defaultPath    = "/metrics"
	DefaultFactory = &Factory{}
	mu             sync.RWMutex
)

type Config struct {
	Port int    `yaml:"port" validate:"gt=0,lte=65535"`
	Path string `yaml:"path"`
}

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	conf := &Config{}
	if err := node.Decode(&conf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("prometheus conf invalided err:%s", err)
	}
	path := conf.Path
	if strings.TrimSpace(path) == "" {
		path = defaultPath
	}
	http.Handle(defaultPath, promhttp.Handler())

	var err error
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	}()
	time.Sleep(time.Millisecond * 3) //等待3毫秒，让 ListenAndServe 完成, 如果失败，3ms足够返回error
	if err != nil {
		return fmt.Errorf("prometheus listen err:%s", err)
	}
	return nil
}

func NewGauge() {
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_hf",
	}, []string{"label1", "label2"})
	g.WithLabelValues("val1", "val2").Set(1.0)
	g.WithLabelValues("val1", "val3").Set(2.0)
	g.WithLabelValues("val2", "val3").Set(3.0)
	prometheus.MustRegister(g)
}
