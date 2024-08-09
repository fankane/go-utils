package sqlite

import (
	"fmt"
	"testing"

	"github.com/fankane/go-utils/plugin"
)

func TestNewDB(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if DB == nil {
		fmt.Println("db is nil")
		return
	}
	rows, err := DB.Query("SELECT id, NAME, AGE FROM user")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(id, name, age)
	}
}
