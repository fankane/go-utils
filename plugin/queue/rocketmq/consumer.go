package rocketmq

import (
	"context"
	"fmt"
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
	PC rocketmq.PushConsumer
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
		PC: c,
	}, nil
}

func (c *Consumer) Start() error {
	if c == nil || c.PC == nil {
		return fmt.Errorf("consumer is nil")
	}
	return c.PC.Start()
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
	start := time.Now()
	for _, topic := range consumerConf.Topics {
		if err = c.PC.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for _, messageExt := range ext {
				if ignore(consumerConf, messageExt, start) {
					continue
				}
				if consumerConf.AsyncConsume {
					go h(setContext(messageExt), messageExt.Body)
				} else {
					h(setContext(messageExt), messageExt.Body)
					//if errHandle := h(setContext(messageExt), messageExt.Body); errHandle != nil {
					//	return -1, errHandle
					//}
				}
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
			return fmt.Errorf("subscribe failed %s", err)
		}
	}
	if err = c.Start(); err != nil {
		return err
	}
	globalConsumers[name] = c
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
