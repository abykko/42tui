package deployment

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunServerPodman(imageName string, secret string, port int, env string) (string, error) {

	fmt.Println("Running podman image:", imageName)

	cmd := exec.Command(
		"podman", "run",
		"-d",
		"-p", fmt.Sprintf("127.0.0.1:%d:%d", port, port),
		"-e", fmt.Sprintf("SERVER_PORT=%d", port),
		"-e", fmt.Sprintf("SECRET=%s", secret),
		fmt.Sprintf("%s:latest", imageName),
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running podman container: %w", err)
	}

	containerID := strings.TrimSpace(string(output))

	// Save the container Id in a env var
	err = os.Setenv(env, containerID)
	if err != nil {
		return "", fmt.Errorf("error setting env var: %w", err)
	}

	fmt.Println("Container started with ID:", containerID)

	return containerID, nil
}