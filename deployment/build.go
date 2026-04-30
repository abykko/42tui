package deployment

import (
	"fmt"
	"os/exec"
)

func BuildServerPodman(imageName string, serverDir string) error {

	cmd := exec.Command(
		"podman", "build",
		"-t", imageName,
		".",
	)

	cmd.Dir = serverDir
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error building podman image: %w", err)
	}

	return nil
}
