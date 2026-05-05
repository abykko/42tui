package deployment

import (
	"os"
	"fmt"
	"os/exec"
	"crypto/rand"
	"encoding/hex"

	"42cli/conf"
)

func generateSecret() string {
	b := make([]byte, 32) // 256 bits

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func Build() error {

	secret := generateSecret()

	secretEnv, err := conf.GetString("secret_env_var_name")
	if err != nil {
		return fmt.Errorf("error getting secret env name: %w", err)
	}

	// We store it to use the api
	if err := os.Setenv(secretEnv, secret); err != nil {
		return fmt.Errorf("error setting SECRET env var: %w", err)
	}

	// Load config
	servDir, err := conf.GetString("server_dir")
	if err != nil {
		return fmt.Errorf("error getting server_dir: %w", err)
	}

	imageName, err := conf.GetString("podman_image_name")
	if err != nil {
		return fmt.Errorf("error getting podman_image_name: %w", err)
	}

	// Remove the containers if already exists for any reason
	StopByImage()

	cmd := exec.Command("podman", "build", "-t", imageName, ".")
	cmd.Dir = servDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error building podman image: %w", err)
	}

	return nil
}
