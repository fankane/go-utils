package http

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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
	MsgNetClosed  = "use of closed network connection"
)

type wsParam struct {
	// 建立连接时
	HandshakeTimeout time.Duration
	ReadBufferSize   int //不填使用 websocket 自带默认值 4096
	WriteBufferSize  int //不填使用 websocket 自带默认值 4096
	ResponseHeader   http.Header
	ReadErrHandler   func(err error)
	WriteErrHandler  func(err error)
	DisablePingTest  bool
	CheckOriginTrue  bool //处理跨域
	CheckOrigin      func(r *http.Request) bool
}

type WSOption func(p *wsParam)

type HandleWSParam struct {
	W http.ResponseWriter
	R *http.Request
}

func ServerHandleWS(param HandleWSParam, opts ...WSOption) (*WSConnection, error) {
	wp := &wsParam{}
	for _, opt := range opts {
		opt(wp)
	}

	upgrader := getUpgrader(wp)
	c, err := upgrader.Upgrade(param.W, param.R, wp.ResponseHeader)
	if err != nil {
		return nil, fmt.Errorf("upgrade err:%s", err)
	}
	result := &WSConnection{Conn: c, Lock: &sync.Mutex{}}
	if wp.ReadErrHandler != nil {
		result.ReadErrHandler = wp.ReadErrHandler
	}
	if wp.WriteErrHandler != nil {
		result.WriteErrHandler = wp.WriteErrHandler
	}
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
func ReadErrHandler(handler func(err error)) WSOption {
	return func(p *wsParam) {
		p.ReadErrHandler = handler
	}
}
func WriteErrHandler(handler func(err error)) WSOption {
	return func(p *wsParam) {
		p.WriteErrHandler = handler
	}
}

func CheckOriginTrue() WSOption {
	return func(p *wsParam) {
		p.CheckOriginTrue = true
	}
}

// CheckOrigin 会覆盖CheckOriginTrue
func CheckOrigin(co func(r *http.Request) bool) WSOption {
	return func(p *wsParam) {
		p.CheckOrigin = co
	}
}
func DisablePingTest(disable bool) WSOption {
	return func(p *wsParam) {
		p.DisablePingTest = disable
	}
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
	if wp.CheckOriginTrue {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	if wp.CheckOrigin != nil {
		upgrader.CheckOrigin = wp.CheckOrigin
	}
	return upgrader
}
