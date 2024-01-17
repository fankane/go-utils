package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Producer struct {
	P rocketmq.Producer
}

func NewProducer(conf *ProducerConf) (*Producer, error) {
	opts := make([]producer.Option, 0)
	opts = append(opts, producer.WithNameServer(conf.NameServerAddrs))
	opts = append(opts, producer.WithGroupName(conf.GroupName))
	opts = append(opts, producer.WithNamespace(conf.NameSpace))
	if conf.Retries > 0 {
		opts = append(opts, producer.WithRetry(conf.Retries))
	}
	if conf.SendMsgTimeoutMS > 0 {
		opts = append(opts, producer.WithSendMsgTimeout(time.Millisecond*time.Duration(conf.SendMsgTimeoutMS)))
	}
	p, err := rocketmq.NewProducer(opts...)
	if err != nil {
		return nil, err
	}
	return &Producer{P: p}, nil
}

func (p *Producer) Start() error {
	if p == nil || p.P == nil {
		return fmt.Errorf("producer is nil")
	}
	return p.P.Start()
}

func (p *Producer) Shutdown() error {
	if p == nil || p.P == nil {
		return fmt.Errorf("producer is nil")
	}
	return p.P.Shutdown()
}

type SendMsgParams struct {
	Tag        string
	Keys       []string
	DelayLevel int
}

type SendMsgOption func(params *SendMsgParams)

func WithTag(tag string) SendMsgOption {
	return func(params *SendMsgParams) {
		params.Tag = tag
	}
}
func WithKeys(keys []string) SendMsgOption {
	return func(params *SendMsgParams) {
		params.Keys = keys
	}
}

/*
*
delay level definition:
1  2   3    4    5   6   7   8   9  10  11  12  13  14   15  16  17 18
1s 5s  10s  30s  1m  2m  3m  4m  5m 6m  7m  8m  9m  10m  20m 30m 1h 2h
*/
func WithDelayLevel(delayLevel int) SendMsgOption {
	return func(params *SendMsgParams) {
		params.DelayLevel = delayLevel
	}
}

func (p *Producer) SendSync(ctx context.Context, topic string, msg []byte, opts ...SendMsgOption) (*primitive.SendResult, error) {
	if p == nil || p.P == nil {
		return nil, fmt.Errorf("producer is nil")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic is empty")
	}

	params := &SendMsgParams{}
	for _, opt := range opts {
		opt(params)
	}

	sendMsg := primitive.NewMessage(topic, msg)
	if params.Tag != "" {
		sendMsg.WithTag(params.Tag)
	}
	if len(params.Keys) > 0 {
		sendMsg.WithKeys(params.Keys)
	}
	if params.DelayLevel > 0 {
		sendMsg.WithDelayTimeLevel(params.DelayLevel)
	}
	return p.P.SendSync(ctx, sendMsg)
}
