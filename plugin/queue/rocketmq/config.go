package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"time"
)

const (
	ConsumeFromLastOffset  = int(consumer.ConsumeFromLastOffset)
	ConsumeFromFirstOffset = int(consumer.ConsumeFromFirstOffset)
	ConsumeFromTimestamp   = int(consumer.ConsumeFromTimestamp)
)

var ConsumeFromList = []int{
	ConsumeFromLastOffset,
	ConsumeFromFirstOffset,
	ConsumeFromTimestamp,
}

type Config struct {
	Producers map[string]*ProducerConf `yaml:"producers"`
	Consumers map[string]*ConsumerConf `yaml:"consumers"`
}

type ProducerConf struct {
	NameServerAddrs  []string `yaml:"name_server"` //required
	NameSpace        string   `yaml:"name_space"`  //required
	GroupName        string   `yaml:"group_name"`  //required
	Retries          int      `yaml:"retries"`
	SendMsgTimeoutMS int64    `yaml:"send_msg_timeout_ms"`
}

type ConsumerConf struct {
	Topics               []string `yaml:"topics"`                  //required
	AsyncConsume         bool     `yaml:"async_consume"`           //异步消费
	NameServerAddrs      []string `yaml:"name_server"`             //required
	NameSpace            string   `yaml:"name_space"`              //required
	GroupName            string   `yaml:"group_name"`              //required
	ConsumeFrom          int      `yaml:"consume_from"`            //[0,1,2]
	ConsumeTimestamp     string   `yaml:"consume_timestamp"`       //format：yyyyMMddHHmmss
	FilterHistoryForInit bool     `yaml:"filter_history_for_init"` //当ConsumeFrom=[0,2]时有效，过滤历史消息

	cts time.Time
}
