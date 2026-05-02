package deployment

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"42cli/conf"
)

func Run() (string, error) {
	// Load config
	port, err := conf.GetString("server_port")
	if err != nil {
		return "", fmt.Errorf("getting server_port: %w", err)
	}

	imageName, err := conf.GetString("podman_image_name")
	if err != nil {
		return "", fmt.Errorf("getting podman_image_name: %w", err)
	}

	containerIDEnv, err := conf.GetString("container_id_env_var_name")
	if err != nil {
		return "", fmt.Errorf("getting container_id_env_var_name: %w", err)
	}

	secretEnv, err := conf.GetString("secret_env_var_name")
	if err != nil {
		return "", fmt.Errorf("getting secret_env_var_name: %w", err)
	}

	secret := os.Getenv(secretEnv)
	if secret == "" {
		return "", fmt.Errorf("environment variable %q is empty", secretEnv)
	}

	// Build command
	cmd := exec.Command(
		"podman", "run",
		"-d",
		"-p", fmt.Sprintf("127.0.0.1:%s:%s", port, port),
		"-e", fmt.Sprintf("SERVER_PORT=%s", port),
		"-e", fmt.Sprintf("SECRET=%s", secret),
		fmt.Sprintf("%s:latest", imageName),
	)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("running podman container: %w - %s", err, strings.TrimSpace(string(output)))
	}

	containerID := strings.TrimSpace(string(output))

	// Save container ID in env var
	if err := os.Setenv(containerIDEnv, containerID); err != nil {
		return "", fmt.Errorf("setting env var %q: %w", containerIDEnv, err)
	}

	return containerID, nil
}