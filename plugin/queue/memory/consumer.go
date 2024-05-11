package memory

import (
	"context"
	"fmt"
	"log"

	"github.com/fankane/go-utils/goroutine"
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
		defer goroutine.Recover()
		for {
			select {
			case msg := <-msgTopic.consumerChan:
				if msg == nil {
					return
				}
				ctx := context.Background()
				ctx = context.WithValue(ctx, CtxPushTs, msg.pushTs)
				ctx = context.WithValue(ctx, CtxDelay, msg.Delay)
				h(ctx, msg.Body)
				go delMessage(topic, msg) //消费一个消息，在全局记录里面，减去这个消息占用的数量，内存大小等信息
			}
		}
	}()
}
