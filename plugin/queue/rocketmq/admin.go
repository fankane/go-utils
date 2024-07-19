package rocketmq

import (
	"context"
	"fmt"
	"github.com/fankane/go-utils/slice"

	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RocketAdmin struct {
}

func CreateTopic(ctx context.Context, nameServerAddrs []string, topic, brokerAddr string, opts ...admin.OptionCreate) error {
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}
	a, err := NewAdmin(nameServerAddrs)
	if err != nil {
		return fmt.Errorf("new admin err:%s, addr:%v", err, nameServerAddrs)
	}
	defer a.Close()
	optsReq := make([]admin.OptionCreate, 0)
	optsReq = append(optsReq, admin.WithTopicCreate(topic))
	optsReq = append(optsReq, admin.WithBrokerAddrCreate(brokerAddr))
	optsReq = append(optsReq, opts...)
	if err = a.CreateTopic(ctx, optsReq...); err != nil {
		return fmt.Errorf("create topic err:%s", err)
	}
	return nil
}

func ExistTopic(ctx context.Context, nameServerAddrs []string, topic string) (bool, error) {
	topicList, err := TopicList(ctx, nameServerAddrs)
	if err != nil {
		return false, err
	}
	if slice.InStrings(topic, topicList) {
		return true, nil
	}
	return false, nil
}

func TopicList(ctx context.Context, nameServerAddrs []string) ([]string, error) {
	a, err := NewAdmin(nameServerAddrs)
	if err != nil {
		return nil, fmt.Errorf("new admin err:%s, addr:%v", err, nameServerAddrs)
	}
	defer a.Close()
	res, err := a.FetchAllTopicList(ctx)
	if err != nil {
		return nil, err
	}
	return res.TopicList, nil
}

func DeleteTopic(ctx context.Context, nameServerAddrs []string, brokerAddr, topic string) error {
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}
	if brokerAddr == "" {
		return fmt.Errorf("brokerAddr is empty")
	}
	a, err := NewAdmin(nameServerAddrs)
	if err != nil {
		return fmt.Errorf("new admin err:%s, addr:%v", err, nameServerAddrs)
	}
	defer a.Close()
	if err = a.DeleteTopic(ctx, admin.WithTopicDelete(topic), admin.WithBrokerAddrDelete(brokerAddr)); err != nil {
		return fmt.Errorf("create topic err:%s", err)
	}
	return nil
}

func NewAdmin(nameServerAddrs []string) (admin.Admin, error) {
	if len(nameServerAddrs) == 0 {
		return nil, fmt.Errorf("nameServerAddrs is empty")
	}
	a, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(nameServerAddrs)))
	if err != nil {
		return nil, fmt.Errorf("new admin err:%s, addr:%v", err, nameServerAddrs)
	}
	return a, nil
}
