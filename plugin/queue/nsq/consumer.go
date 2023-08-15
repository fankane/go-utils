package nsq

import (
	"context"
	"fmt"
	go_nsq "github.com/nsqio/go-nsq"
)

const (
	defaultConMax  = 1000
	ctxAttempts    = "ctx_attempts"
	ctxTs          = "ctx_timestamp"
	ctxNSQDAddress = "ctx_nsqd_address"
)

type Handler func(ctx context.Context, body []byte) error

type nsqHandler struct {
	h    Handler
	conf *ConsumerConf
}

func RegisterHandler(name string, h Handler) error {
	consumerConf, ok := globalConsumerMap[name]
	if !ok {
		return fmt.Errorf("not found consumer config of [%s]", name)
	}
	c, err := go_nsq.NewConsumer(consumerConf.Topic, consumerConf.Channel, go_nsq.NewConfig())
	if err != nil {
		return err
	}
	nsqH := &nsqHandler{
		h:    h,
		conf: consumerConf,
	}
	if consumerConf.ConcurrencyConsume {
		conMax := consumerConf.ConcurrencyMax
		if conMax <= 0 {
			conMax = defaultConMax
		}
		c.AddConcurrentHandlers(nsqH, conMax)
	} else {
		c.AddHandler(nsqH)
	}
	if err = c.ConnectToNSQLookupds(consumerConf.Addrs); err != nil {
		return fmt.Errorf("ConnectToNSQLookupds err:%s", err)
	}
	return nil
}

func (n nsqHandler) HandleMessage(message *go_nsq.Message) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxTs, message.Timestamp)
	ctx = context.WithValue(ctx, ctxAttempts, message.Attempts)
	ctx = context.WithValue(ctx, ctxNSQDAddress, message.NSQDAddress)
	return n.h(ctx, message.Body)
}

func NSQDAddress(ctx context.Context) string {
	return ctx.Value(ctxNSQDAddress).(string)
}

func Attempts(ctx context.Context) uint16 {
	return ctx.Value(ctxAttempts).(uint16)
}

func Timestamp(ctx context.Context) int64 {
	return ctx.Value(ctxTs).(int64)
}
