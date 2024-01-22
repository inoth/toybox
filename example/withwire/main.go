package main

import (
	"fmt"
	"os"
)

func main() {
	tb := initApp()
	if err := tb.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
