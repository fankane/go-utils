package prometheus

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	g := GetGaugeVec("test1")
	c := GetCounterVec("http_num")
	h := GetHistogram("yyy")
	if h == nil {
		fmt.Println("h is nil")
		return
	}
	s := GetSummary("summary_test1")
	if h == nil {
		fmt.Println("s is nil")
		return
	}
	fmt.Println(RegisteredCollTypeList())
	fmt.Println(RegisteredCollNameList(CollCounter))
	for i := 0; i < 1000; i++ {
		g.WithLabelValues("val1", "val2").Set(121.0)
		g.WithLabelValues("val1", "val3").Set(122.0)
		g.WithLabelValues("val2", "val3").Set(123.0)

		c.WithLabelValues("cnt1", "http").Inc()
		h.Observe(66.66)
		s.Observe(77.77)
		time.Sleep(time.Millisecond * 50)
	}

	time.Sleep(time.Hour)
}
