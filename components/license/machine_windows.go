package license

func GetUnique() (string, error) {
	c := exec.Command("cmd", "/C", "wmic diskdrive get serialnumber")
	output, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.Replace(string(output), "\n", "", -1), nil
}
