package deployment

import (
	"fmt"
	"os/exec"
)

func StopServerPodman(container string) error {

	fmt.Println("Stopping podman container:", container)

	cmd := exec.Command(
		"podman", "kill", container,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error stopping container: %w - %s", err, string(output))
	}

	fmt.Println("Container stopped:", string(output))

	return nil
}