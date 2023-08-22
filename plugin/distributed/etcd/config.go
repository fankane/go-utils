package etcd

type Config struct {
	Username      string      `yaml:"username" `
	Password      string      `yaml:"password" `
	Endpoints     []string    `yaml:"endpoints" validate:"required"`
	DialTimeoutMs int64       `yaml:"dial_timeout_ms"` //单位：毫秒
	OpenDiscovery bool        `yaml:"open_discovery"`  //是否打开服务发现功能
	SInfo         *ServerInfo `yaml:"server_info"`
}

type ServerInfo struct {
	ServerName    string `yaml:"server_name" validate:"required"`
	ServerID      string `yaml:"server_id" validate:"required"`
	Region        string `yaml:"region"`
	Host          string `yaml:"host"`
	CheckInterval int64  `yaml:"check_interval"` //服务检测间隔，单位：秒

	stop bool
}
