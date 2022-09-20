package runshell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 运行shell命令
// @command 命令
// @args 参数
func RunShell(command string, args ...string) (string, error) {
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

// 远程机器运行shell命令
// @user 登陆用户
// @host 机器地址
// @passwd 密码
// @sshKey ras_key
// @command 命令
// @args 参数
func RemoteRunShell(user, host, passwd, sshKey, command string, args ...string) (string, error) {
	cfg := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	n := len(sshKey)
	if n > 0 {
		var signer ssh.Signer
		if n < 1000 {
			key, err := ioutil.ReadFile(sshKey)
			if err != nil {
				return "", err
			}
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return "", err
			}
		} else {
			key := []byte(sshKey)
			var err error
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return "", err
			}
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		cfg.Auth = []ssh.AuthMethod{ssh.Password(passwd)}
	}
	sshClient, err := ssh.Dial("tcp", host, cfg)
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

func CopyToRemoteMachine(user, host, passwd, sshKey, remotePath string, content []byte) error {
	cfg := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	n := len(sshKey)
	if n > 0 {
		var signer ssh.Signer
		if n < 1000 {
			key, err := ioutil.ReadFile(sshKey)
			if err != nil {
				return err
			}
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return err
			}
		} else {
			key := []byte(sshKey)
			var err error
			signer, err = ssh.ParsePrivateKey(key)
			if err != nil {
				return err
			}
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		cfg.Auth = []ssh.AuthMethod{ssh.Password(passwd)}
	}
	sshClient, err := ssh.Dial("tcp", host, cfg)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	dstFile, err := sftp.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	origin := bytes.NewReader(content)
	if _, err := dstFile.ReadFrom(origin); err != nil {
		return err
	}
	return nil
}
