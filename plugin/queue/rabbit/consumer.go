package rabbit

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fankane/go-utils/goroutine"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ctxExchange        = "ctx_exchange"
	ctxTs              = "ctx_timestamp"
	ctxRoutingKey      = "ctx_routing_key"
	ctxContentType     = "ctx_content_type"
	ctxContentEncoding = "ctx_content_encoding"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Conf    *ConsumerConf

	lock     *sync.Mutex
	queueMap map[string]amqp.Queue
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
	for _, queueName := range consumerConf.QueueNames {
		if err = c.HandleMessage(name, queueName, h); err != nil {
			return err
		}
	}
	return nil
}

func NewConsumer(conf *ConsumerConf) (*Consumer, error) {
	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{
		Conf:     conf,
		Conn:     conn,
		Channel:  ch,
		lock:     &sync.Mutex{},
		queueMap: map[string]amqp.Queue{},
	}, nil
}

func (c *Consumer) HandleMessage(name, queueName string, h Handler) error {
	var (
		q   amqp.Queue
		ok  bool
		err error
	)
	q, ok = c.queueMap[queueName]
	if !ok {
		q, err = c.Channel.QueueDeclare(queueName, c.Conf.Durable, c.Conf.AutoDelete, c.Conf.Exclusive, c.Conf.NoWait, nil)
		if err != nil {
			return fmt.Errorf("QueueDeclare err:%s", err)
		}
		c.lock.Lock()
		c.queueMap[name] = q
		c.lock.Unlock()
	}
	msgs, err := c.Channel.Consume(q.Name, queueName, false, c.Conf.Exclusive, false, c.Conf.NoWait, nil)
	if err != nil {
		return fmt.Errorf("consume err:%s, queueName:%s", err, q.Name)
	}
	log.Println(queueName, " start handle messages...")
	go func() {
		defer goroutine.Recover()
		for msg := range msgs {
			ctx := context.Background()
			ctx = context.WithValue(ctx, ctxTs, msg.Timestamp)
			ctx = context.WithValue(ctx, ctxExchange, msg.Exchange)
			ctx = context.WithValue(ctx, ctxRoutingKey, msg.RoutingKey)
			ctx = context.WithValue(ctx, ctxContentType, msg.ContentType)
			ctx = context.WithValue(ctx, ctxContentEncoding, msg.ContentEncoding)
			h(ctx, msg.Body)
			if err = msg.Ack(true); err != nil {
				log.Println("ack err:", err)
			}
		}
	}()
	return nil
}

func CtxExchange(ctx context.Context) string {
	return ctx.Value(ctxExchange).(string)
}

func CtxRoutingKey(ctx context.Context) string {
	return ctx.Value(ctxRoutingKey).(string)
}
func CtxContentType(ctx context.Context) string {
	return ctx.Value(ctxContentType).(string)
}
func CtxContentEncoding(ctx context.Context) string {
	return ctx.Value(ctxContentEncoding).(string)
}
func CtxTimestamp(ctx context.Context) int64 {
	return ctx.Value(ctxTs).(int64)
}
