package redis

type Config struct {
	Addr               string `yaml:"addr"`
	DB                 int    `yaml:"db"`
	User               string `yaml:"user"`
	Pwd                string `yaml:"pwd"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec"`
	ConnMaxIdleTimeSec int    `yaml:"conn_max_idle_time_sec"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
	MinIdleConn        int    `yaml:"min_idle_conn"`
	DialTimeoutMils    int    `yaml:"dial_timeout_mils"`
	PingTimeoutMils    int    `yaml:"ping_timeout_mils"`
}
