package log

import (
	"context"
	"strings"

	"github.com/fankane/go-utils/str"
	"go.uber.org/zap"
)

const ctxTrace = "_ctx_trace"

type Log struct {
	log *zap.SugaredLogger
}

func NewTraceCtx() context.Context {
	traceID := strings.ReplaceAll(str.UUID(), "-", "")
	if len(traceID) > 12 {
		traceID = traceID[0:12] //取uuid前12位
	}
	return context.WithValue(context.Background(), ctxTrace, traceID)
}

func getTrace(ctx context.Context) string {
	t := ctx.Value(ctxTrace)
	if t == nil {
		return ""
	}
	return t.(string)
}

func (l *Log) Debugf(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Debugf(template, args...)
}
func (l *Log) Infof(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Infof(template, args...)
}
func (l *Log) Warnf(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Warnf(template, args...)
}
func (l *Log) Errorf(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Errorf(template, args...)
}
func (l *Log) Panicf(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Panicf(template, args...)
}
func (l *Log) Fatalf(template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	l.log.Fatalf(template, args...)
}

func (l *Log) DebugfCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Debugf(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Debugf(template, args...)
}

func (l *Log) InfofCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Infof(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Infof(template, args...)
}
func (l *Log) WarnfCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Warnf(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Warnf(template, args...)
}
func (l *Log) ErrorfCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Errorf(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Errorf(template, args...)
}
func (l *Log) PanicfCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Panicf(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Panicf(template, args...)
}
func (l *Log) FatalfCtx(ctx context.Context, template string, args ...interface{}) {
	if l == nil || l.log == nil {
		return
	}
	traceVal := getTrace(ctx)
	if traceVal == "" {
		l.log.WithOptions(zap.AddCallerSkip(1)).Fatalf(template, args...)
		return
	}
	l.log.WithOptions(zap.AddCallerSkip(1), zap.Fields(zap.String(ctxTrace, getTrace(ctx)))).Fatalf(template, args...)
}