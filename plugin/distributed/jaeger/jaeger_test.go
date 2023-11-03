package jaeger

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	if Tracer == nil {
		fmt.Println("tracer is nil")
		return
	}
	//defer Tracer.Closer.Close()
	Root()
	time.Sleep(time.Second)
}

func Root() {
	ctx := GenTraceCTX(context.Background())
	span := Tracer.StartSpan(ctx, "root", Tags(map[string]interface{}{"tag1": "test1"}))
	defer SpanFinish(span)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		A1(CloneTraceCTX(ctx)) //并发处理时 Clone ctx
	}()
	go func() {
		defer wg.Done()
		A2(CloneTraceCTX(ctx))
	}()
	wg.Wait()
	fmt.Println("done, traceID:", TraceID(span))
}

func A1(ctx context.Context) {
	span := Tracer.StartSpan(ctx, "A1", Tags(map[string]interface{}{"tagA1": "testA1"}))
	defer SpanFinish(span)
	time.Sleep(time.Millisecond * 600)
	B(ctx)
}

func A2(ctx context.Context) {
	span := Tracer.StartSpan(ctx, "A2", Tags(map[string]interface{}{"tagA2": "test2"}))
	defer SpanFinish(span)
	time.Sleep(time.Millisecond * 600)
	D(ctx)
}

func B(ctx context.Context) {
	span := Tracer.StartSpan(ctx, "B", Logs(map[string]string{"log1": "log 001"}))
	defer SpanFinish(span)
	time.Sleep(time.Millisecond * 500)
	D(ctx)
}

func C(ctx context.Context) {
	span := Tracer.StartSpan(ctx, "C", Logs(map[string]string{"log1": "log 002"}), Tags(map[string]interface{}{"tagA2": "test2"}))
	defer SpanFinish(span)
	time.Sleep(time.Millisecond * 400)
}

func D(ctx context.Context) {
	span := Tracer.StartSpan(ctx, "D", Logs(map[string]string{"log1": "log 001"}))
	defer SpanFinish(span)
	time.Sleep(time.Millisecond * 200)
	C(ctx)
}
