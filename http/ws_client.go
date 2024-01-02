package http

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

func NewWSClient(url url.URL, opts ...WSOption) (*WSConnection, error) {
	log.Printf("client connecting to %s", url.String())
	wp := &wsParam{}
	for _, opt := range opts {
		opt(wp)
	}

	dialer := getDialer(wp)
	c, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial err:%s", err)
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

func getDialer(wp *wsParam) *websocket.Dialer {
	dialer := websocket.DefaultDialer
	if wp.HandshakeTimeout > 0 {
		dialer.HandshakeTimeout = wp.HandshakeTimeout
	}
	if wp.ReadBufferSize > 0 {
		dialer.ReadBufferSize = wp.ReadBufferSize
	}
	if wp.WriteBufferSize > 0 {
		dialer.WriteBufferSize = wp.WriteBufferSize
	}
	return dialer
}
