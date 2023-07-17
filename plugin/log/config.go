package log

const (
	Development = "dev"
	Production  = "pro"
)

type ConfContent struct {
	Log *Config `yaml:"log"`
}

type Config struct {
	Level        string `yaml:"level"`
	Environment  string
	EnableStdout bool `yaml:"enable_stdout"`
	file         logFileConf
}

type logFileConf struct {
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"` // 单位：兆字节
	MaxAge     int    `yaml:"max_age"`  //单位：天
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
}
