package deployment

import (
	"fmt"
	"strconv"
	"42cli/conf"
)

func DeployServer() (string, error) {
	// Generate secret key
	secret := GenerateSecret()

	cfg, err := conf.LoadConfig("conf/.conf", false)
	if err != nil {
		return "", fmt.Errorf("error reading settings file: %w", err)
	}

	port, err := strconv.Atoi(cfg["server_port"])
	if err != nil {
		return "", fmt.Errorf("port conversion error: %w", err)
	}

	imageName := cfg["podman_image_name"]
	serverDir := cfg["server_dir"]

	// Build image
	if err := BuildServerPodman(imageName, serverDir); err != nil {
		return "", fmt.Errorf("error building image: %w", err)
	}

	env := cfg["container_id_env_var_name"]
	
	// Run container
	podmanId, err := RunServerPodman(imageName, secret, port, env)
	if err != nil {
		return "", fmt.Errorf("error running podman server: %w", err)
	}

	defer func() {
		if err != nil {
			_ = StopServerPodman(podmanId)
		}
	}()

	return podmanId, nil
}
