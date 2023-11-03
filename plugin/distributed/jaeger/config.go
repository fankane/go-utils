package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"io"
)

const (
	traceInfo         = "jaeger_trace_info"
	defaultIntervalMs = 1000
)

type Config struct {
	ServiceName           string  `yaml:"service_name"`
	SamplerType           string  `yaml:"sampler_type"`
	SamplerParam          float64 `yaml:"sampler_param"`
	LogSpans              bool    `yaml:"log_spans"`
	CollectorEndpoint     string  `yaml:"collector_endpoint"`
	User                  string  `yaml:"user"`
	Password              string  `yaml:"password"`
	BufferFlushIntervalMs int64   `json:"buffer_flush_interval_ms"` //强制刷新时间间隔：毫秒
}

type TraceInfo struct {
	Span opentracing.Span
}

type TraceClient struct {
	Tracer opentracing.Tracer
	Closer io.Closer
}

type OptParams struct {
	tag  map[string]interface{}
	logs map[string]string
}

type Option func(params *OptParams)

func Tags(tag map[string]interface{}) Option {
	return func(params *OptParams) {
		params.tag = tag
	}
}

func Logs(logs map[string]string) Option {
	return func(params *OptParams) {
		params.logs = logs
	}
}
