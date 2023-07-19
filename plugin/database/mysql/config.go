package mysql

type Config struct {
	DSN                string `yaml:"dsn"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec"`
	ConnMaxIdleTimeSec int    `yaml:"conn_max_idle_time_sec"`
	MaxOpenConn        int    `yaml:"max_open_conn"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
}
