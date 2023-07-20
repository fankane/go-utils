package postgres

type Config struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	User               string `yaml:"user"`
	Pwd                string `yaml:"pwd"`
	DBName             string `yaml:"db_name"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec"`
	ConnMaxIdleTimeSec int    `yaml:"conn_max_idle_time_sec"`
	MaxOpenConn        int    `yaml:"max_open_conn"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
}
