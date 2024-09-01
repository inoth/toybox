package file

import (
	"context"
	"fmt"
	"testing"
)

func TestTailFile(t *testing.T) {
	lines := make(chan string)

	go TailFile(context.Background(), "log/info.log", lines)
	i := 0
	for line := range lines {
		i++
		fmt.Println(i, line)
	}
	t.Log("ok")
}
