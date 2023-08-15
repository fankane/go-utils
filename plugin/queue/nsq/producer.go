package nsq

import (
	"fmt"
	go_nsq "github.com/nsqio/go-nsq"
)

type Producer interface {
	SendMsg(topic string, body []byte) error
}

type nsqProducer struct {
	p *go_nsq.Producer
}

func NewProducer(conf *ProducerConf) (Producer, error) {
	p, err := go_nsq.NewProducer(conf.Addr, go_nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	if err = p.Ping(); err != nil {
		return nil, fmt.Errorf("producer ping err:%s", err)
	}
	return &nsqProducer{p: p}, nil
}

func (p *nsqProducer) SendMsg(topic string, body []byte) error {
	if p == nil || p.p == nil {
		return fmt.Errorf("producer is nil")
	}
	return p.p.Publish(topic, body)
}
