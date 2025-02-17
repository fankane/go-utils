package limiter

import (
	"sync"
	"time"

	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/utime"
)

/**
滑动窗口限流器

漏桶算法：go.uber.org/ratelimit
令牌桶算法：golang.org/x/time/rate
*/

type Request struct {
	timestamp time.Time
}

type SlidingWindowLimiter struct {
	mu          *sync.Mutex
	bufferSize  int
	maxRequests int
	buffer      []*Request
	head, tail  int
	count       int
	windowSize  time.Duration
}

func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	result := &SlidingWindowLimiter{
		bufferSize:  maxRequests + 1,
		maxRequests: maxRequests,
		buffer:      make([]*Request, maxRequests+1),
		windowSize:  windowSize,
		mu:          &sync.Mutex{},
	}
	go result.clean()
	return result
}

func (s *SlidingWindowLimiter) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.removeOldData(now.Add(-s.windowSize))
	if s.count >= s.maxRequests {
		return false
	}
	s.buffer[s.tail] = &Request{timestamp: now}
	s.tail = (s.tail + 1) % s.bufferSize
	s.count++
	return true
}

func (s *SlidingWindowLimiter) clean() {
	defer goroutine.Recover()
	utime.TickerDo(time.Second, func() error {
		s.removeOldData(time.Now().Add(-s.windowSize))
		return nil
	})
}

var cleanL = &sync.Mutex{}

func (s *SlidingWindowLimiter) removeOldData(start time.Time) {
	cleanL.Lock()
	defer cleanL.Unlock()
	for s.count > 0 && s.buffer[s.head].timestamp.Before(start) {
		s.buffer[s.head] = nil //释放空间
		s.head = (s.head + 1) % s.bufferSize
		s.count--
	}
}
