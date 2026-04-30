package main

import (
	"fmt"
	"os"
	"42cli/conf"
	"42cli/deployment"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/term"
)

type s_loading struct {
	header		string
	subheader	string
	status		string
}

type model struct{
	loading 		s_loading
	activeWindow	string
	serverRunning	bool
}

func initialModel() model {
	return model{
		loading: s_loading{
			header: "42cli",
			subheader: "by iamrani-",
			status: "Press Enter to start...",
		},
		serverRunning: false,
		activeWindow: "loadingScreen",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			if m.serverRunning == true {
				m.loading.status = "Stopping server"
			}
			return m, tea.Quit

		case "enter":
			if m.serverRunning == false {
				_, err := deployment.DeployServer()
				if err != nil {
					fmt.Println("Error:", err)
				}
				m.serverRunning = true
				m.loading.status = "Server running"
			}
		}
	}
	return m, nil
}

func LoadingScreen(m model) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(physicalWidth)
	header := headerStyle.Render(m.loading.header)

	subheader := headerStyle.Render(m.loading.subheader)

	status := headerStyle.Render(m.loading.status)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		subheader,
		status,
	)

	return content
}

func (m model) View() tea.View {

	if m.activeWindow == "loadingScreen" {
		v := tea.NewView(LoadingScreen(m))
		return v
	}
	
	v := tea.NewView("")
	v.AltScreen = true
	return v
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

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running UI:", err)
		os.Exit(1)
	}
}