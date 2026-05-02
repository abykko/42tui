package deployment

import (
	"fmt"
	"os"
	"os/exec"

	"42cli/conf"
)

func Stop() error {

	// Get the running container id
	contEnvId, err := conf.GetString("container_id_env_var_name")
	if err != nil {
		return fmt.Errorf("error getting container env var name: %w", err)
	}

	containerId := os.Getenv(contEnvId)
	if containerId == "" {
		return fmt.Errorf("container id env var is empty")
	}

	// Stop
	cmd := exec.Command("podman", "kill", containerId)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error stopping container: %w - %s", err, string(output))
	}

	// Leave the env variable empty
	err = os.Unsetenv(contEnvId)
	if err != nil {
		return fmt.Errorf("error unsetting env var: %w", err)
	}

	return nil
}