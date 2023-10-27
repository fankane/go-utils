package rabbit

import (
	"context"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Conf    *ProducerConf

	lock     *sync.Mutex
	queueMap map[string]amqp.Queue
}

func NewProducer(conf *ProducerConf) (*Producer, error) {
	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Producer{
		Conf:     conf,
		Conn:     conn,
		Channel:  ch,
		lock:     &sync.Mutex{},
		queueMap: map[string]amqp.Queue{},
	}, nil
}

type ProducerOptions func(param *ProducerParam)

type ProducerParam struct {
	args      amqp.Table
	exchange  string
	mandatory bool
	immediate bool
}

func Args(args map[string]interface{}) ProducerOptions {
	return func(param *ProducerParam) {
		param.args = args
	}
}
func Exchange(exchange string) ProducerOptions {
	return func(param *ProducerParam) {
		param.exchange = exchange
	}
}
func Mandatory(mandatory bool) ProducerOptions {
	return func(param *ProducerParam) {
		param.mandatory = mandatory
	}
}
func Immediate(immediate bool) ProducerOptions {
	return func(param *ProducerParam) {
		param.immediate = immediate
	}
}

func (p *Producer) SendMsg(ctx context.Context, queueName string, body []byte, opts ...ProducerOptions) error {
	cParam := &ProducerParam{}
	for _, opt := range opts {
		opt(cParam)
	}

	var (
		q   amqp.Queue
		ok  bool
		err error
	)
	q, ok = p.queueMap[queueName]
	if !ok {
		q, err = p.Channel.QueueDeclare(queueName, p.Conf.Durable, p.Conf.AutoDelete, p.Conf.Exclusive, p.Conf.NoWait, cParam.args)
		if err != nil {
			return err
		}
		p.lock.Lock()
		p.queueMap[queueName] = q
		p.lock.Unlock()
	}

	return p.Channel.PublishWithContext(ctx, cParam.exchange, q.Name, cParam.mandatory, cParam.immediate,
		amqp.Publishing{
			ContentType: ContentTypeText,
			Body:        body,
		})
}
