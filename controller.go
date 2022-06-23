package main

import (
	"context"
	"fmt"
	"geektime-web/framework"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(10*time.Second))
	defer cancel()

	//finish负责通知结束
	finish := make(chan struct{}, 1)
	//panicChan负责通知panic异常
	panicChan := make(chan interface{}, 1)

	go func() {
		//异常处理
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		//具体业务
		time.Sleep(5 * time.Second)
		c.Json(200, "ok")
		//结束时通过一个finish通道告知父goroutine
		finish <- struct{}{}
	}()
	select {
	case <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}
