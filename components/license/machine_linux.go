package license

func GetUnique() (string, error) {
	///usr/sbin/dmidecode -t system|grep "Serial Number"
	c := exec.Command("/usr/sbin/dmidecode", "-t", "system")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New("获取机器码失败 " + err.Error() + " - " + string(out))
	}
	var sn string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		linStr := strings.TrimSpace(line)
		if !strings.HasPrefix(linStr, "Serial Number") {
			continue
		}
		snL := strings.Split(linStr, ":")
		if len(snL) == 2 {
			sn = snL[1]
		}
	}
	if sn == "" {
		return "", errors.New("未找到 Serial Number")
	}
	return sn, nil
}
