package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch          = 1735660800000                              //时间戳起点 2025-01-01 00:00:00
	datacenterBits = 5                                          //数据中心占用 位数
	workerBits     = 5                                          //节点占用 位数
	sequenceBits   = 12                                         //序号占用 位数
	bitsCnt        = datacenterBits + workerBits + sequenceBits //可自定义位数的总和，固定的

	maxDatacenterID = -1 ^ (-1 << datacenterBits) //最大数据中心ID，掩码
	maxWorkerID     = -1 ^ (-1 << workerBits)     //最大数据节点ID，掩码
	maxSequenceID   = -1 ^ (-1 << sequenceBits)   //最大序号ID，掩码

	timestampLeftShift  = datacenterBits + workerBits + sequenceBits //时间戳偏移量
	datacenterLeftShift = workerBits + sequenceBits                  //数据中心偏移量
	workerLeftShift     = sequenceBits                               //节点偏移量
)

type Snowflake struct {
	l             *sync.Mutex
	lastTimestamp int64
	datacenterID  int64
	workerID      int64
	sequence      int64
	customizeParam
}

type customizeParam struct {
	epoch          int64
	datacenterBits int64
	workerBits     int64
	sequenceBits   int64

	maxDatacenterID int64
	maxWorkerID     int64
	maxSequenceID   int64

	timestampLeftShift  int64
	datacenterLeftShift int64
	workerLeftShift     int64
}
type CPOptions func(c *customizeParam)

var defaultBits = customizeParam{
	epoch:               epoch,
	datacenterBits:      datacenterBits,
	workerBits:          workerBits,
	sequenceBits:        sequenceBits,
	maxDatacenterID:     maxDatacenterID,
	maxWorkerID:         maxWorkerID,
	maxSequenceID:       maxSequenceID,
	timestampLeftShift:  timestampLeftShift,
	datacenterLeftShift: datacenterLeftShift,
	workerLeftShift:     workerLeftShift,
}

var DefaultSnowFlake = &Snowflake{
	l:              &sync.Mutex{},
	lastTimestamp:  -1,
	datacenterID:   1,
	workerID:       1,
	sequence:       0,
	customizeParam: defaultBits,
}

func CustomizeBits(datacenterBits, workerBits, sequenceBits int64) CPOptions {
	return func(c *customizeParam) {
		c.datacenterBits = datacenterBits
		c.workerBits = workerBits
		c.sequenceBits = sequenceBits
	}
}
func CustomizeEpoch(t time.Time) CPOptions {
	return func(c *customizeParam) {
		c.epoch = t.UnixMilli()
	}
}

/*
NewSnowflake
数据中心 + 节点 一共10位，最多支持1024台机器
序列号 一共12 位，同一毫秒内，最多支持 4096个ID
整体支持同一个毫秒内，最多 1024 * 4096 个唯一ID
*/
func NewSnowflake(datacenterID, workerID int64, opts ...CPOptions) (*Snowflake, error) {
	fmt.Println("default:", fmt.Sprintf("%p", &defaultBits))
	dp := defaultBits
	fmt.Println("dp:", fmt.Sprintf("%p", &dp))
	if len(opts) > 0 {
		fmt.Println("opts has")
		for _, opt := range opts {
			opt(&dp)
		}
		if err := checkParams(&dp); err != nil {
			return nil, err
		}
	}
	fmt.Println(dp.epoch, dp.datacenterBits, dp.workerBits, dp.sequenceBits)
	if datacenterID > dp.maxDatacenterID {
		return nil, fmt.Errorf("datacenterID out of limit")
	}
	if workerID > dp.maxWorkerID {
		return nil, fmt.Errorf("workerID out of limit")
	}
	return &Snowflake{
		l:              &sync.Mutex{},
		lastTimestamp:  -1,
		datacenterID:   datacenterID,
		workerID:       workerID,
		sequence:       0,
		customizeParam: dp,
	}, nil
}

func checkParams(dp *customizeParam) error {
	if dp.datacenterBits == datacenterBits && dp.workerBits == workerBits && dp.sequenceBits == sequenceBits {
		return nil
	}
	// 有自定义内容
	if dp.datacenterBits+dp.workerBits+dp.sequenceBits != bitsCnt {
		return fmt.Errorf("customize bits invalid, sum should be:%d", bitsCnt)
	}
	dp.maxDatacenterID = -1 ^ (-1 << dp.datacenterBits)
	dp.maxWorkerID = -1 ^ (-1 << dp.workerBits)
	dp.maxSequenceID = -1 ^ (-1 << dp.sequenceBits)
	dp.timestampLeftShift = dp.datacenterBits + dp.workerBits + dp.sequenceBits
	dp.datacenterLeftShift = dp.workerBits + dp.sequenceBits
	dp.workerLeftShift = dp.sequenceBits
	return nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.l.Lock()
	defer s.l.Unlock()
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if timestamp < s.lastTimestamp { //时间小于上次，机器时钟回拨了
		return 0, fmt.Errorf("machine clock moved backwards")
	}
	if timestamp == s.lastTimestamp { //同一毫秒内，sequence确定
		s.sequence = s.sequence + 1
		if s.sequence > s.maxSequenceID { //到达一毫秒内最大数量限制，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / int64(time.Millisecond)
				time.Sleep(time.Microsecond * 2)
			}
		}
	} else {
		s.sequence = 0 //时间过去了，序号置零
	}
	s.lastTimestamp = timestamp
	id := (timestamp-s.epoch)<<s.timestampLeftShift |
		s.datacenterID<<s.datacenterLeftShift |
		s.workerID<<s.workerLeftShift |
		s.sequence
	return id, nil
}

type ParseSnowflake struct {
	Timestamp    time.Time
	DatacenterID int64
	WorkerID     int64
	Sequence     int64
}

// ParseSnowflakeID 从ID里面反解析出生成时的相关信息
func (s *Snowflake) ParseSnowflakeID(id int64) (*ParseSnowflake, error) {
	if s == nil {
		return nil, fmt.Errorf("s is nil")
	}
	if id < 0 {
		return nil, fmt.Errorf("invalid ID")
	}
	timestamp := (id >> s.timestampLeftShift) + s.epoch
	datacenterID := (id >> s.datacenterLeftShift) & s.maxDatacenterID
	workerID := (id >> s.workerLeftShift) & s.maxWorkerID
	sequence := id & s.maxSequenceID
	return &ParseSnowflake{
		Timestamp:    time.UnixMilli(timestamp),
		DatacenterID: datacenterID,
		WorkerID:     workerID,
		Sequence:     sequence,
	}, nil
}
