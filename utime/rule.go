package utime

import "time"

type TickerOption struct {
	Max              int           //最多执行次数, 默认一直执行
	ReturnWhenError  bool          //当出错的时候返回, 默认不返回
	DoExactly        bool          //准确执行：即使上一个任务还没完成，时间到了，也会执行下一次. 默认 false
	FirstImmediately bool          //第一次马上执行，默认false
	BreakDuration    time.Duration //到达指定时长后，就停止, 默认不返回
}

type TkOptions func(opt *TickerOption)

// DelayDo 延迟 delay 时间后执行 do
func DelayDo(delay time.Duration, do func() error) error {
	t := time.NewTimer(delay)
	defer t.Stop()
	for {
		<-t.C
		return do()
	}
}

// TickerDo 每隔 every 时间后执行 do
func TickerDo(every time.Duration, do func() error, opts ...TkOptions) error {
	tOpt := &TickerOption{}
	for _, opt := range opts {
		opt(tOpt)
	}
	t := time.NewTicker(every)
	defer t.Stop()
	i := 0              //计数器
	start := time.Now() //开始时间，计时

	var (
		err    error
		finish bool
	)
	doFunc := func() error {
		if tOpt.BreakDuration.Nanoseconds() > 0 &&
			time.Since(start).Nanoseconds() >= tOpt.BreakDuration.Nanoseconds() { //有配置时间限制
			finish = true
			return nil
		}
		i++
		err = do()
		if err != nil && tOpt.ReturnWhenError {
			finish = true
			return err
		}
		if tOpt.Max > 0 && i >= tOpt.Max { //如果配置了且到达次数了就退出
			finish = true
			return nil
		}
		return nil
	}

	once := func() {
		if tOpt.DoExactly {
			go doFunc()
		} else {
			doFunc()
		}
	}
	if tOpt.FirstImmediately { //进来就第一次执行，减少第一次等待的时间
		once()
		if finish {
			return err
		}
	}
	for range t.C {
		once()
		if finish {
			return err
		}
	}
	return nil
}

func WithMax(max int) TkOptions {
	return func(opt *TickerOption) {
		opt.Max = max
	}
}

func WithReturn(needReturn bool) TkOptions {
	return func(opt *TickerOption) {
		opt.ReturnWhenError = needReturn
	}
}

func WithBreakDuration(duration time.Duration) TkOptions {
	return func(opt *TickerOption) {
		opt.BreakDuration = duration
	}
}

func WithDoExactly(exactly bool) TkOptions {
	return func(opt *TickerOption) {
		opt.DoExactly = exactly
	}
}

func WithFirstImmediately(firstImmediately bool) TkOptions {
	return func(opt *TickerOption) {
		opt.FirstImmediately = firstImmediately
	}
}
