package memory

type Config struct {
	BufferSize  int    `yaml:"buffer_size"`   //排队channel容量，默认1000；超过容量再入消息会阻塞，每个Topic有一个独立的channel
	MaxSize     int64  `yaml:"max_size"`      //最多占用字节，单位：B, 超过无法发送消息到队列
	MaxLen      int64  `yaml:"max_len"`       //最多消息条数，超过无法发送消息到队列
	LoadAtBegin bool   `yaml:"load_at_begin"` //启动时加载数据
	LoadFile    string `yaml:"load_file"`     //加载的文件
}
