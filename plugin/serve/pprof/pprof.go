package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/fankane/go-utils/plugin"
	"gopkg.in/yaml.v3"
)

const (
	pluginType = "serve"
	pluginName = "pprof"
)

var (
	DefaultFactory = &Factory{}
)

type Config struct {
	Addr string `yaml:"addr"`
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
	if err := node.Decode(conf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	go func() {
		http.ListenAndServe(conf.Addr, nil)
	}()
	return nil
}
