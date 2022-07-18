package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	// 服務器設定
	srv := &http.Server{Addr: ":8081"}

	// 建立父Ctx
	parentCtx := context.Background()

	// 透過父Ctx產生帶有Cancel的Ctx
	cancelCtx, cancel := context.WithCancel(parentCtx)

	// 註冊errGroup
	eg, ctx := errgroup.WithContext(cancelCtx)
	eg.Go(func() error {
		return startServer(srv)
	})

	eg.Go(func() error {
		return closeServer(srv, ctx)
	})

	// 宣告buffer channel為接收信號所使用
	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	// 開始做信號接收
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-chanel:
				cancel()
			}
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Println("get error: ", err)
	}
	fmt.Println("finish")

}

// startServer 啟動服務器
func startServer(srv *http.Server) error {
	http.HandleFunc("/hello", helloHandle)
	fmt.Println("server start!")
	fmt.Println("hello api url: http://127.0.0.1:8081/hello")
	return srv.ListenAndServe()
}

// closeServer 關閉服務器
func closeServer(srv *http.Server, ctx context.Context) error {
	<-ctx.Done()
	return srv.Shutdown(ctx)
}

// helloHandle hello API
func helloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
