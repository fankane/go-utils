package oracle

type Config struct {
	Host               string `yaml:"host"  validate:"required"`
	Port               int    `yaml:"port"  validate:"required"`
	User               string `yaml:"user"  validate:"required"`
	Pwd                string `yaml:"pwd"  validate:"required"`
	Sid                string `yaml:"sid"`
	ServerName         string `yaml:"server_name"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec"`
	ConnMaxIdleTimeSec int    `yaml:"conn_max_idle_time_sec"`
	MaxOpenConn        int    `yaml:"max_open_conn"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
	PingTimeoutMs      int    `yaml:"ping_timeout_ms"` //建立连接ping的超时时间，单位：ms
}
