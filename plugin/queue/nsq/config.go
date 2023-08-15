package nsq

type Config struct {
	Producers map[string]*ProducerConf `yaml:"producers"`
	Consumers map[string]*ConsumerConf `yaml:"consumers"`
}

type ProducerConf struct {
	Addr     string `yaml:"addr"`
	SendType string `yaml:"send_type"`
}

type ConsumerConf struct {
	Addrs              []string `yaml:"addrs"`
	Topic              string   `yaml:"topic"`
	Channel            string   `yaml:"channel"`
	ConcurrencyConsume bool     `yaml:"concurrency_consume"` //并发消费
	ConcurrencyMax     int      `yaml:"concurrency_max"`     //并发数，当concurrency_consume=true时有效
}
