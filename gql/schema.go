package gql

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

// NewSchema 初始化，Schema需要在服务启动的时候确认
func NewSchema() (*graphql.Schema, error) {

	return nil, nil
}

// NewSchemaWithTableInfo 根据表信息，创建Schema
func NewSchemaWithTableInfo(tables map[string]*TableGraphInfo) (graphql.Schema, error) {
	if len(tables) == 0 {
		return graphql.Schema{}, fmt.Errorf("tables is empty")
	}
	tableFields := make(graphql.Fields)
	for name, table := range tables {
		tableFields[name] = NewTableQueryField(table)
	}
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "root",
		Fields: tableFields,
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		return graphql.Schema{}, fmt.Errorf("new schema err:%s", err)
	}
	return schema, nil
}
