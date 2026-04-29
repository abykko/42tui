package deployment

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunServerPodman(imageName string, secret string, port int) (string, error) {

	fmt.Println("Running podman image:", imageName)

	cmd := exec.Command(
		"podman", "run",
		"-d",
		"-e", fmt.Sprintf("SERVER_PORT=%d", port),
		"-e", fmt.Sprintf("SECRET=%s", secret),
		fmt.Sprintf("%s:latest", imageName),
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running podman container: %w", err)
	}

	containerID := strings.TrimSpace(string(output))

	fmt.Println("Container started with ID:", containerID)

	return containerID, nil
}