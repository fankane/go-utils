package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fankane/go-utils/goroutine"
	"github.com/gorilla/websocket"
)

/**
websocket 服务端
*/

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
	CloseMessage  = websocket.CloseMessage
	PingMessage   = websocket.PingMessage
	PongMessage   = websocket.PongMessage
)

var (
	ErrConnClosed = errors.New("conn closed")
)

type wsParam struct {
	// 建立连接时
	HandshakeTimeout time.Duration
	ReadBufferSize   int //不填使用 websocket 自带默认值 4096
	WriteBufferSize  int //不填使用 websocket 自带默认值 4096
	ResponseHeader   http.Header
}

type WSOption func(p *wsParam)

// WsMessageHandler 服务端出来websocket方法，函数返回的Byte数组是写回客户端的
// needResponse: 是否需要把 body 写回客户端
// closeConn: 是否断开连接，不再读取数据
type WsMessageHandler func(ctx context.Context, messageType int, p []byte) (needResponse, closeConn bool, body []byte)

type HandleWSParam struct {
	W http.ResponseWriter
	R *http.Request
	F WsMessageHandler
}

type WSCommonInfo struct {
	Conn *websocket.Conn
}

var (
	lock         = &sync.Mutex{}
	globalWSConn = map[string]*WSCommonInfo{}
)

func ServerHandleWS(param HandleWSParam, opts ...WSOption) (*WSCommonInfo, error) {
	wp := &wsParam{}
	for _, opt := range opts {
		opt(wp)
	}

	upgrader := getUpgrader(wp)
	c, err := upgrader.Upgrade(param.W, param.R, wp.ResponseHeader)
	if err != nil {
		return nil, fmt.Errorf("upgrade err:%s", err)
	}
	result := &WSCommonInfo{Conn: c}
	handleMessage(result, param.F)
	return result, nil
}

func HandshakeTimeout(timeout time.Duration) WSOption {
	return func(p *wsParam) {
		p.HandshakeTimeout = timeout
	}
}
func ReadBufferSize(size int) WSOption {
	return func(p *wsParam) {
		p.ReadBufferSize = size
	}
}
func WriteBufferSize(size int) WSOption {
	return func(p *wsParam) {
		p.WriteBufferSize = size
	}
}
func ResponseHeader(header http.Header) WSOption {
	return func(p *wsParam) {
		p.ResponseHeader = header
	}
}

func (h *WSCommonInfo) Close() error {
	if h == nil || h.Conn == nil {
		return ErrConnClosed
	}
	if err := h.Conn.Close(); err != nil {
		return err
	}
	return nil
}

func (h *WSCommonInfo) WriteMessage(messageType int, data []byte) error {
	if h == nil || h.Conn == nil {
		return ErrConnClosed
	}
	return h.Conn.WriteMessage(messageType, data)
}

func getUpgrader(wp *wsParam) websocket.Upgrader {
	var upgrader = websocket.Upgrader{}
	if wp.HandshakeTimeout > 0 {
		upgrader.HandshakeTimeout = wp.HandshakeTimeout
	}
	if wp.ReadBufferSize > 0 {
		upgrader.ReadBufferSize = wp.ReadBufferSize
	}
	if wp.WriteBufferSize > 0 {
		upgrader.WriteBufferSize = wp.WriteBufferSize
	}
	return upgrader
}

func handleMessage(wInfo *WSCommonInfo, f WsMessageHandler) {
	go func() {
		defer goroutine.Recover()
		for {
			mt, message, err := wInfo.Conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure,
					websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					return
				}
				log.Printf("read message err:%s", err)
				return
			}
			ok, closed, resp := f(context.Background(), mt, message)
			if !ok {
				continue
			}
			if closed { //需要断开连接，不再读取数据
				wInfo.Conn.Close()
				return
			}
			err = wInfo.Conn.WriteMessage(mt, resp)
			if err != nil {
				log.Printf("read message err:%s", err)
				return
			}
		}
	}()

}
