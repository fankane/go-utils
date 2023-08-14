package kafka

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/fankane/go-utils/plugin"
	"gopkg.in/yaml.v3"
)

const (
	defaultName = "default"
	pluginType  = "queue"
	pluginName  = "kafka"
)

var (
	DefaultFactory    = &Factory{}
	DefaultProducer   Producer
	globalProducers   = make(map[string]Producer)
	mu                = sync.RWMutex{}
	globalConsumerMap map[string]*ConsumerConf
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetProducer(name string) Producer {
	mu.RLock()
	defer mu.RUnlock()
	return globalProducers[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	kafkaConf := &Config{}
	if err := node.Decode(&kafkaConf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(kafkaConf.Consumers) == 0 && len(kafkaConf.Producers) == 0 {
		return fmt.Errorf("kafka config is emtpy")
	}
	if err := initProducers(kafkaConf.Producers); err != nil {
		return err
	}
	globalConsumerMap = kafkaConf.Consumers
	return nil
}

func initProducers(producers map[string]*ProducerConf) error {
	if len(producers) == 0 {
		return nil
	}
	for name, proConf := range producers {
		var (
			tempPro Producer
			err     error
		)
		if proConf.SendType == sendTypeSync {
			tempPro, err = NewSyncProducer(proConf)
			if err != nil {
				return err
			}
			log.Println(fmt.Sprintf("kafka producer-sync [%s] init success", name))
		} else if proConf.SendType == sendTypeAsync {
			tempPro, err = NewAsyncProducer(proConf)
			if err != nil {
				return err
			}
			log.Println(fmt.Sprintf("kafka producer-async [%s] init success", name))
		} else {
			return fmt.Errorf("unknown send_type:%s", proConf.SendType)
		}
		globalProducers[name] = tempPro
		if name == defaultName {
			DefaultProducer = tempPro
		}
	}
	return nil
}

func getDefaultConf() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Version = sarama.V1_1_1_0
	return conf
}
