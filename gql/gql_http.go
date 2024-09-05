package gql

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// 提供Http服务

type ServeHTTP func(w http.ResponseWriter, r *http.Request)

// GraphGinServe Gin 框架里面的Handle方法，直接设置在router里面即可
func GraphGinServe(schema *graphql.Schema) gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GraphHttpServe(schema *graphql.Schema) ServeHTTP {
	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func GraphGinServeWithTableInfos(tables map[string]*TableGraphInfo) (gin.HandlerFunc, error) {
	schema, err := NewSchemaWithTableInfo(tables)
	if err != nil {
		return nil, err
	}
	return GraphGinServe(&schema), nil
}

func GraphHttpServeWithTableInfos(tables map[string]*TableGraphInfo) (ServeHTTP, error) {
	schema, err := NewSchemaWithTableInfo(tables)
	if err != nil {
		return nil, err
	}
	return GraphHttpServe(&schema), nil
}
