package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
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
