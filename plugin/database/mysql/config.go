package mysql

import "database/sql"

const (
	TableTypeView = "VIEW"
	TableTypeBase = "BASE TABLE"
	NullYes       = "YES"
	NullNo        = "NO"
)

type Config struct {
	Host               string `yaml:"host"  validate:"required"`
	Port               int    `yaml:"port"  validate:"required"`
	User               string `yaml:"user"  validate:"required"`
	Pwd                string `yaml:"pwd"  validate:"required"`
	DBName             string `yaml:"db_name"`
	Params             string `yaml:"params"`
	ConnMaxLifeTimeSec int    `yaml:"conn_max_life_time_sec"`
	ConnMaxIdleTimeSec int    `yaml:"conn_max_idle_time_sec"`
	MaxOpenConn        int    `yaml:"max_open_conn"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
}

// TableColumn SHOW FULL COLUMNS FROM db_name.table_name 返回结果
type TableColumn struct {
	Field      string         `json:"field"`
	Type       sql.NullString `json:"type"`
	Collation  string         `json:"collation"`
	Null       string         `json:"null"`
	Key        string         `json:"key"` // [PRI,unique,index]
	Default    sql.NullString `json:"default"`
	Extra      string         `json:"extra"`
	Privileges string         `json:"privileges"`
	Comment    string         `json:"comment"`
}

// TableStatus show TABLE STATUS; 返回结果
type TableStatus struct {
	Name          string         `json:"Name"`
	Engine        string         `json:"Engine"`
	Version       sql.NullInt64  `json:"Version"`
	RowFormat     sql.NullString `json:"Row_format"`
	Rows          sql.NullInt64  `json:"Rows"`
	AvgRowLength  sql.NullInt64  `json:"Avg_row_length"`
	DataLength    sql.NullInt64  `json:"Data_length"`
	MaxDataLength sql.NullInt64  `json:"Max_data_length"`
	IndexLength   sql.NullInt64  `json:"Index_length"`
	DataFree      sql.NullInt64  `json:"Data_free"`
	AutoIncrement sql.NullInt64  `json:"Auto_increment"`
	CreateTime    sql.NullTime   `json:"Create_time"`
	UpdateTime    sql.NullTime   `json:"Update_time"`
	CheckTime     sql.NullTime   `json:"Check_time"`
	Collation     sql.NullString `json:"Collation"`
	CheckSum      sql.NullString `json:"Checksum"`
	CreateOptions sql.NullString `json:"Create_options"`
	Comment       string         `json:"Comment"`
}
