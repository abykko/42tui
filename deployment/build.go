package deployment

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildServerPodman(imageName string, serverDir string) error {

	fmt.Println("Building podman image", imageName)

	cmd := exec.Command(
		"podman", "build",
		"-t", imageName,
		".",
	)

	cmd.Dir = serverDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error building podman image: %w", err)
	}

	return nil
}