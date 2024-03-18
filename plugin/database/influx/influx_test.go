package influx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/machine"
	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	tags := map[string]string{
		"taghf_test1": "taghf_val1",
	}

	for i := 0; i < 1; i++ {
		fields := map[string]interface{}{
			"cpu": machine.GetCPUPercent(),
			"mem": machine.GetMemUsed(),
		}
		fmt.Println(str.ToJSON(fields))
		point := write.NewPoint("hhh", tags, fields, time.Now())
		if err := Cli.WritePoint(context.Background(), point); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second * 5)
	query := `from(bucket: "hufan")
            |> range(start: -10m)
            |> filter(fn: (r) => r._measurement == "hhh")`
	results, err := Cli.QueryClient().Query(context.Background(), query)
	if err != nil {
		fmt.Println("query err:", err)
		return
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err = results.Err(); err != nil {
		fmt.Println("query err:", err)
		return
	}

	fmt.Println("success")
}
