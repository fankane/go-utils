package utime

import (
	"github.com/robfig/cron/v3"
)

// CronDo 定期执行函数
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
