package shelltool

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/inoth/toybox/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ShellCmd struct {
	Host   string // 地址
	Port   string // 端口；default 22
	User   string // 用户
	Passwd string // 密码
	IdRsa  string // .ssh 下密钥, len<1000 为地址，len>1000 为密钥
}

func NewShellCmd(host, user, passwd string, idRsa ...string) *ShellCmd {
	return &ShellCmd{
		Host:   host,
		Port:   ":22",
		User:   user,
		Passwd: passwd,
		IdRsa:  utils.FirstParam("", idRsa),
	}
}

// 本地命令执行
func LocalCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var buf []byte
	res := bytes.NewBuffer(buf)

	cmd.Stdout = io.MultiWriter(os.Stdout, res)
	cmd.Stderr = io.MultiWriter(os.Stderr, res)

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

// 远程命令执行
func (sc *ShellCmd) RemoteCmd(command string, args ...string) (string, error) {
	cfg := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sc.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	n := len(sc.IdRsa)
	if n > 0 {
		var signer ssh.Signer
		if n < 1000 {
			key, err := os.ReadFile(sc.IdRsa)
			if err != nil {
				return "", err
			}
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return "", err
			}
		} else {
			key := []byte(sc.IdRsa)
			var err error
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return "", err
			}
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		cfg.Auth = []ssh.AuthMethod{ssh.Password(sc.Passwd)}
	}
	sshClient, err := ssh.Dial("tcp", sc.Host+sc.Port, cfg)
	if err != nil {
		return "", err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	command = fmt.Sprintf("%v %v", command, strings.Join(args, " "))
	// combo, err := session.CombinedOutput(command)
	var buf []byte
	res := bytes.NewBuffer(buf)
	session.Stdout = io.MultiWriter(os.Stdout, res)
	session.Stderr = io.MultiWriter(os.Stderr, res)

	err = session.Run(command)
	if err != nil {
		return "", err
	}
	return string(res.String()), nil
}

// 发送内存内容到远程机器
func (sc *ShellCmd) RemoteCopyStream(fileName, remotepath string, content []byte) error {
	cfg := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sc.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	n := len(sc.IdRsa)
	if n > 0 {
		var signer ssh.Signer
		if n < 1000 {
			key, err := os.ReadFile(sc.IdRsa)
			if err != nil {
				return err
			}
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return err
			}
		} else {
			key := []byte(sc.IdRsa)
			var err error
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return err
			}
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		cfg.Auth = []ssh.AuthMethod{ssh.Password(sc.Passwd)}
	}
	sshClient, err := ssh.Dial("tcp", sc.Host+sc.Port, cfg)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer client.Close()

	client.MkdirAll(remotepath)
	// leave your mark
	remoteFilePath := fmt.Sprintf("%v/%v", remotepath, fileName)
	f, err := client.Create(remoteFilePath)
	if err != nil {
		return err
	}
	_ = client.Chmod(remoteFilePath, 0655)
	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}
	f.Close()

	// check it's there
	fi, err := client.Lstat(remoteFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("目标[%v]写入完成;[%v]文件大小: %v\n", sc.Host, remoteFilePath, fi.Size())
	return nil
}

// 从文件获取内容发送到远程
func (sc *ShellCmd) RemoteCopyFile(orginfile, remotepath string) error {
	paths := strings.Split(orginfile, "/")
	fileName := paths[len(paths)-1]
	content, err := os.ReadFile(orginfile)
	if err != nil {
		return err
	}
	return sc.RemoteCopyStream(fileName, remotepath, content)
}
