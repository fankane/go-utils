package gql

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

var (
	Int          ValidType = graphql.Int
	Float        ValidType = graphql.Float
	String       ValidType = graphql.String
	Boolean      ValidType = graphql.Boolean
	DateTime     ValidType = graphql.DateTime
	ListInt      ValidType = graphql.NewList(graphql.Int)
	ListFloat    ValidType = graphql.NewList(graphql.Float)
	ListString   ValidType = graphql.NewList(graphql.String)
	ListBoolean  ValidType = graphql.NewList(graphql.Boolean)
	ListDateTime ValidType = graphql.NewList(graphql.DateTime)
)

type GType graphql.Type
type ValidType interface{}
type TableQuery func(args map[string]interface{}) (interface{}, error)

type TableGraphInfo struct {
	Name         string
	Fields       map[string]ValidType
	CustomFields map[string]GType //自定义非基础字段类型的字段，可以为空
	QueryFunc    TableQuery
}

// NewTableQueryField Schema定义里面的Fields信息 Resolve方法是提供数据的业务函数
func NewTableQueryField(info *TableGraphInfo) *graphql.Field {
	return &graphql.Field{
		Name: fmt.Sprintf("Query%s", info.Name),
		Type: graphql.NewList(NewTableObject(info)),
		Args: NewTableQueryArgs(info.Fields),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return info.QueryFunc(p.Args)
		},
	}
}

// NewTableQueryArgs 查询表数据时，查询参数定义
func NewTableQueryArgs(fieldM map[string]ValidType) map[string]*graphql.ArgumentConfig {
	res := make(map[string]*graphql.ArgumentConfig)
	for s, validType := range fieldM {
		res[s] = &graphql.ArgumentConfig{Type: validType.(*graphql.Scalar)}
	}
	return res
}

// NewTableObject 对表字段进行封装，提供 graphql 格式的对象
func NewTableObject(info *TableGraphInfo) *graphql.Object {
	fields := NewTableFields(info.Fields)
	if fields == nil {
		fields = make(graphql.Fields)
	}
	if info.CustomFields != nil {
		AppendCustomFields(fields, info.CustomFields)
	}
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:   info.Name,
		Fields: fields,
	})
	return obj
}

func NewTableFields(fieldM map[string]ValidType) graphql.Fields {
	fields := make(graphql.Fields)
	for fieldName, validType := range fieldM {
		fields[fieldName] = &graphql.Field{Type: validType.(*graphql.Scalar)}
	}
	return fields
}

func AppendCustomFields(fields graphql.Fields, customFieldM map[string]GType) {
	if fields == nil {
		return
	}
	for fieldName, t := range customFieldM {
		fields[fieldName] = &graphql.Field{Type: t.(graphql.Type)}
	}
}
