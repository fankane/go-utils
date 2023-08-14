package kafka

const (
	sendTypeSync  = "sync"
	sendTypeAsync = "async"
)

type Config struct {
	Producers map[string]*ProducerConf `yaml:"producers"`
	Consumers map[string]*ConsumerConf `yaml:"consumers"`
}

type ProducerConf struct {
	Addrs    []string `yaml:"addrs"`
	SendType string   `yaml:"send_type"`
}

type ConsumerConf struct {
	Addrs              []string        `yaml:"addrs"`
	Topics             []string        `yaml:"topics"`
	GroupID            string          `yaml:"group_id"`
	ConcurrencyConsume bool            `yaml:"concurrency_consume"` //并发消费
	ConcurrencyMax     int             `yaml:"concurrency_max"`     //并发数，当concurrency_consume=true时有效
	OffsetInitial      int64           `yaml:"offset_initial"`      //consumer时的offset设置
	OffsetInfo         []*OffsetSingle `yaml:"reset_offset_info"`   //consumerGroup时的offset设置
}

type OffsetSingle struct {
	Topic             string             `yaml:"topic"`
	Offset            int64              `yaml:"offset"`             // -1 newest, [不填默认-1]
	SetForAll         bool               `yaml:"set_for_all"`        // offset 对所有partition生效
	PartitionsSetting []*PartitionOffset `yaml:"partitions_setting"` // 当 SetForAll = false时配置
}

type PartitionOffset struct {
	Partition int32 `yaml:"partition"`
	Offset    int64 `yaml:"offset"` // -2:oldest, -1 newest, [不填默认-1]
}
