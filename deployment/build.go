package deployment

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildServerDocker(port int) error {

	fmt.Println("Building docker server with port", port)

	runCommand := "gunicorn app.main:app --workers=1"

	cmd := exec.Command(
		"docker", "build",
		"-t", "api-app",
		"--build-arg", fmt.Sprintf("SERVER_PORT=%d", port),
		"--build-arg", fmt.Sprintf("RUN_COMMAND=%s", runCommand),
		".",
	)

	// Ver output en tiempo real
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error building docker image: %w", err)
	}

	return nil
}