package internal

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
		fmt.Printf("Listen:%d\n", syscall.SIGKILL)
		Quit(callFunc...)

	case syscall.SIGINT:
		fmt.Printf("Listen:%d\n", syscall.SIGINT)
		Quit(callFunc...)

	case syscall.SIGTERM:
		fmt.Printf("Listen:%d\n", syscall.SIGTERM)
		Quit(callFunc...)

	case syscall.SIGQUIT:
		fmt.Printf("Listen:%d\n", syscall.SIGQUIT)
		Quit(callFunc...)
	}
}

func Quit(callFunc ...func()) {
	fmt.Println("Starting to quit...")
	fmt.Println("Executing cleanup...")
	for _, item := range callFunc {
		item()
	}
	fmt.Println("Successfully quit...")
	os.Exit(0)
}
