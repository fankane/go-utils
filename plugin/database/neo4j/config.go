package neo4j

type Config struct {
	Target       string `yaml:"target"`
	User         string `yaml:"user"`
	Pwd          string `yaml:"pwd"`
	Realm        string `yaml:"realm"`
	DatabaseName string `yaml:"database_name"`
	AccessMode   int    `yaml:"access_mode"` //0: 写；1：只读[默认0]
}
