package jaeger

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// GenTraceCTX 生成带 trace信息的上下文，在各级Span创建时，自动创建child span
func GenTraceCTX(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceInfo, &TraceInfo{})
}

// CloneTraceCTX 复制 上下文，获取新的同内容的 TraceInfo
// 使用场景：当并发多个child function需要进行 span创建的时候，先Clone上下文，让并发时创建的child span父子关系不发生错乱
func CloneTraceCTX(ctx context.Context) context.Context {
	ti := GetTraceInfo(ctx)
	return context.WithValue(ctx, traceInfo, &TraceInfo{Span: ti.Span})
}

// Close 关闭后不可使用，同时会flush未上报的数据
func (cli *TraceClient) Close() {
	if cli == nil || cli.Closer == nil {
		return
	}
	cli.Closer.Close()
}

func (cli *TraceClient) StartSpan(ctx context.Context, name string, opts ...Option) opentracing.Span {
	if cli == nil || cli.Tracer == nil {
		return nil
	}
	optParam := &OptParams{}
	for _, opt := range opts {
		opt(optParam)
	}
	tags := make([]opentracing.StartSpanOption, 0)
	for k, v := range optParam.tag {
		tags = append(tags, opentracing.Tag{
			Key:   k,
			Value: v,
		})
	}
	parentSpan := GetTraceInfo(ctx)
	if parentSpan != nil && parentSpan.Span != nil {
		tags = append(tags, opentracing.ChildOf(parentSpan.Span.Context()))
	}
	span := cli.Tracer.StartSpan(name, tags...)
	SpanLog(span, optParam.logs)
	SetTraceInfo(ctx, span)
	return span
}

func SpanFinish(span opentracing.Span) {
	if span != nil {
		span.Finish()
	}
}

func SetTraceInfo(ctx context.Context, parentSpan opentracing.Span) {
	v, ok := ctx.Value(traceInfo).(*TraceInfo)
	if !ok {
		return
	}
	v.Span = parentSpan
}

func GetTraceInfo(ctx context.Context) *TraceInfo {
	v, ok := ctx.Value(traceInfo).(*TraceInfo)
	if !ok {
		return nil
	}
	return v
}

func TraceID(span opentracing.Span) string {
	if span == nil {
		return ""
	}
	spanContext, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}
	return spanContext.TraceID().String()
}

func SpanLog(span opentracing.Span, logs map[string]string) {
	if span == nil {
		return
	}
	for k, v := range logs {
		span.LogFields(log.String(k, v))
	}
}

// HttpTransport http调用的时候，client 设置 transport，可以跨服务链路追踪，同 InjectHttpHeader 配合使用
func HttpTransport(trans *http.Transport) *otelhttp.Transport {
	if trans == nil {
		return otelhttp.NewTransport(http.DefaultTransport)
	}
	return otelhttp.NewTransport(trans)
}

// InjectHttpHeader http调用的时候，header构造，可以跨服务链路追踪，同 HttpTransport 配合使用
func InjectHttpHeader(tracer opentracing.Tracer, span opentracing.Span) (http.Header, error) {
	if tracer == nil {
		return nil, fmt.Errorf("tracer is nil")
	}
	if span == nil {
		return nil, fmt.Errorf("span is nil")
	}
	baseHeader := http.Header{}
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(baseHeader))
	if err != nil {
		return nil, fmt.Errorf("inject err:%s", err)
	}
	return baseHeader, nil
}

// StartSpanFromHttpHeader 服务端收到Http请求的时候，从Header里面解析出span，让链路跨服务连接起来
func (cli *TraceClient) StartSpanFromHttpHeader(ctx context.Context, header http.Header, name string, opts ...Option) (opentracing.Span, error) {
	if cli == nil || cli.Tracer == nil {
		return nil, fmt.Errorf("tracer is nil")
	}
	optParam := &OptParams{}
	for _, opt := range opts {
		opt(optParam)
	}
	tags := make([]opentracing.StartSpanOption, 0)
	for k, v := range optParam.tag {
		tags = append(tags, opentracing.Tag{
			Key:   k,
			Value: v,
		})
	}
	if GetTraceInfo(ctx) == nil {
		ctx = GenTraceCTX(ctx)
	}
	parentSpan := GetTraceInfo(ctx)
	if parentSpan.Span != nil {
		tags = append(tags, opentracing.ChildOf(parentSpan.Span.Context()))
	}
	clientContext, err := cli.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		return nil, fmt.Errorf("extract err:%s", err)
	}
	tags = append(tags, ext.RPCServerOption(clientContext))
	span := cli.Tracer.StartSpan(name, tags...)
	SpanLog(span, optParam.logs)
	SetTraceInfo(ctx, span)
	return span, nil
}
