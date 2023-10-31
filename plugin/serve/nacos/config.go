package nacos

import "github.com/nacos-group/nacos-sdk-go/clients/config_client"

type Config struct {
	EndPoint            string `yaml:"end_point"`
	NamespaceID         string `yaml:"namespace_id"`
	TimeoutMs           uint64 `yaml:"timeout_ms"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	DataId              string `yaml:"data_id"`
	Group               string `yaml:"group"`
	NotLoadCacheAtStart bool   `yaml:"not_load_cache_at_start"`
}

type Client struct {
	Cli  config_client.IConfigClient
	conf *Config
}

type MarshalFunc func(newData string, v interface{}) error
