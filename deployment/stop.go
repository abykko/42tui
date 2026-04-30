package deployment

import (
	"fmt"
	"os/exec"
)

func StopServerPodman(container string) error {

	cmd := exec.Command(
		"podman", "kill", container,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error stopping container: %w - %s", err, string(output))
	}

	return nil
}