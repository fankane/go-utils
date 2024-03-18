package influx

type Config struct {
	URL        string `yaml:"url"`
	Token      string `yaml:"token"`
	Org        string `yaml:"org"`
	Bucket     string `yaml:"bucket"`
	WriteAsync bool   `yaml:"write_async"` //异步写入
}
