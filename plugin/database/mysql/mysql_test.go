package mysql

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"testing"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if DB == nil {
		fmt.Println("db is nil")
		return
	}
	rows, err := DB.Query("select name from big_test")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var temp string
		if err = rows.Scan(&temp); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("temp:", temp)
	}
}
