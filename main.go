package main

import (
	"fmt"
	"os"
	"42cli/conf"
	"42cli/deployment"
	"42cli/components"
	tea "charm.land/bubbletea/v2"
)

type model struct{
	loadingPage components.LoadingPage
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			// _, err := deployment.DeployServer()
			// if err != nil {
			// 	fmt.Println("Error:", err)
			// }
			m.loadingPage.Text = "hola"
		}
	}
	return m, nil
}

func (m model) View() tea.View	 {
	return tea.NewView(fmt.Sprintf("Hello %s", m.loadingPage.Text))
}

func main() {

	cfg, err := conf.LoadConfig("conf/.conf", false)
	if err != nil {
		fmt.Println("error reading settings file:", err)
		os.Exit(1)
	}

	envVarName := cfg["container_id_env_var_name"]

	defer func() {
		containerID := os.Getenv(envVarName)
		if containerID != "" {
			fmt.Println("Stopping container from env:", containerID)
			if err := deployment.StopServerPodman(containerID); err != nil {
				fmt.Println("Error stopping container:", err)
			}
		}
	}()

	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running UI:", err)
		os.Exit(1)
	}
}