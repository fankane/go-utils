package mongodb

import (
	"context"
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	document := bson.M{"name": "tom", "age": 22, "address": "wuhan"}
	insertResult, err := Cli.Cli.Database("testDB").Collection("testCollection").InsertOne(context.Background(), document)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
