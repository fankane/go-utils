package prometheus

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
	"time"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_hf",
	}, []string{"label1", "label2"})
	prometheus.MustRegister(g)
	for i := 0; i < 100; i++ {
		g.WithLabelValues("val1", "val2").Set(1.0)
		g.WithLabelValues("val1", "val3").Set(2.0)
		g.WithLabelValues("val2", "val3").Set(3.0)
		time.Sleep(time.Millisecond * 50)
	}

	time.Sleep(time.Hour)
}
