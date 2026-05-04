package deployment

import (
    "fmt"
    "os"
    "os/exec"
	
    "42cli/conf"
)

func Stop() error {
    contEnvId, err := conf.GetString("container_id_env_var_name")
    if err != nil {
        return fmt.Errorf("error getting container env var name: %w", err)
    }

    containerId := os.Getenv(contEnvId)
    if containerId == "" {
        return nil 
    }

    killCmd := exec.Command("podman", "kill", containerId)
    if err := killCmd.Run(); err != nil {
        // Solo imprimimos aviso, ya que si el contenedor ya estaba parado, kill fallará
        fmt.Printf("Aviso: No se pudo ejecutar kill (quizás ya estaba detenido): %v\n", err)
    }

	rmCmd := exec.Command("podman", "rm", "-f", containerId)
    output, err := rmCmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("error removing container: %w - %s", err, string(output))
    }

    if err := os.Unsetenv(contEnvId); err != nil {
        return fmt.Errorf("error unsetting env var: %w", err)
    }

    return nil
}