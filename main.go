package main

import (
	"os"
	"fmt"
	"strconv"
	"42cli/conf"
	// "42cli/server"
	"42cli/deployment"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

type s_loading struct {
	header		string
	status		string
}

type s_status struct {
	container	bool
	api			bool
	session		bool
}

type model struct{
	loading 		s_loading
	status			s_status
	selected		int
}

func initialModel() model {
	return model{
		loading: s_loading{header: "42cli",status: "Press Enter to start...",},
		status: s_status{container: false, api: false, session: false,},
		selected: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "1", "2", "3":
			m.selected, _ = strconv.Atoi(msg.String())
			return m, nil

		case "q":
			m.status.container = false
			m.loading.status = "Server stopped."
			deployment.Stop()
			return m, nil

		case "enter":
			m.status.container = true
			m.loading.status = "Server running"
			deployment.Build()
			deployment.Run()
			return m, nil
		}
	}

	return m, nil
}

func Temp(m model) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(physicalWidth)
	header := headerStyle.Render(m.loading.header)

	status := headerStyle.Render(m.loading.status)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		status,
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder())
	
	return borderStyle.Render(content)
}

func ServerStatusView(m model) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	textStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(physicalWidth)
	
	containerLabel := textStyle.Render("container")
	apiLabel := textStyle.Render("api")
	sessionLabel := textStyle.Render("session")

	content := lipgloss.JoinVertical(lipgloss.Center,
		containerLabel,
		apiLabel,
		sessionLabel,
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder())
	
	return borderStyle.Render(content)
}

func (m model) View() tea.View {

	content := lipgloss.JoinVertical(lipgloss.Center,
		Temp(m),
		ServerStatusView(m),
	)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func main() {

	envVarName, err := conf.GetString("container_id_env_var_name")
	if err != nil {
		fmt.Println("error obtaining the name of the env variable with the container id")
		os.Exit(1)
	}

	defer func() {
		containerID := os.Getenv(envVarName)
		if containerID != "" {
			fmt.Println("Stopping container from main() defer function. Container:", containerID)
			if err := deployment.Stop(); err != nil {
				fmt.Println("Error stopping container:", err)
			}
		}
		os.Exit(1)
	}()

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running UI:", err)
		os.Exit(1)
	}
}
