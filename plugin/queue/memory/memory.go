package memory

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/utime"
	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
	"sync"
	"time"
)

const (
	pluginType = "queue"
	pluginName = "memory"
)

var (
	DefaultFactory = &Factory{}
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	memConf := &Config{}
	if err := node.Decode(&memConf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	InitQueue(memConf)
	return nil
}

// CachedSize 在内存里面的数据大小
func CachedSize() int64 {
	return globalMemQueue.cachedSize.Load()
}

// CachedLen 在内存里面的消息条数
func CachedLen() int64 {
	return globalMemQueue.cachedLen.Load()
}

type MemQueue struct {
	conf            *Config
	lock            *sync.RWMutex
	topicChannelMap map[string]chan *Message //每个topic的消息通道
	topicMsgSlice   map[string]MessageList   //每个topic的待处理的消息
	topicMsgNotify  map[string]chan bool     //提醒topic 里的消息需要处理了
	cachedSize      *atomic.Int64            //占用空间
	cachedLen       *atomic.Int64            //消息条数
}

var globalMemQueue *MemQueue

func InitQueue(conf *Config) {
	if conf.BufferSize <= 0 {
		conf.BufferSize = 1000
	}
	globalMemQueue = &MemQueue{
		conf:            conf,
		lock:            &sync.RWMutex{},
		topicChannelMap: make(map[string]chan *Message),
		topicMsgSlice:   make(map[string]MessageList),
		topicMsgNotify:  make(map[string]chan bool),
		cachedSize:      atomic.NewInt64(0),
		cachedLen:       atomic.NewInt64(0),
	}
}

func createTopicInfo(topic string) MessageList {
	globalMemQueue.lock.Lock()
	msgList := NewMessageList(topic)
	globalMemQueue.topicChannelMap[topic] = make(chan *Message, globalMemQueue.conf.BufferSize)
	globalMemQueue.topicMsgSlice[topic] = msgList
	globalMemQueue.topicMsgNotify[topic] = make(chan bool)
	globalMemQueue.lock.Unlock()
	go daemonTopicMsg(topic, msgList)
	return msgList
}

// 将到点需要处理的消息，推送到对应 channel
func daemonTopicMsg(topic string, msgList MessageList) {
	for {
		//fmt.Println("daemonTopicMsg, topic:", topic)
		select {
		case <-globalMemQueue.topicMsgNotify[topic]: //需要处理数据了
			//fmt.Println("开始处理队列 list, len:", msgList.Len())
			for msgList.Len() > 0 { //收到一次消息通知，把所有可以处理的消息全部处理，然后等待下一次通知
				msg := msgList.First()
				if msgCanExec(msg) {
					consumerMsg := msgList.pop()
					if consumerMsg == nil {
						fmt.Println("first is nil")
						continue
					}
					//fmt.Println("进消费队列，msg:", string(consumerMsg.Body), ", topic:", topic)
					globalMemQueue.topicChannelMap[topic] <- consumerMsg
					continue
				}
				notifyDealMsg(topic, msg)
				break
			}
			continue
		}
	}
}

func notifyDealMsg(topic string, msg *Message) {
	go func() {
		delay := time.Duration(msg.expectTs - time.Now().UnixNano())
		if msg.Delay == 0 || delay <= 0 {
			globalMemQueue.topicMsgNotify[topic] <- true
			return
		}
		// 当到了下一条需要执行的消息间隔时间后，再次提醒消息处理
		utime.DelayDo(delay, func() error {
			globalMemQueue.topicMsgNotify[topic] <- true
			return nil
		})
	}()
}

func isLenFull() bool {
	if globalMemQueue.conf.MaxLen <= 0 {
		return false
	}
	return globalMemQueue.cachedLen.Load() >= globalMemQueue.conf.MaxLen
}

func isSizeFull() bool {
	if globalMemQueue.conf.MaxSize <= 0 {
		return false
	}
	return globalMemQueue.cachedSize.Load() >= globalMemQueue.conf.MaxSize
}

func addMessage(msg *Message) {
	globalMemQueue.cachedSize.Add(sizeOfMessage(msg))
	globalMemQueue.cachedLen.Inc()
}

func delMessage(msg *Message) {
	globalMemQueue.cachedSize.Sub(sizeOfMessage(msg))
	globalMemQueue.cachedLen.Dec()
}
