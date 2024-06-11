package mongodb

import "go.mongodb.org/mongo-driver/mongo"

const (
	DefaultConnTimeoutMs = 10000 //默认10s 超时
)

type Config struct {
	Host             string `yaml:"host"  validate:"required"`
	Port             int    `yaml:"port"  validate:"required"`
	User             string `yaml:"user"`
	Pwd              string `yaml:"pwd"`
	ConnectTimeoutMs int    `yaml:"connect_timeout_ms"`
}

type MongoCli struct {
	Cli *mongo.Client
}
