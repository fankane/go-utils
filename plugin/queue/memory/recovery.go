package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fankane/go-utils/file"
	"github.com/fankane/go-utils/str"
)

type backData struct {
	Topic    string        `json:"topic"`
	Body     []byte        `json:"body"`
	Delay    time.Duration `json:"delay"`
	PushTs   int64         `json:"push_ts"`
	ExpectTs int64         `json:"expect_ts"`
}

func LoadFileData(f string) error {
	lines, err := file.ReadLine(f)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if err = backLine(line); err != nil {
			return err
		}
	}
	return nil
}

// StopAndBackOnFile 备份内存数据到文件，会导致所有消费者停止，所有生产者也停止
func StopAndBackOnFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file err:%s", err)
	}
	defer f.Close()
	StopAllProducer()
	for topic, _ := range globalMemQueue.topicInfo {
		msgList, err := getTopicCachedMsg(topic)
		if err != nil {
			return err
		}
		if err = writeTopicData(f, topic, msgList); err != nil {
			return fmt.Errorf("back failed err:%s", err)
		}
	}
	return nil
}

func backLine(line string) error {
	data := &backData{}
	if err := json.Unmarshal([]byte(line), data); err != nil {
		return err
	}

	msgTopic, ok := globalMemQueue.topicInfo[data.Topic]
	if !ok {
		msgTopic = createTopicInfo(data.Topic)
	}
	msgTopic.msgSlice.addMessage(&Message{
		Body:     data.Body,
		Delay:    data.Delay,
		pushTs:   data.PushTs,
		expectTs: data.ExpectTs,
	})
	return nil
}

func getTopicCachedMsg(topic string) ([]*Message, error) {
	existConsumer := true
	if _, ok := consumers[topic]; ok {
		if err := StopConsumer(topic); err != nil {
			return nil, err
		}
	} else {
		// 没有创建消费者，
		existConsumer = false
	}

	// 先获取尚未执行的 chan 里面，再获取排在列表里面的，保证顺序
	tpInfo, ok := globalMemQueue.topicInfo[topic]
	if !ok {
		return nil, fmt.Errorf("topic:%s not exist", topic)
	}
	if !existConsumer {
		tpInfo.consumerChan <- nil //写入一条nil，表示关闭，不再消费
	}

	res := make([]*Message, 0)
	chanFinish := false
	for {
		if chanFinish {
			fmt.Println("chan 获取完毕")
			break
		}
		select {
		case msg := <-tpInfo.consumerChan:
			if msg == nil {
				fmt.Println("get nil msg")
				chanFinish = true
				break
			}
			fmt.Println("get data:", str.ToJSON(msg))
			res = append(res, msg)
		}
	}

	for tpInfo.msgSlice.Len() > 0 {
		res = append(res, tpInfo.msgSlice.pop())
	}
	return res, nil
}

// 按行写入，恢复时可按行恢复，预留以后：防止文件太大的时候，读取文件可以不读取全量
func writeTopicData(f *os.File, topic string, msgList []*Message) error {
	for _, message := range msgList {
		tempData := &backData{
			Topic:    topic,
			Body:     message.Body,
			Delay:    message.Delay,
			PushTs:   message.pushTs,
			ExpectTs: message.expectTs,
		}
		tempBytes, _ := json.Marshal(tempData)
		f.Write(tempBytes)
		f.WriteString("\n")
	}
	return nil
}
