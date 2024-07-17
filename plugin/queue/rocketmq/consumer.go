package rocketmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/fankane/go-utils/slice"
	"github.com/fankane/go-utils/utime"
)

const (
	ctx_rocketmq_msg_id          = "ctx_rocketmq_msg_id"
	ctx_rocketmq_offset_msg_id   = "ctx_rocketmq_offset_msg_id"
	ctx_rocketmq_queue_offset    = "ctx_rocketmq_queue_offset"
	ctx_rocketmq_store_timestamp = "ctx_rocketmq_store_timestamp"
	ctx_rocketmq_born_timestamp  = "ctx_rocketmq_born_timestamp"
)

type Consumer struct {
	PC   rocketmq.PushConsumer
	conf *ConsumerConf
}

type TopicCreateInfo struct {
	BrokerAddr string
	H          Handler
}

type ConsumeParams struct {
	AutoCreateTopic bool // auto create topic when not exists
	TopicHandler    map[string]TopicCreateInfo
}

type ConsumeOption func(params *ConsumeParams)

func AutoCreateTopic(a bool) ConsumeOption {
	return func(params *ConsumeParams) {
		params.AutoCreateTopic = a
	}
}

func TopicHandler(topic, brokerAddr string, h Handler) ConsumeOption {
	return func(params *ConsumeParams) {
		params.TopicHandler[topic] = TopicCreateInfo{
			BrokerAddr: brokerAddr,
			H:          h,
		}
	}
}

func NewConsumer(conf *ConsumerConf) (*Consumer, error) {
	opts := make([]consumer.Option, 0)
	opts = append(opts, consumer.WithGroupName(conf.GroupName))
	opts = append(opts, consumer.WithNameServer(conf.NameServerAddrs))
	opts = append(opts, consumer.WithNamespace(conf.NameSpace))
	if slice.InInts(conf.ConsumeFrom, ConsumeFromList) {
		opts = append(opts, consumer.WithConsumeFromWhere(consumer.ConsumeFromWhere(conf.ConsumeFrom)))
	}
	if conf.ConsumeTimestamp != "" {
		cTS, err := time.ParseInLocation(utime.LayYMDHms3, conf.ConsumeTimestamp, utime.GetUTC8Loc())
		if err == nil {
			opts = append(opts, consumer.WithConsumeTimestamp(conf.ConsumeTimestamp))
			conf.cts = cTS
		} else {
			conf.cts = utime.GetUTC8Time().Add(-30 * time.Minute) //默认半小时前
		}
	}
	c, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		PC:   c,
		conf: conf,
	}, nil
}

func (c *Consumer) Start(opts ...ConsumeOption) error {
	if c == nil || c.PC == nil {
		return fmt.Errorf("consumer is nil")
	}
	param := &ConsumeParams{
		AutoCreateTopic: false,
		TopicHandler:    make(map[string]TopicCreateInfo),
	}
	for _, opt := range opts {
		opt(param)
	}
	if len(param.TopicHandler) > 0 {
		if err := c.consumeTopics(param); err != nil {
			return err
		}
	}
	return c.PC.Start()
}

func (c *Consumer) consumeTopics(param *ConsumeParams) error {
	ctx := context.Background()
	for topic, temp := range param.TopicHandler {
		if param.AutoCreateTopic {
			exist, err := ExistTopic(ctx, c.conf.NameServerAddrs, topic)
			if err != nil {
				return err
			}
			if !exist {
				if strings.TrimSpace(temp.BrokerAddr) == "" {
					return fmt.Errorf("broker addr is empty")
				}
				if err = CreateTopic(ctx, c.conf.NameServerAddrs, topic, temp.BrokerAddr); err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 100) //等待一下，否则监听的时候，可能新创建的Topic还没同步导致启动失败
			}
		}
		if err := consumeTopic(c.PC, topic, temp.H, c.conf); err != nil {
			return err
		}
	}
	return nil
}

type Handler func(ctx context.Context, value []byte) error

func RegisterHandler(name string, h Handler) error {
	consumerConf, ok := globalConsumerConfs[name]
	if !ok {
		return fmt.Errorf("not found consumer config of [%s]", name)
	}
	c, err := NewConsumer(consumerConf)
	if err != nil {
		return err
	}

	for _, topic := range consumerConf.Topics {
		if err = consumeTopic(c.PC, topic, h, consumerConf); err != nil {
			return err
		}
	}
	if err = c.Start(); err != nil {
		return err
	}
	globalConsumers[name] = c
	return nil
}

func consumeTopic(pc rocketmq.PushConsumer, topic string, h Handler, consumerConf *ConsumerConf) error {
	start := time.Now()
	if err := pc.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, messageExt := range ext {
			if ignore(consumerConf, messageExt, start) {
				continue
			}
			if consumerConf.AsyncConsume {
				go h(setContext(messageExt), messageExt.Body)
			} else {
				h(setContext(messageExt), messageExt.Body)
			}
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		return fmt.Errorf("subscribe failed %s, topic:[%s]", err, topic)
	}
	return nil
}

func ignore(consumerConf *ConsumerConf, msg *primitive.MessageExt, start time.Time) bool {
	if !consumerConf.FilterHistoryForInit {
		return false
	}
	// 如果是第一次启动消费者，会因为找不到历史commit offset 而从头开始消费,
	// 根据消息生产时间，过滤掉消费者创建时的历史消息
	bornTime := time.UnixMilli(msg.BornTimestamp)
	switch consumerConf.ConsumeFrom {
	case ConsumeFromLastOffset:
		return bornTime.Before(start)
	case ConsumeFromTimestamp:
		return bornTime.Before(consumerConf.cts)
	default:
		return false
	}
}

func setContext(ext *primitive.MessageExt) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctx_rocketmq_msg_id, ext.MsgId)
	ctx = context.WithValue(ctx, ctx_rocketmq_offset_msg_id, ext.OffsetMsgId)
	ctx = context.WithValue(ctx, ctx_rocketmq_queue_offset, ext.QueueOffset)
	ctx = context.WithValue(ctx, ctx_rocketmq_store_timestamp, ext.StoreTimestamp)
	ctx = context.WithValue(ctx, ctx_rocketmq_born_timestamp, ext.BornTimestamp)
	return ctx
}

func GetMsgID(ctx context.Context) string {
	v, ok := ctx.Value(ctx_rocketmq_msg_id).(string)
	if !ok {
		return ""
	}
	return v
}
func GetOffsetMsgID(ctx context.Context) string {
	v, ok := ctx.Value(ctx_rocketmq_offset_msg_id).(string)
	if !ok {
		return ""
	}
	return v
}
func GetQueueOffset(ctx context.Context) int64 {
	v, ok := ctx.Value(ctx_rocketmq_queue_offset).(int64)
	if !ok {
		return 0
	}
	return v
}
func GetStoreTimestamp(ctx context.Context) int64 {
	v, ok := ctx.Value(ctx_rocketmq_store_timestamp).(int64)
	if !ok {
		return 0
	}
	return v
}
func GetBornTimestamp(ctx context.Context) int64 {
	v, ok := ctx.Value(ctx_rocketmq_born_timestamp).(int64)
	if !ok {
		return 0
	}
	return v
}
