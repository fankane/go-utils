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

type ValidType interface{}
type TableQuery func(args map[string]interface{}) (interface{}, error)

type TableGraphInfo struct {
	Name      string
	Fields    map[string]ValidType
	QueryFunc TableQuery
}

func NewTableQueryField(info *TableGraphInfo) *graphql.Field {
	return &graphql.Field{
		Name: fmt.Sprintf("Query%s", info.Name),
		Type: graphql.NewList(NewTableObject(info.Name, info.Fields)),
		Args: NewTableQueryArgs(info.Fields),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return info.QueryFunc(p.Args)
		},
	}
}

func NewTableQueryArgs(fieldM map[string]ValidType) map[string]*graphql.ArgumentConfig {
	res := make(map[string]*graphql.ArgumentConfig)
	for s, validType := range fieldM {
		res[s] = &graphql.ArgumentConfig{Type: validType.(*graphql.Scalar)}
	}
	return res
}

func NewTableObject(tableName string, fieldM map[string]ValidType) *graphql.Object {
	fields := make(graphql.Fields)
	for fieldName, validType := range fieldM {
		fields[fieldName] = &graphql.Field{Type: validType.(*graphql.Scalar)}
	}
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:   tableName,
		Fields: fields,
	})
	return obj
}
