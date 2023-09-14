package utime

import (
	"context"
	"errors"
	"github.com/robfig/cron/v3"
)

var (
	ErrCronEmpty = errors.New("cron is nil")
)

// CronDo 定期执行函数, 快速开始，无法停止；完整版请使用下面的 NewCron().Do() 按需开始和停止
func CronDo(cronStr string, fs ...func()) error {
	c := cron.New()
	for _, f := range fs {
		_, err := c.AddFunc(cronStr, f)
		if err != nil {
			return err
		}
	}
	c.Start()
	return nil
}

type Cron struct {
	cron    *cron.Cron
	started bool
}

func NewCron() *Cron {
	return &Cron{cron: cron.New()}
}

func (c *Cron) Start() error {
	if c == nil || c.cron == nil {
		return ErrCronEmpty
	}
	c.cron.Start()
	c.started = true
	return nil
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
// A context is returned so the caller can wait for running jobs to complete.
func (c *Cron) Stop() context.Context {
	if c == nil || c.cron == nil {
		return nil
	}
	defer func() {
		c.started = false
	}()
	return c.cron.Stop()
}

// Do 添加定时任务，如果定时器没开始则自动开始
func (c *Cron) Do(cronStr string, fs ...func()) ([]int, error) {
	entryIDs, err := c.AddFunc(cronStr, fs...)
	if err != nil {
		return nil, err
	}
	if !c.started {
		if err = c.Start(); err != nil {
			c.removeEntryIDs(entryIDs) //已经创建成功的删除
			return nil, err
		}
	}
	return entryIDs, nil
}

// AddFunc 添加定时任务，全部成功返回entryID列表，如果存在失败，则返回错误且部分成功的部分被删除
// 需要调用 c.Start() 后任务才会真正开始执行
func (c *Cron) AddFunc(cronStr string, fs ...func()) ([]int, error) {
	entryIDs := make([]int, 0, len(fs))
	for _, f := range fs {
		entryID, err := c.cron.AddFunc(cronStr, f)
		if err != nil {
			c.removeEntryIDs(entryIDs) //已经创建成功的删除
			return nil, err
		}
		entryIDs = append(entryIDs, int(entryID))
	}
	return entryIDs, nil
}

func (c *Cron) removeEntryIDs(ids []int) {
	for _, id := range ids {
		c.cron.Remove(cron.EntryID(id))
	}
}
