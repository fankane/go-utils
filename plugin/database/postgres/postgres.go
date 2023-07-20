package postgres

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

const (
	defaultDBName = "default"
	pluginType    = "database"
	pluginName    = "postgres"
)

var (
	DB             *sql.DB
	DefaultFactory = &Factory{}
	dbs            = make(map[string]*sql.DB)
	mu             = sync.RWMutex{}
)

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

func GetDB(name string) *sql.DB {
	mu.RLock()
	defer mu.RUnlock()
	return dbs[name]
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	postgresMap := make(map[string]*Config)
	if err := node.Decode(&postgresMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(postgresMap) == 0 {
		return fmt.Errorf("postgres config is emtpy")
	}
	for confName, config := range postgresMap {
		db, err := NewDB(config)
		if err != nil {
			return err
		}
		if confName == defaultDBName {
			DB = db
		}
		mu.Lock()
		dbs[confName] = db
		mu.Unlock()
	}
	return nil
}

func NewDB(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Pwd, config.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open err:%s", err)
	}

	if config.MaxOpenConn > 0 {
		db.SetMaxOpenConns(config.MaxOpenConn)
	}
	if config.MaxIdleConn > 0 {
		db.SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.ConnMaxIdleTimeSec > 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(config.ConnMaxIdleTimeSec))
	}
	if config.ConnMaxLifeTimeSec > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifeTimeSec))
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping err:%s, dsn:%s", err, dsn)
	}
	return db, nil
}
