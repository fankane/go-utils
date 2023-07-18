package log

const (
	FormatJSON = "json"
)

type Config struct {
	Level        string `yaml:"level"`
	EnableStdout bool   `yaml:"enable_stdout"`
	Filename     string `yaml:"filename"`
	MaxSize      int    `yaml:"max_size"` // 单位：兆字节
	MaxAge       int    `yaml:"max_age"`  //单位：天
	MaxBackups   int    `yaml:"max_backups"`
	Format       string `yaml:"format"` //[console, json], 默认 console
	Compress     bool   `yaml:"compress"`
}
