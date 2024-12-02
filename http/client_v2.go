package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fankane/go-utils/str"
)

type DoParams struct {
	BaseHeader     http.Header //发起请求前设置基础 header
	ResponseHeader http.Header
	AppendHeader   map[string]string //在baseHeader基础上增加 key-value 对
}

type DoOption func(params *DoParams)

func BaseHeader(header http.Header) DoOption {
	return func(params *DoParams) {
		params.BaseHeader = header
	}
}

// RespHeader header should not be nil, otherwise can't get response header
func RespHeader(header http.Header) DoOption {
	return func(params *DoParams) {
		params.ResponseHeader = header
	}
}

func AppendHeader(append map[string]string) DoOption {
	return func(params *DoParams) {
		params.AppendHeader = append
	}
}

func (c *Client) CTXPost(ctx context.Context, url string, data []byte, opts ...DoOption) (int, []byte, error) {
	return c.doRequestV2(ctx, http.MethodPost, url, data, opts...)
}

func (c *Client) CTXPostJSONByte(ctx context.Context, url string, data []byte, opts ...DoOption) (int, []byte, error) {
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeJSON}))
	return c.CTXPost(ctx, url, data, opts...)
}

func (c *Client) CTXPostJSON(ctx context.Context, url, json string, opts ...DoOption) (int, []byte, error) {
	return c.CTXPostJSONByte(ctx, url, str.ToBytes(json), opts...)
}

func (c *Client) CTXPostForm(ctx context.Context, url string, data url.Values, opts ...DoOption) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeForm}))
	return c.CTXPost(ctx, url, body, opts...)
}

func (c *Client) CTXGet(ctx context.Context, url string, opts ...DoOption) (int, []byte, error) {
	return c.doRequestV2(ctx, http.MethodGet, url, nil, opts...)
}

func (c *Client) CTXDeleteJSON(ctx context.Context, url, json string, opts ...DoOption) (int, []byte, error) {
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeJSON}))
	return c.doRequestV2(ctx, http.MethodDelete, url, str.ToBytes(json), opts...)
}

func (c *Client) CTXDeleteForm(ctx context.Context, url string, data url.Values, opts ...DoOption) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeForm}))
	return c.doRequestV2(ctx, http.MethodDelete, url, body, opts...)
}

func (c *Client) CTXPutJSON(ctx context.Context, url, json string, opts ...DoOption) (int, []byte, error) {
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeJSON}))
	return c.doRequestV2(ctx, http.MethodPut, url, str.ToBytes(json), opts...)
}

func (c *Client) CTXPutForm(ctx context.Context, url string, data url.Values, opts ...DoOption) (int, []byte, error) {
	body, err := io.ReadAll(strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	opts = append(opts, AppendHeader(map[string]string{ContentType: ContentTypeForm}))
	return c.doRequestV2(ctx, http.MethodPut, url, body, opts...)
}

func (c *Client) doRequestV2(ctx context.Context, method, url string, data []byte, opts ...DoOption) (int, []byte, error) {
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
	case http.MethodDelete:
		req, err = http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer(data))
	case http.MethodPut:
		req, err = http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	}
	if doParams.BaseHeader != nil {
		req.Header = doParams.BaseHeader
	}
	for key, val := range doParams.AppendHeader {
		req.Header.Add(key, val)
	}
	c.setShortConn(req)
	resp, err := c.cli.Do(req)
	if err != nil {
		return 0, nil, err
	}
	setRespHeader(resp, doParams)
	return parseDoResp(resp)
}

func setRespHeader(resp *http.Response, doParams *DoParams) {
	if doParams.ResponseHeader == nil {
		return
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			doParams.ResponseHeader.Add(k, vv)
		}
	}
}
