/*!
 * <信号监听>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package profile

import (
	pub "common/public_yfm"
	"os"
	"os/signal"
	"syscall"
)

type SignalBackF func()

func Signal(calback SignalBackF) {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) //, syscall.SIGUSR1, syscall.SIGUSR2
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				pub.PrintLog("进程退出 :", s)
				calback()
			// case syscall.SIGUSR1:
			// 	pub.PrintLog("usr1", s)
			// case syscall.SIGUSR2:
			// 	pub.PrintLog("usr2", s)
			default:
				pub.PrintLog("other", s)
			}
		}
	}()
}
