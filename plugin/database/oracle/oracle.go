package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	_ "github.com/godror/godror"
	"gopkg.in/yaml.v3"
)

const (
	defaultOracleName = "default"
	pluginType        = "database"
	pluginName        = "oracle"
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
	oracleMap := make(map[string]*Config)
	if err := node.Decode(&oracleMap); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if len(oracleMap) == 0 {
		return fmt.Errorf("oracle config is emtpy")
	}
	for confName, conf := range oracleMap {
		db, err := NewDB(conf)
		if err != nil {
			return err
		}
		if confName == defaultOracleName {
			DB = db
		}
		mu.Lock()
		dbs[confName] = db
		mu.Unlock()
	}
	return nil
}

func NewDB(config *Config) (*sql.DB, error) {
	dsn := WrapDSN(config)
	db, err := sql.Open("godror", dsn)
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
	ctx := context.Background()
	if config.PingTimeoutMs > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Millisecond*time.Duration(config.PingTimeoutMs))
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[%s] ping err:%s", dsn, err)
	}
	return db, nil
}

func WrapDSN(config *Config) string {
	if config == nil {
		return ""
	}
	dsn := fmt.Sprintf("%s/%s@%s:%d/", config.User, config.Pwd, config.Host, config.Port)
	if config.Sid != "" {
		dsn = fmt.Sprintf("%s%s", dsn, config.Sid)
	} else if config.ServerName != "" {
		dsn = fmt.Sprintf("%s%s", dsn, config.ServerName)
	}
	return dsn
}
