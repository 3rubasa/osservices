package osservices

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
)

type OSServices struct {
}

func NewOSServicesProvider() *OSServices {
	return &OSServices{}
}

func (o OSServices) Reboot() error {
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func (o OSServices) GetIPFromMAC(mac string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command := fmt.Sprintf("arp | grep %s", mac)
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("Failed to exec command: ", err)
		return "", err
	}

	if len(stderr.String()) > 0 {
		fmt.Println("STDERR: ", stderr.String())
		return "", fmt.Errorf("error, STDERR not empty: %s", stderr.String())
	}

	re := regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\b`)

	ip := re.FindString(stdout.String())

	if len(ip) == 0 {
		fmt.Println("STDOUT: ", stdout.String())
		return "", fmt.Errorf("error, IP address not found, STDOUT: %s", stdout.String())
	}

	return ip, nil
}
