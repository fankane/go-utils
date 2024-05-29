package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fankane/go-utils/plugin/log"
)

/**
SSE 后台接口
*/

var ErrRequestClosed = errors.New("request closed")

type SSEParam struct {
	Header map[string]string
	Log    *log.Log
}

type SSEOption func(client *SSEParam)

func Header(h map[string]string) SSEOption {
	return func(client *SSEParam) {
		client.Header = h
	}
}
func Log(l *log.Log) SSEOption {
	return func(client *SSEParam) {
		client.Log = l
	}
}
func RegisterSSE(w http.ResponseWriter, r *http.Request, dataChan chan []byte, opts ...SSEOption) error {
	if w == nil {
		return fmt.Errorf("responseWriter is nil")
	}
	p := &SSEParam{}
	for _, opt := range opts {
		opt(p)
	}
	if len(p.Header) > 0 {
		for k, v := range p.Header {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	for {
		select {
		case data := <-dataChan:
			if p.Log != nil {
				p.Log.Debugf("send data:[%s]", string(preData)+string(data)+splitMsg)
			}

			fmt.Fprintf(w, string(preData)+string(data)+splitMsg)
			// 检查连接是否关闭
			f, ok := w.(http.Flusher)
			if !ok {
				return fmt.Errorf("ResponseWriter assert Flusher failed")
			}
			f.Flush()
		case <-r.Context().Done(): //客户端断开连接，则停止发送消息
			if p.Log != nil {
				p.Log.Debugf("client close conn")
			}
			return ErrRequestClosed
		}
	}
	return nil
}
