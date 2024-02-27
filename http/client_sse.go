package http

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/fankane/go-utils/goroutine"
	"net/http"
	"time"
)

/**

参考：https://www.ruanyifeng.com/blog/2017/05/server-sent_events.html

每一次发送的信息，由若干个message组成，每个message之间用\n\n分隔。每个message内部由若干行组成，每一行都是如下格式。
[field]: value\n

上面的field可以取四个值。

data: 数据内容
event：event字段表示自定义的事件类型，默认是message事件
id： 数据标识符用id字段表示，相当于每一条数据的编号。
retry：服务器可以用retry字段，指定浏览器重新发起连接的时间间隔


此外，还可以有冒号开头的行，表示注释。通常，服务器每隔一段时间就会向浏览器发送一个注释，保持连接不中断。
: This is a comment
*/

var (
	preID      = []byte("id:")
	preData    = []byte("data:")
	preEvent   = []byte("event:")
	preRetry   = []byte("retry:")
	preComment = []byte(":")
)

type SSEEvent struct {
	Timestamp time.Time
	ID        []byte
	Data      []byte
	Event     []byte
	Retry     []byte
	Comment   []byte
}

func (c *Client) SSEGet(ctx context.Context, url string, opts ...DoOption) (int, chan *SSEEvent, error) {
	return c.doSSERequest(ctx, http.MethodGet, url, nil, opts...)
}
func (c *Client) SSEPost(ctx context.Context, url string, data []byte, opts ...DoOption) (int, chan *SSEEvent, error) {
	return c.doSSERequest(ctx, http.MethodPost, url, data, opts...)
}
func (c *Client) SSEPut(ctx context.Context, url string, data []byte, opts ...DoOption) (int, chan *SSEEvent, error) {
	return c.doSSERequest(ctx, http.MethodPut, url, data, opts...)
}

func (c *Client) doSSERequest(ctx context.Context, method, url string, data []byte, opts ...DoOption) (int, chan *SSEEvent, error) {
	var (
		req *http.Request
		err error
	)
	doParams := &DoParams{}
	for _, opt := range opts {
		opt(doParams)
	}
	switch method {
	case http.MethodGet:
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	case http.MethodPost:
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	case http.MethodPut:
		req, err = http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	default:
		return 0, nil, fmt.Errorf("unsupport method:%s", method)
	}

	if doParams.BaseHeader != nil {
		req.Header = doParams.BaseHeader
	}
	for key, val := range doParams.AppendHeader {
		req.Header.Add(key, val)
	}
	req.Header.Set("Accept", "text/event-stream")
	resp, err := c.cli.Do(req)
	if err != nil {
		return 0, nil, err
	}
	return ReadResponse(resp)
}

func ReadResponse(resp *http.Response) (int, chan *SSEEvent, error) {
	if resp == nil {
		return 0, nil, fmt.Errorf("response is nil")
	}
	resultChan := make(chan *SSEEvent, 100)
	go func() {
		defer goroutine.Recover()
		scanner := bufio.NewScanner(resp.Body)
		split := func(data []byte, atEOF bool) (int, []byte, error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			// We have a full event payload to parse.
			if i, nlen := containsDoubleNewline(data); i >= 0 {
				return i + nlen, data[0:i], nil
			}
			// If we're at EOF, we have all of the data.
			if atEOF {
				return len(data), data, nil
			}
			// Request more data.
			return 0, nil, nil
		}
		scanner.Split(split)
		for scanner.Scan() {
			resultChan <- getEvent(scanner.Bytes())
		}
		resultChan <- nil //表示结束
	}()
	return resp.StatusCode, resultChan, nil
}

func getEvent(line []byte) *SSEEvent {
	var e SSEEvent
	switch {
	case bytes.HasPrefix(line, preID):
		e.ID = append([]byte(nil), trimHeader(len(preID), line)...)
	case bytes.HasPrefix(line, preData):
		e.Data = append(e.Data[:], append(trimHeader(len(preData), line), byte('\n'))...)
	// The spec says that a line that simply contains the string "data" should be treated as a data field with an empty body.
	case bytes.Equal(line, bytes.TrimSuffix(preData, []byte(":"))):
		e.Data = append(e.Data, byte('\n'))
	case bytes.HasPrefix(line, preEvent):
		e.Event = append([]byte(nil), trimHeader(len(preEvent), line)...)
	case bytes.HasPrefix(line, preRetry):
		e.Retry = append([]byte(nil), trimHeader(len(preRetry), line)...)
	default:
		// Ignore any garbage that doesn't match what we're looking for.
	}
	//fmt.Println("eeeee:", string(e.Data))
	return &e
}

func trimHeader(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}

func containsDoubleNewline(data []byte) (int, int) {
	// Search for each potentially valid sequence of newline characters
	crcr := bytes.Index(data, []byte("\r\r"))
	lflf := bytes.Index(data, []byte("\n\n"))
	crlflf := bytes.Index(data, []byte("\r\n\n"))
	lfcrlf := bytes.Index(data, []byte("\n\r\n"))
	crlfcrlf := bytes.Index(data, []byte("\r\n\r\n"))
	// Find the earliest position of a double newline combination
	minPos := minPosInt(crcr, minPosInt(lflf, minPosInt(crlflf, minPosInt(lfcrlf, crlfcrlf))))
	// Detemine the length of the sequence
	nlen := 2
	if minPos == crlfcrlf {
		nlen = 4
	} else if minPos == crlflf || minPos == lfcrlf {
		nlen = 3
	}
	return minPos, nlen
}

// Returns the minimum non-negative value out of the two values. If both
// are negative, a negative value is returned.
func minPosInt(a, b int) int {
	if a < 0 {
		return b
	}
	if b < 0 {
		return a
	}
	if a > b {
		return b
	}
	return a
}
