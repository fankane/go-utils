package memory

import (
	"errors"
	"time"
)

var (
	ErrClosed       = errors.New("producer has closed")
	ErrOutOfMaxLen  = errors.New("out of max message len")
	ErrOutOfMaxSize = errors.New("out of max message size")
)

type Producer interface {
	SendMessage(topic string, value []byte, opts ...Opts) error
	Close()
}

type memProducer struct {
	closed bool
}

func NewProducer() Producer {
	return &memProducer{}
}

type Opts func(*option)

type option struct {
	Delay time.Duration
}

func (p memProducer) SendMessage(topic string, value []byte, opts ...Opts) error {
	if p.closed {
		return ErrClosed
	}
	if isLenFull() {
		return ErrOutOfMaxLen
	}
	if isSizeFull() {
		return ErrOutOfMaxSize
	}
	optParams := &option{}
	for _, opt := range opts {
		opt(optParams)
	}
	globalMemQueue.lock.RLock()
	topicMsg, ok := globalMemQueue.topicMsgSlice[topic]
	globalMemQueue.lock.RUnlock()
	if !ok { //新建topic
		topicMsg = createTopicInfo(topic)
	}
	// topic 存在，将消息放到指定位置
	topicMsg.addMessage(wrapMessage(value, optParams.Delay))
	return nil
}
func (p memProducer) Close() {
	p.closed = true
}

func Delay(delay time.Duration) Opts {
	return func(o *option) {
		o.Delay = delay
	}
}
