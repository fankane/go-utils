package http

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestBLSse(t *testing.T) {
	//url := "http://192.168.0.17:7861/chat/chat"
	url := "http://192.168.99.38:9001/g_hf_management/wechat/sse/status"
	ctx := context.Background()
	code, resChan, err := NewClient().SSEGet(ctx, url)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("statusCode:", code)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range resChan {
			if e == nil {
				fmt.Println("--------Done--------")
				return
			}
			fmt.Println(string(e.Data))
			//tempRes := &TTR{}
			//if err = json.Unmarshal(e.Data, tempRes); err != nil {
			//	fmt.Println("json err:", err)
			//	return
			//}
			//fmt.Println(tempRes.Text)
		}
	}()
	wg.Wait()
}

type TTR struct {
	Text string `json:"text"`
}

var tJSON = `{
  "query": "hello",
  "history": [{"role": "user","content": "遥感的定义"}, {"role": "assistant","content": "遥感（Remote Sensing）是一门技术，它使用电磁波或粒子来探测、识别和收集地表和地球大气层的信息，而无需直接接触目标。这些信息可以用来了解地球的物理、化学和生物学特性，以及它们如何随时间变化。遥感技术可以用于各种应用，包括农业、林业、环境监测、地质勘探、城市规划和灾害管理。"}],
  "stream": true,
  "model_name": "Qwen-14B-Chat",
  "system_messages": "You are a helpful assistant.",
  "temperature": 0.7,
  "max_tokens": 8192,
  "prompt_name": "default"
}`
