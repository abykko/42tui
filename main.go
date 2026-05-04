package main

import (
    "fmt"
    "os"
    "42cli/conf"
    "42cli/server-deployment"
    "42cli/ui"
    tea "charm.land/bubbletea/v2"
)

func main() {
    envVarName, err := conf.GetString("container_id_env_var_name")
    if err != nil {
        fmt.Println("Error config")
        os.Exit(1)
    }

    defer func() {
        containerID := os.Getenv(envVarName)
        if containerID != "" {
            deployment.Stop()
        }
    }()

    p := tea.NewProgram(ui.InitialModel())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v", err)
        os.Exit(1)
    }
}