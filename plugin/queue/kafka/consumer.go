package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/str"
)

var (
	sessionCloseError = errors.New("kafka consumer group session close")
	messageCloseError = errors.New("kafka consumer group claim message close")
	defaultConMax     = 1000
	ctxTopic          = "ctx_topic"
	ctxPartition      = "ctx_partition"
	ctxOffset         = "ctx_offset"
	ctxTs             = "ctx_timestamp"
)

type Handler func(ctx context.Context, key, value []byte) error

func RegisterHandler(name string, h Handler) error {
	consumerConf, ok := globalConsumerMap[name]
	if !ok {
		return fmt.Errorf("not found consumer config of [%s]", name)
	}
	consumerGroup, err := NewConsumerGroup(consumerConf)
	if err != nil {
		return err
	}
	go func() {
		err = consumerGroup.Consume(context.Background(), consumerConf.Topics, &consumerHandler{
			cg:      consumerGroup,
			handler: h,
			conf:    consumerConf,
		})
	}()
	time.Sleep(time.Millisecond * 100) //等待consumerGroup.Consume完成
	if err != nil {
		return err
	}
	return nil
}

func NewConsumerGroup(conf *ConsumerConf) (sarama.ConsumerGroup, error) {
	defConf := getDefaultConf()
	if conf.GroupID == "" {
		conf.GroupID = str.UUID()
	}
	cg, err := sarama.NewConsumerGroup(conf.Addrs, conf.GroupID, defConf)
	if err != nil {
		return nil, err
	}
	return cg, nil
}

type consumerHandler struct {
	handler Handler
	cg      sarama.ConsumerGroup
	conf    *ConsumerConf
	tm      *goroutine.TaskManager
}

func (h *consumerHandler) Setup(s sarama.ConsumerGroupSession) error {
	if h.conf.OffsetInfo == nil || len(h.conf.OffsetInfo) == 0 {
		return nil
	}
	topicPartition := s.Claims()

	fmt.Println("topicPartition:", topicPartition)
	for _, topicSetting := range h.conf.OffsetInfo {

		if _, ok := topicPartition[topicSetting.Topic]; !ok {
			log.Printf("topic:[%s] not exists...", topicSetting.Topic)
			continue
		}
		if topicSetting.Offset == sarama.OffsetNewest {
			continue
		}
		if topicSetting.SetForAll {
			partitions, ok := topicPartition[topicSetting.Topic]
			if !ok {
				return fmt.Errorf("setting topic:[%s] partition offset failed: cannot found partition info", topicSetting.Topic)
			}
			for _, partition := range partitions {
				fmt.Println(fmt.Sprintf("ResetOffset topic:%s, partition:%d, offset:%d", topicSetting.Topic, partition, topicSetting.Offset))
				s.ResetOffset(topicSetting.Topic, partition, topicSetting.Offset, "")
			}
		} else {
			for _, offset := range topicSetting.PartitionsSetting {
				fmt.Println(fmt.Sprintf("Single ResetOffset topic:%s, partition:%d, offset:%d", topicSetting.Topic, offset.Partition, offset.Offset))
				s.ResetOffset(topicSetting.Topic, offset.Partition, offset.Offset, "")
			}
		}
	}
	return nil
}
func (h *consumerHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	if h.tm != nil {
		h.tm.GracefulRelease()
	}
	return nil
}
func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if h.conf.ConcurrencyConsume {
		max := h.conf.ConcurrencyMax
		if max <= 0 {
			max = defaultConMax
		}
		taskManger, err := goroutine.NewTaskManager(goroutine.WithRunnerNum(max))
		if err != nil {
			return fmt.Errorf("new concurrency manager failed err:%s", err)
		}
		h.tm = taskManger
	}
	for {
		select {
		case <-sess.Context().Done(): // 判断session是否结束
			return sessionCloseError
		case msg, ok := <-claim.Messages(): // 监听消息
			if !ok { // msg close  and return
				fmt.Println("messageCloseError")
				return messageCloseError
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, ctxTopic, msg.Topic)
			ctx = context.WithValue(ctx, ctxPartition, msg.Partition)
			ctx = context.WithValue(ctx, ctxOffset, msg.Offset)
			ctx = context.WithValue(ctx, ctxTs, msg.Timestamp)
			if h.tm != nil { //并发执行
				h.tm.AddTask(func() {
					h.handler(ctx, msg.Key, msg.Value)
					sess.MarkMessage(msg, "")
				})
			} else {
				h.handler(ctx, msg.Key, msg.Value)
				sess.MarkMessage(msg, "")
			}
		}
	}
	return nil
}

func Topic(ctx context.Context) string {
	return ctx.Value(ctxTopic).(string)
}
func Partition(ctx context.Context) int32 {
	return ctx.Value(ctxPartition).(int32)
}
func Offset(ctx context.Context) int64 {
	return ctx.Value(ctxOffset).(int64)
}
func Timestamp(ctx context.Context) time.Time {
	return ctx.Value(ctxTs).(time.Time)
}
