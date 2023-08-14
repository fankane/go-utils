package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(topic string, key, value []byte) (partition int32, offset int64, err error)
	Close()
}

func NewSyncProducer(conf *ProducerConf) (Producer, error) {
	defaultConf := getDefaultConf()
	defaultConf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(conf.Addrs, defaultConf)
	if err != nil {
		return nil, err
	}
	return &syncProducer{producer: producer}, nil
}

func NewAsyncProducer(conf *ProducerConf) (Producer, error) {
	producer, err := sarama.NewAsyncProducer(conf.Addrs, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case er := <-errors:
				if er != nil {
					log.Printf("async err:%s", er)
				}
			case <-success:
			}
		}
	}(producer)
	return &asyncProducer{producer: producer}, nil
}

type syncProducer struct {
	producer sarama.SyncProducer
}

func (s *syncProducer) SendMessage(topic string, key, value []byte) (partition int32, offset int64, err error) {
	return s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
}

func (s *syncProducer) Close() {
	if s.producer == nil {
		return
	}
	s.producer.Close()
}

type asyncProducer struct {
	producer sarama.AsyncProducer
}

func (s *asyncProducer) SendMessage(topic string, key, value []byte) (partition int32, offset int64, err error) {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	s.producer.Input() <- message
	return
}

func (s *asyncProducer) Close() {
	if s.producer == nil {
		return
	}
	s.producer.Close()
}
