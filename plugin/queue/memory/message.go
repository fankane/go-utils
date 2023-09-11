package memory

import (
	"container/list"
	"encoding/json"
	"time"
)

type MessageList interface {
	addMessage(msg *Message)
	pop() *Message
	Len() int
	First() *Message
}

type messageListImpl struct {
	topic string
	list  *list.List
}

type Message struct {
	Body  []byte
	Delay time.Duration //延迟时间

	pushTs   int64 //推送时间 单位： 纳秒
	expectTs int64 //预期执行时间，单位： 纳秒
}

func NewMessageList(topic string) MessageList {
	return &messageListImpl{list: list.New(), topic: topic}
}

func (m *messageListImpl) addMessage(msg *Message) {
	defer notifyDealMsg(m.topic, msg)
	addMessage(m.topic, msg)
	if m.list.Len() == 0 {
		m.list.PushBack(msg)
		return
	}

	// 逆序遍历，将msg按expectTs时间先后顺序排列，expectTs相同的，先到的排在前面
	for i := m.list.Back(); i != nil; i = i.Prev() {
		temp, ok := i.Value.(*Message)
		if !ok {
			return
		}
		if msg.expectTs >= temp.expectTs { //根据预期执行时间找到地方
			m.list.InsertAfter(msg, i)
			return
		}
	}
	m.list.PushFront(msg)
}

func (m *messageListImpl) pop() *Message {
	if m.list.Len() == 0 {
		return nil
	}
	first := m.list.Front()
	msg, ok := first.Value.(*Message)
	if !ok {
		return nil
	}
	m.list.Remove(first)
	return msg
}

func (m *messageListImpl) First() *Message {
	if m.list.Len() == 0 {
		return nil
	}
	first := m.list.Front()
	msg, ok := first.Value.(*Message)
	if !ok {
		return nil
	}
	return msg
}

func (m *messageListImpl) Len() int {
	return m.list.Len()
}

func wrapMessage(value []byte, delay time.Duration) *Message {
	return &Message{
		Body:     value,
		Delay:    delay,
		pushTs:   time.Now().UnixNano(),
		expectTs: time.Now().Add(delay).UnixNano(),
	}
}

func msgCanExec(msg *Message) bool {
	if msg == nil {
		return false
	}
	return msg.expectTs <= time.Now().UnixNano()
}

func sizeOfMessage(msg *Message) int64 {
	if msg == nil {
		return 0
	}
	b, _ := json.Marshal(msg)
	return int64(len(b))
}
