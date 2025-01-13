package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/plugin/log"
	"github.com/fankane/go-utils/str"
	"github.com/gin-gonic/gin"
)

const CtxTrace = "_ctx_trace"

type ListenEngine struct {
	Addr            string
	Engine          *gin.Engine
	Signal          []os.Signal   //需要监听的信号量
	ShutdownTimeout time.Duration //停止服务时，最长等待时间，不填 默认1分钟
}

// StartEngine 启动gin engine，收到指定信号量则停止服务
func StartEngine(param *ListenEngine) error {
	if param == nil || param.Engine == nil {
		return fmt.Errorf("gin engine is nil")
	}

	signalChan := make(chan os.Signal, 1)
	hasSignal := false
	if len(param.Signal) > 0 {
		// 创建信号通道，用于捕获系统信号
		signal.Notify(signalChan, param.Signal...)
		hasSignal = true
	}
	// 创建一个自定义的 HTTP 服务器
	srv := &http.Server{
		Addr:    param.Addr,
		Handler: param.Engine,
	}

	if !hasSignal { //不监听信息，直接启动
		return srv.ListenAndServe()
	}

	// 启动 HTTP 服务器在一个新的 Goroutine 中
	go func() {
		if errRun := srv.ListenAndServe(); errRun != nil && errRun != http.ErrServerClosed {
			log.Logger.Panicf("Failed to start server: %v", errRun)
		}
	}()

	sg := &sync.WaitGroup{}
	sg.Add(1)
	go func() {
		defer goroutine.Recover()
		sig := <-signalChan
		log.Logger2.Infof("received signal: %s", sig)
		sg.Done()
	}()
	sg.Wait()
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if param.ShutdownTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), param.ShutdownTimeout)
	}
	srv.Shutdown(ctx)
	return nil
}

func MidSetTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := strings.ReplaceAll(str.UUID(), "-", "")
		if len(traceID) > 12 {
			traceID = traceID[0:12] //取uuid前12位
		}
		c.Set(CtxTrace, traceID)
		c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
	traceID := c.GetString(CtxTrace)
	if traceID == "" {
		traceID = str.UUID()
	}
	return traceID
}
