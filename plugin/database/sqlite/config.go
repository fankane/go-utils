package sqlite

type Config struct {
	DBFile string `yaml:"db_file"  validate:"required"`
	Mode   string `yaml:"mode"`  // [ro, rw, rwc, memory] 等
	Cache  string `yaml:"cache"` // [shared,private]
}
