package es

type Config struct {
	Addr []string `yaml:"addr"`
	User string   `yaml:"user"`
	Pwd  string   `yaml:"pwd"`
}
