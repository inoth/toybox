package license

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func GetUnique() (string, error) {
	cmd := exec.Command("bash", "-c", "system_profiler SPHardwareDataType | awk '/Serial/ {print $4}'")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取机器码失败, err: %s, output: %s", err, output)
	}

	sn := strings.TrimSpace(string(output))
	if sn == "" {
		return "", errors.New("未找到 Serial Number")
	}
	return sn, nil
}
