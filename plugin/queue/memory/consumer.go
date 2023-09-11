package memory

import (
	"context"
	"fmt"
	"log"
)

const (
	CtxPushTs = "_ctx_push_ts"
	CtxDelay  = "_ctx_delay"
)

type Handler func(ctx context.Context, value []byte) error

type memConsumer struct {
}

var consumers = make(map[string]*memConsumer)

func RegisterHandler(topic string, h Handler) error {
	if _, ok := consumers[topic]; ok {
		return fmt.Errorf("topic:%s already has consumer", topic)
	}

	tempCon := &memConsumer{}
	lock.Lock()
	consumers[topic] = tempCon
	lock.Unlock()
	daemonConsume(topic, h)
	return nil
}

func daemonConsume(topic string, h Handler) {
	msgTopic, ok := globalMemQueue.topicInfo[topic]
	if !ok {
		msgTopic = createTopicInfo(topic)
	}
	if msgTopic == nil {
		log.Printf("create topic info failed")
		return
	}
	go func() {
		for {
			select {
			case msg := <-msgTopic.consumerChan:
				if msg == nil {
					return
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, CtxPushTs, msg.pushTs)
				ctx = context.WithValue(ctx, CtxDelay, msg.Delay)
				go delMessage(topic, msg)
				h(ctx, msg.Body)
			}
		}
	}()
}
