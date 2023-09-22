package neo4j

import (
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

	res, err := Cli.Session.Run("match (n:Person) return n", nil)
	if err != nil {
		fmt.Println("run err:", err)
		return
	}

	records, err := res.Collect()
	if err != nil {
		fmt.Println("Collect err:", err)
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
	}
	return
}
func TestCreate(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	res, err := Cli.Session.Run("create (hf:Person{birday:\"1993-09-11\",gender:\"male\",phone:\"0755\"})", nil)
	//res, err := Session.Session.Run("match (n:Person) where id(n) = 5 or id(n)=6 delete n", nil)
	if err != nil {
		fmt.Println("run err:", err)
		return
	}

	records, err := res.Collect()
	if err != nil {
		fmt.Println("Collect err:", err)
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
