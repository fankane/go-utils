package rocketmq

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"gopkg.in/yaml.v3"
	"sync"
)

const (
	defaultName = "default"
	pluginType  = "queue"
	pluginName  = "rocketmq"
)

var (
	DefaultFactory      = &Factory{}
	DefaultProducer     *Producer
	mu                  = sync.RWMutex{}
	muConsumer          = sync.RWMutex{}
	globalProducers     = make(map[string]*Producer)
	globalConsumers     = make(map[string]*Consumer)
	globalConsumerConfs = make(map[string]*ConsumerConf)
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetProducer(name string) *Producer {
	mu.RLock()
	defer mu.RUnlock()
	return globalProducers[name]
}
func GetConsumer(name string) *Consumer {
	muConsumer.RLock()
	defer muConsumer.RUnlock()
	return globalConsumers[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	nsqConf := &Config{}
	if err := node.Decode(&nsqConf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(nsqConf.Consumers) == 0 && len(nsqConf.Producers) == 0 {
		return fmt.Errorf("rocketmq config is emtpy")
	}
	if err := initProducers(nsqConf.Producers); err != nil {
		return err
	}
	globalConsumerConfs = nsqConf.Consumers
	return nil
}

func initProducers(confM map[string]*ProducerConf) error {
	for name, conf := range confM {
		p, err := NewProducer(conf)
		if err != nil {
			return fmt.Errorf("new producer err:%s, name:%s", err, name)
		}
		if err = p.Start(); err != nil {
			return fmt.Errorf("%s start producer failed %s", name, err)
		}
		if name == defaultName {
			DefaultProducer = p
		}
		globalProducers[name] = p
	}
	return nil
}
