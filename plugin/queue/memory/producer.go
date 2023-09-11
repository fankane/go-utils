package memory

import (
	"errors"
	"time"
)

var (
	ErrClosed           = errors.New("producer has closed")
	ErrInternal         = errors.New("internal error")
	ErrOutOfMaxLen      = errors.New("out of max message len")
	ErrOutOfMaxSize     = errors.New("out of max message size")
	ErrConsumerNotFound = errors.New("consumer not found")
)
var producers = make([]Producer, 0)

type Producer interface {
	SendMessage(topic string, value []byte, opts ...Opts) error
	Close()
}

type memProducer struct {
	closed bool
}

func NewProducer() Producer {
	res := &memProducer{}
	lock.Lock()
	producers = append(producers, res)
	lock.Unlock()
	return res
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
	topicInfo, ok := globalMemQueue.topicInfo[topic]
	globalMemQueue.lock.RUnlock()
	if !ok { //新建topic
		topicInfo = createTopicInfo(topic)
	}
	if topicInfo == nil {
		return ErrInternal
	}
	// topic 存在，将消息放到指定位置
	topicInfo.msgSlice.addMessage(wrapMessage(value, optParams.Delay))
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
