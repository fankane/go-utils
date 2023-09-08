package memory

import (
	"context"
	"fmt"
)

const (
	CtxPushTs = "_ctx_push_ts"
	CtxDelay  = "_ctx_delay"
)

type Handler func(ctx context.Context, value []byte) error

var consumers = make(map[string]struct{})

func RegisterHandler(topic string, h Handler) error {
	if _, ok := consumers[topic]; ok {
		return fmt.Errorf("topic:%s already has consumer", topic)
	}
	consumers[topic] = struct{}{}
	daemonConsume(topic, h)
	return nil
}

func daemonConsume(topic string, h Handler) {
	go func() {
		for {
			if _, ok := globalMemQueue.topicChannelMap[topic]; !ok {
				createTopicInfo(topic)
			}
			select {
			case msg := <-globalMemQueue.topicChannelMap[topic]:
				ctx := context.Background()
				ctx = context.WithValue(ctx, CtxPushTs, msg.pushTs)
				ctx = context.WithValue(ctx, CtxDelay, msg.Delay)
				h(ctx, msg.Body)
			}
		}
	}()
}
