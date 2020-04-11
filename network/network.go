package network

import "os/exec"

func GetIP() string {
	output, _ := exec.Command("sh", "-c", "ifconfig wlan0 | grep 'inet addr' | cut -d: -f2 | awk '{print $1}' | tr -d '\n'").Output()

	return string(output)
}
