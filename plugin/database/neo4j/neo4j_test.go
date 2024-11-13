package neo4j

import (
	"context"
	"fmt"
	"testing"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
)

func TestMatch(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	//cql := "MATCH (p1:`项目名称`)-[r]->(p2:`申请情况`) RETURN p1,r, p2"
	cql := "MATCH (n) RETURN n"
	//cql := "MATCH (p1:Person)-[r]->(p2:Person) RETURN p1, r, p2" //本地数据库
	ctx := context.Background()
	records, err := Dri.Run(ctx, cql, nil)
	if err != nil {
		fmt.Println("run err:", err)
		return
	}

	fmt.Println("len:", len(records))
	fmt.Println("detail:", str.ToJSON(records))
	for _, record := range records {
		//fmt.Println("keys len", len(record.Keys))
		fmt.Println("record:", str.ToJSON(record))
		b, ok := record.Get("birday")
		if ok {
			fmt.Println("dirthday:", b)
		} else {
			fmt.Println("birthday not exist")
		}
		for _, key := range record.Keys {
			v, ok := record.Get(key)
			if !ok {
				fmt.Println(key, "不存在")
				continue
			}
			fmt.Println("key:", key, ", val:", str.ToJSON(v))
		}
	}
	return
}

func TestCreate(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	ctx := context.Background()
	records, err := Dri.Run(ctx, "create (hf2:Person{birday:\"1993-09-11\",gender:\"male\",phone:\"3333\"})", nil)
	if err != nil {
		fmt.Println("run err:", err)
		return
	}
	fmt.Println(str.ToJSON(records))
	for _, record := range records {
		//fmt.Println("keys", record.Keys)
		for _, key := range record.Keys {
			v, ok := record.Get(key)
			if !ok {
				fmt.Println(key, "不存在")
				continue
			}
			fmt.Println("key:", key, ", val:", str.ToJSON(v))
		}
		//fmt.Println("values", record.Values)
	}
	return
}
