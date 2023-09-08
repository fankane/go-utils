package memory

type Config struct {
	BufferSize int   `yaml:"buffer_size"`
	MaxSize    int64 `yaml:"max_size"` //最多占用字节，单位：B, 超过无法发送消息到队列
	MaxLen     int64 `yaml:"max_len"`  //最多消息条数，超过无法发送消息到队列
}
