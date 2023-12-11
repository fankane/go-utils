package uerr

import "fmt"

type Err struct {
	Code    int
	Msg     string
	ShowMsg string //不想显示原始信息，替换的字段
}

type Option func(e *Err)

func (e *Err) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("code:%d, msg:%s, showMsg:%s", e.Code, e.Msg, e.ShowMsg)
}
func (e *Err) ECode() int {
	if e == nil {
		return 0
	}
	return e.Code
}
func (e *Err) EMsg() string {
	if e == nil {
		return ""
	}
	return e.Msg
}
func (e *Err) EShowMsg() string {
	if e == nil {
		return ""
	}
	return e.ShowMsg
}

// EMsgPriority 优先显示showMsg,如果没有，展示 Emg
func (e *Err) EMsgPriority() string {
	if e == nil {
		return ""
	}
	if e.ShowMsg != "" {
		return e.ShowMsg
	}
	return e.Msg
}

func NewErr(msg string, opts ...Option) *Err {
	e := &Err{Msg: msg}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Copy 从 err 复制一个对象，如果有opts，则覆盖原对象的值
func Copy(err *Err, opts ...Option) *Err {
	if err == nil {
		return nil
	}
	e := &Err{
		Code:    err.Code,
		Msg:     err.Msg,
		ShowMsg: err.ShowMsg,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// FromError 从 err 创建
func FromError(err error, opts ...Option) *Err {
	if err == nil {
		return nil
	}
	e := &Err{
		Msg: err.Error(),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func ToErr(err error) *Err {
	if err == nil {
		return nil
	}
	e, ok := err.(*Err)
	if ok {
		return e
	}
	return NewErr(err.Error())
}

func WithCode(code int) Option {
	return func(e *Err) {
		e.Code = code
	}
}

func WithMsg(msg string) Option {
	return func(e *Err) {
		e.Msg = msg
	}
}

func WithShowMsg(showMsg string) Option {
	return func(e *Err) {
		e.ShowMsg = showMsg
	}
}
