package http

import (
	"context"
	"strings"
	"sync"

	"github.com/fankane/go-utils/goroutine"
	"github.com/gorilla/websocket"
)

type WSConnection struct {
	Conn            *websocket.Conn
	Lock            *sync.Mutex
	ReadErrHandler  func(err error)
	WriteErrHandler func(err error)

	closed bool
}

func (h *WSConnection) Close() error {
	if h == nil || h.Conn == nil {
		return ErrConnClosed
	}
	if err := h.Conn.Close(); err != nil {
		return err
	}
	h.Conn = nil //关闭后，连接置空
	h.closed = true
	return nil
}

func (h *WSConnection) ReadMessage() (int, []byte, error) {
	if h == nil || h.Conn == nil {
		return 0, nil, ErrConnClosed
	}
	mt, data, err := h.Conn.ReadMessage()
	if err != nil {
		if h.ReadErrHandler != nil {
			h.ReadErrHandler(err)
		}
		h.Close()
	}
	return mt, data, err
}
func (h *WSConnection) WriteMessage(messageType int, data []byte) error {
	if h == nil || h.Conn == nil {
		return ErrConnClosed
	}
	h.Lock.Lock()
	defer h.Lock.Unlock()

	err := h.Conn.WriteMessage(messageType, data)
	if err != nil && h.WriteErrHandler != nil {
		h.WriteErrHandler(err)
	}
	return err
}

// ListenMessage 循环不断读取数据流
func (h *WSConnection) ListenMessage(ctx context.Context, handler func(ctx context.Context, mt int, data []byte)) error {
	if h == nil || h.Conn == nil {
		return ErrConnClosed
	}
	go func() {
		defer goroutine.Recover()
		for {
			if h.closed {
				return
			}
			messageType, message, err := h.ReadMessage()
			if err != nil {
				return
			}
			go handler(ctx, messageType, message)
		}
	}()
	return nil
}

func IsNetWorkConnectClosed(err error) bool {
	return strings.Contains(err.Error(), MsgNetClosed)
}
