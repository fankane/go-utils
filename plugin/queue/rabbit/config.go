package rabbit

const (
	ContentTypeOS   = "application/octet-stream"
	ContentTypeText = "text/plain"
)

type Config struct {
	Producers map[string]*ProducerConf `yaml:"producers"`
	Consumers map[string]*ConsumerConf `yaml:"consumers"`
}

type ProducerConf struct {
	URL        string `yaml:"url"`
	Durable    bool   `yaml:"durable"`
	AutoDelete bool   `yaml:"auto_delete"`
	Exclusive  bool   `yaml:"exclusive"`
	NoWait     bool   `yaml:"no_wait"`
}

type ConsumerConf struct {
	URL        string   `yaml:"url"`
	Durable    bool     `yaml:"durable"`
	AutoDelete bool     `yaml:"auto_delete"`
	Exclusive  bool     `yaml:"exclusive"`
	NoWait     bool     `yaml:"no_wait"`
	QueueNames []string `yaml:"queue_names"`
}
