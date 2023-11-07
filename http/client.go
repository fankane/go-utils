package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fankane/go-utils/str"
)

const (
	Connection         = "Connection"
	ContentType        = "Content-Type"
	ContentTypeJSON    = "application/json"
	ContentTypeXML     = "application/xml"
	ContentTypeForm    = "application/x-www-form-urlencoded"
	ContentTypeMulForm = "multipart/form-data"
)

type Client struct {
	Host      string
	Timeout   time.Duration
	ShortConn bool //是否使用短连接
	cli       *http.Client
	Transport http.RoundTripper
}

type ClientOption func(client *Client)

func NewClient(opts ...ClientOption) *Client {
	cli := &Client{}
	for _, opt := range opts {
		opt(cli)
	}
	c := &http.Client{}
	if cli.Timeout.Nanoseconds() > 0 {
		c.Timeout = cli.Timeout
	}
	cli.cli = c
	return cli
}

func WithHost(host string) ClientOption {
	return func(client *Client) {
		client.Host = host
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.Timeout = timeout
	}
}

func WithShortConn(shortConn bool) ClientOption {
	return func(client *Client) {
		client.ShortConn = shortConn
	}
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(client *Client) {
		client.Transport = transport
	}
}

func (c *Client) Post(url string, header map[string]string, data []byte) (int, []byte, error) {
	return c.doRequest(http.MethodPost, url, data, header)
}

func (c *Client) PostJSONByte(url string, data []byte) (int, []byte, error) {
	return c.Post(url, map[string]string{ContentType: ContentTypeJSON}, data)
}

func (c *Client) PostJSON(url, json string) (int, []byte, error) {
	return c.Post(url, map[string]string{ContentType: ContentTypeJSON}, str.ToBytes(json))
}

func (c *Client) PostForm(url string, data url.Values) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	return c.Post(url, map[string]string{ContentType: ContentTypeForm}, body)
}

func (c *Client) Get(url string) (int, []byte, error) {
	return c.doRequest(http.MethodGet, url, nil, nil)
}

func (c *Client) GetWithHeader(url string, header map[string]string) (int, []byte, error) {
	return c.doRequest(http.MethodGet, url, nil, header)
}

func (c *Client) DeleteJSON(url, json string) (int, []byte, error) {
	return c.doRequest(http.MethodDelete, url, str.ToBytes(json), map[string]string{ContentType: ContentTypeJSON})
}

func (c *Client) DeleteForm(url string, data url.Values) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	return c.doRequest(http.MethodDelete, url, body, map[string]string{ContentType: ContentTypeForm})
}

func (c *Client) PutJSON(url, json string) (int, []byte, error) {
	return c.doRequest(http.MethodPut, url, str.ToBytes(json), map[string]string{ContentType: ContentTypeJSON})
}

func (c *Client) PutForm(url string, data url.Values) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	return c.doRequest(http.MethodPut, url, body, map[string]string{ContentType: ContentTypeForm})
}

func (c *Client) doRequest(method, url string, data []byte, header map[string]string) (int, []byte, error) {
	var (
		req *http.Request
		err error
	)
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	case http.MethodDelete:
		req, err = http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(data))
	case http.MethodPut:
		req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	}
	for key, val := range header {
		req.Header.Add(key, val)
	}
	c.setShortConn(req)
	resp, err := c.cli.Do(req)
	if err != nil {
		return 0, nil, err
	}
	return parseDoResp(resp)
}

func (c *Client) setShortConn(req *http.Request) {
	if c.ShortConn {
		req.Close = true
		req.Header.Set(Connection, "close")
	}
}

func parseDoResp(resp *http.Response) (int, []byte, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}
