package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func WaitForShutdown(hooks ...Hook) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case sig := <-signalChan:
		switch sig {
		case syscall.SIGHUP:
			fmt.Printf("received signal SIGHUP\n")
		case syscall.SIGINT:
			fmt.Printf("received signal SIGINT\n")
		case syscall.SIGQUIT:
			fmt.Printf("received signal SIGQUIT\n")
		}
		fmt.Printf("shutting down\n")
		time.AfterFunc(time.Minute, func() {
			os.Exit(0)
		})

		for _, hook := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			err := hook(ctx)
			if err != nil {
				fmt.Printf("failed to run hook: %v\n", err)
			}
			cancel()
		}
	}
	os.Exit(0)
}

type GracefulShutdown struct {
	closing    int32
	reqCnt     int64
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c Context) {
		// 开始拒绝所有的请求
		cl := atomic.LoadInt32(&g.closing)
		if cl > 0 {
			c.W().WriteHeader(http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		n := atomic.AddInt64(&g.reqCnt, -1)
		// 已经开始关闭了，而且请求数为0，
		if cl > 0 && n == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	atomic.AddInt32(&g.closing, 1)
	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}
	done := ctx.Done()
	select {
	case <-done:
		fmt.Println("timeout")
		return errors.New("timeout")
	case <-g.zeroReqCnt:
		fmt.Println("all requests finished")
	}
	return nil
}
