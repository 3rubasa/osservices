package osservices

import (
	"bytes"
	"fmt"
	"log"
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

	log.Println("Debug: running command: '", cmd.Path, "' with args '", cmd.Args, "'")
	err := cmd.Run()

	if err != nil {
		log.Println("ERROR: Failed to execute command, error: ", err)
		return "", fmt.Errorf("failed to get IP from MAC-address: %v", err)
	}

	if len(stderr.String()) > 0 {
		log.Println("ERROR: Command has been run succesfully, but stderr is not empty: ", stderr.String())
		return "", fmt.Errorf("failed to get IP from MAC-address, stderr not empty: ", stderr.String())
	}

	re := regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\b`)

	ip := re.FindString(stdout.String())

	if len(ip) == 0 {
		log.Println("Debug: IP for MAC ", mac, " not found, stdout: ", stdout.String())
		return "", ErrNotFound
	}

	return ip, nil
}
