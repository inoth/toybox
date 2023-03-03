package test

import (
	"fmt"
	shell "github/inoth/ino-toybox/tools/shell_tool"
	"testing"
)

func TestRemoteCmd(t *testing.T) {
	sc := shell.NewShellCmd("192.168.1.100", "ubuntu", "12345678")
	result, err := sc.RemoteCmd("ls", "-l", "~/")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	t.Log()
}

func TestRemoteCopyStream(t *testing.T) {
	sc := shell.NewShellCmd("192.168.1.100", "ubuntu", "12345678")
	err := sc.RemoteCopyStream("test.txt", "/var/script/test", []byte("test txt by bytes"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
	t.Log()
}

func TestRemoteCopyFile(t *testing.T) {
	sc := shell.NewShellCmd("192.168.1.100", "ubuntu", "12345678")
	err := sc.RemoteCopyFile("test.txt", "/var/script/test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
	t.Log()
}
