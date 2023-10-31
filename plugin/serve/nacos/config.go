package nacos

import "github.com/nacos-group/nacos-sdk-go/clients/config_client"

type Config struct {
	EndPoint            string       `yaml:"end_point"`
	NamespaceID         string       `yaml:"namespace_id"`
	TimeoutMs           uint64       `yaml:"timeout_ms"`
	Username            string       `yaml:"username"`
	Password            string       `yaml:"password"`
	DataId              string       `yaml:"data_id"`
	Group               string       `yaml:"group"`
	NotLoadCacheAtStart bool         `yaml:"not_load_cache_at_start"`
	LogDir              string       `yaml:"log_dir"`
	ServerConfs         []ServerConf `yaml:"server_confs"`
}

type ServerConf struct {
	IpAddr      string `yaml:"ip_addr"`      //the nacos server address
	Port        uint64 `yaml:"port"`         //the nacos server port
	Scheme      string `yaml:"scheme"`       //the nacos server scheme
	ContextPath string `yaml:"context_path"` //the nacos server contextpath
}

type Client struct {
	Cli  config_client.IConfigClient
	conf *Config
}

type MarshalFunc func(newData string, v interface{}) error
