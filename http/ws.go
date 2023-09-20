package http

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func NewWSClient(url url.URL, handler WsMessageHandler, opts ...WSOption) (*WSCommonInfo, error) {
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
	result := &WSCommonInfo{Conn: c}
	handleMessage(result, handler)
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
