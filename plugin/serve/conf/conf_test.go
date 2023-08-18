package conf

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
	"testing"
	"time"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	x := &TT{}
	if err := Unmarshal(x); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(str.ToJSON(x))

	utime.TickerDo(time.Second*2, func() error {
		fmt.Println(time.Now(), str.ToJSON(x))
		return nil
	})

	time.Sleep(time.Minute)
}

type AB struct {
	A int    `yaml:"a"`
	B string `yaml:"b"`
}

type TT struct {
	ResultPath  string     `yaml:"result_path" json:"result_path"`
	ProjectName string     `yaml:"project_name" json:"project_name"`
	Menus       []MenuTree `yaml:"menus" json:"menus"`
}

type MenuTree struct {
	Name      string     `yaml:"name" json:"name"`             //目录名
	FileNames []string   `yaml:"file_names" json:"file_names"` //目录下的文件名
	Children  []MenuTree `yaml:"children" json:"children"`     //子目录
}
