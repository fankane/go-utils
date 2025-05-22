package memory

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/utime"
	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
)

const (
	pluginType = "queue"
	pluginName = "memory"
)

var (
	DefaultFactory = &Factory{}
	lock           = sync.Mutex{}
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
	if memConf.LoadAtBegin && memConf.LoadFile == "" {
		return fmt.Errorf("load_file empty")
	}
	if err := InitQueue(memConf); err != nil {
		return fmt.Errorf("init err:%s", err)
	}
	return nil
}

// CachedSize 在内存里面的数据大小
func CachedSize() int64 {
	return globalMemQueue.cachedSize.Load()
}

// CachedLen 在内存里面的消息条数, 还没开始到消费channel的，以及已经在channel里面，但是还没执行的数量
func CachedLen() int64 {
	return globalMemQueue.cachedLen.Load()
}

// TopicCachedLen 在内存里面的消息条数, 还没开始到消费channel的，以及已经在channel里面，但是还没执行的数量
func TopicCachedLen(topic string) int64 {
	for name, tp := range globalMemQueue.topicInfo {
		if name == topic {
			return tp.msgLen.Load()
		}
	}
	return 0
}

// TopicCachedSize 在内存里面的数据大小
func TopicCachedSize(topic string) int64 {
	for name, tp := range globalMemQueue.topicInfo {
		if name == topic {
			return tp.dataSize.Load()
		}
	}
	return 0
}

// Topics topic列表
func Topics() []string {
	result := make([]string, 0)
	for topic, _ := range globalMemQueue.topicInfo {
		result = append(result, topic)
	}
	return result
}

// StopAllProducer 停止所有的生产者，不允许生产数据
func StopAllProducer() {
	for _, producer := range producers {
		producer.Close()
	}
}

// StopAllConsumer 停止所有消费者
func StopAllConsumer() error {
	for topic, _ := range consumers {
		if err := StopConsumer(topic); err != nil {
			return err
		}
	}
	return nil
}

// StopConsumer 停止消费某个topic的数据
func StopConsumer(topic string) error {
	_, ok := consumers[topic]
	if !ok {
		return ErrConsumerNotFound
	}
	delete(consumers, topic)

	topicInfo, ok := globalMemQueue.topicInfo[topic]
	if !ok {
		return errors.New("can not stop, not found topic info")
	}
	topicInfo.consumerChan <- nil //写入一条nil，表示关闭，不再消费
	return nil
}

// ClearTopicMsg 清空Topic所有未到消费时间的消息
func ClearTopicMsg(topic string) error {
	topicInfo, ok := globalMemQueue.topicInfo[topic]
	if !ok {
		return errors.New("not found topic info")
	}
	topicInfo.msgSlice.Clear()
	return nil
}

// ClearAllTopicMsg 清空所有未到消费时间的消息
func ClearAllTopicMsg() error {
	funcList := make([]func() error, 0)
	errMsg := make([]string, 0)
	cl := &sync.Mutex{}
	for topic, _ := range globalMemQueue.topicInfo {
		clearTopicName := topic
		funcList = append(funcList, func() error {
			if err := ClearTopicMsg(clearTopicName); err != nil {
				cl.Lock()
				errMsg = append(errMsg, fmt.Sprintf("%s clear failed %s", clearTopicName, err))
				cl.Unlock()
			}
			return nil
		})
	}
	if err := goroutine.Exec(funcList); err != nil {
		return fmt.Errorf("concurrent clear err:%s", err)
	}
	if len(errMsg) > 0 {
		return fmt.Errorf(strings.Join(errMsg, ";"))
	}
	return nil
}

type MemQueue struct {
	conf       *Config
	lock       *sync.RWMutex
	topicInfo  map[string]*memTopic //每个topic的消息
	cachedSize *atomic.Int64        //占用空间
	cachedLen  *atomic.Int64        //消息条数
}

type memTopic struct { //topic信息
	msgSlice     MessageList   //消息链表，所有消息先进链表，到时间后，进 consumerChan
	consumerChan chan *Message //到延迟时间了，需要被消费的数据
	msgNotify    chan bool     //提醒topic 里的消息需要处理了
	dataSize     *atomic.Int64 //占用空间
	msgLen       *atomic.Int64 //消息条数
}

var (
	globalMemQueue *MemQueue
	inited         bool
)

func InitQueue(conf *Config) error {
	if inited {
		return nil
	}
	if conf.BufferSize <= 0 {
		conf.BufferSize = 1000
	}
	globalMemQueue = &MemQueue{
		conf:       conf,
		lock:       &sync.RWMutex{},
		topicInfo:  make(map[string]*memTopic),
		cachedSize: atomic.NewInt64(0),
		cachedLen:  atomic.NewInt64(0),
	}
	if conf.LoadAtBegin {
		return LoadFileData(conf.LoadFile)
	}
	inited = true
	return nil
}

func createTopicInfo(topic string) *memTopic {
	msgList := NewMessageList(topic)
	mt := &memTopic{
		msgSlice:     msgList,
		consumerChan: make(chan *Message, globalMemQueue.conf.BufferSize),
		msgNotify:    make(chan bool),
		dataSize:     atomic.NewInt64(0),
		msgLen:       atomic.NewInt64(0),
	}
	globalMemQueue.lock.Lock()
	globalMemQueue.topicInfo[topic] = mt
	globalMemQueue.lock.Unlock()
	go daemonTopicMsg(topic, msgList)
	return mt
}

// 将到点需要处理的消息，推送到对应 channel
func daemonTopicMsg(topic string, msgList MessageList) {
	msgTopic, ok := globalMemQueue.topicInfo[topic]
	if !ok {
		log.Printf("topic:%s not exist", topic)
		return
	}
	for {
		select {
		case <-msgTopic.msgNotify: //需要处理数据了
			for msgList.Len() > 0 { //收到一次消息通知，把所有可以处理的消息全部处理，然后等待下一次通知
				msg := msgList.First()
				if msgCanExec(msg) {
					consumerMsg := msgList.pop()
					if consumerMsg == nil {
						continue
					}
					msgTopic.consumerChan <- consumerMsg
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
	msgTopic, ok := globalMemQueue.topicInfo[topic]
	if !ok || msgTopic == nil {
		log.Printf("topic:%s not exist", topic)
		return
	}
	go func() {
		delay := time.Duration(msg.expectTs - time.Now().UnixNano())
		if msg.Delay == 0 || delay <= 0 {
			msgTopic.msgNotify <- true
			return
		}
		// 当到了下一条需要执行的消息间隔时间后，再次提醒消息处理
		utime.DelayDo(delay, func() error {
			msgTopic.msgNotify <- true
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

func addMessage(topic string, msg *Message) {
	msgSize := sizeOfMessage(msg)
	globalMemQueue.cachedSize.Add(msgSize)
	globalMemQueue.cachedLen.Inc()

	msgTopic, ok := globalMemQueue.topicInfo[topic]
	if ok {
		msgTopic.dataSize.Add(msgSize)
		msgTopic.msgLen.Inc()
	}
}

func delMessage(topic string, msg *Message) {
	msgSize := sizeOfMessage(msg)
	globalMemQueue.cachedSize.Sub(msgSize)
	globalMemQueue.cachedLen.Dec()
	msgTopic, ok := globalMemQueue.topicInfo[topic]
	if ok {
		msgTopic.dataSize.Sub(msgSize)
		msgTopic.msgLen.Dec()
	}
}
