package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	defaultMysqlName = "default"
	pluginType       = "database"
	pluginName       = "mysql"
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
		if confName == defaultMysqlName {
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.User, config.Pwd,
		config.Host, config.Port)
	if config.DBName != "" {
		dsn = fmt.Sprintf("%s%s", dsn, config.DBName)
	}
	if config.Params != "" {
		dsn = fmt.Sprintf("%s?%s", dsn, config.Params)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("[%s] open err:%s", dsn, err)
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
		return nil, fmt.Errorf("[%s] ping err:%s", dsn, err)
	}
	return db, nil
}
