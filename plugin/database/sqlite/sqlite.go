package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/fankane/go-utils/plugin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

const (
	defaultSQLiteName = "default"
	pluginType        = "database"
	pluginName        = "sqlite"
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
	mysqlMap := make(map[string]*Config)
	if err := node.Decode(&mysqlMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(mysqlMap) == 0 {
		return fmt.Errorf("mysql config is emtpy")
	}
	for confName, config := range mysqlMap {
		db, err := NewDB(config)
		if err != nil {
			return err
		}
		if confName == defaultSQLiteName {
			DB = db
		}
		mu.Lock()
		dbs[confName] = db
		mu.Unlock()
	}
	return nil
}

func NewDB(config *Config) (*sql.DB, error) {
	if err := validator.New().Struct(config); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", WrapDSN(config))
	if err != nil {
		return nil, fmt.Errorf("[%s] open err:%s", WrapDSN(config), err)
	}
	return db, nil
}

func WrapDSN(config *Config) string {
	if config == nil {
		return ""
	}
	// file:test.db?cache=shared&mode=memory
	dsn := fmt.Sprintf("file:%s?_loc=auto", config.DBFile)
	if config.Cache != "" {
		dsn += fmt.Sprintf("&cache=%s", config.Cache)
	}
	if config.Cache != "" {
		dsn += fmt.Sprintf("&mode=%s", config.Mode)
	}
	return dsn
}
