package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(callFunc ...func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	switch <-ch {
	case syscall.SIGKILL:
		fmt.Printf("信号监听:%d\n", syscall.SIGKILL)
		Quit(callFunc...)

	case syscall.SIGINT:
		fmt.Printf("信号监听:%d\n", syscall.SIGINT)
		Quit(callFunc...)

	case syscall.SIGTERM:
		fmt.Printf("信号监听:%d\n", syscall.SIGTERM)
		Quit(callFunc...)

	case syscall.SIGQUIT:
		fmt.Printf("信号监听:%d\n", syscall.SIGQUIT)
		Quit(callFunc...)
	}
}

func Quit(callFunc ...func()) {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	for _, item := range callFunc {
		item()
	}
	fmt.Println("退出成功...")
	os.Exit(0)
}
