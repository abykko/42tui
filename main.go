package main

import (
    "fmt"
    "os"
    // "42cli/conf"
    "42cli/server-deployment"
    "42cli/ui"
    tea "charm.land/bubbletea/v2"
)

func main() {

    defer func() {
		deployment.Stop()
        deployment.StopByImage()
    }()

    p := tea.NewProgram(ui.InitialModel())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v", err)
        os.Exit(1)
    }
}