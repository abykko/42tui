package main

import (
	"os"
	"fmt"
	"strconv"

	"42cli/conf"
	"42cli/server"
	"42cli/deployment"

	"golang.org/x/term"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

type status struct {
	container bool
	api       bool
	session   bool
}

type profile struct {
	name	string
	user	string
}

type model struct {
	selected 	int
	status   	status
	profile		profile
}

func initialModel() model {
	return model{
		selected: 1,
		status: status{
			container: false,
			api:       false,
			session:   false,
		},
		profile: profile{
			name: "",
			user: "",
		},
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
			if m.selected == 2 {
				deployment.Stop()
				m.status.container = false
				return m, nil
			}

		case "enter":
			if m.selected == 2 {
				err := deployment.Build()
				if err != nil {
					fmt.Println(err)
					return m, nil
				}

				_, err = deployment.Run()
				if err != nil {
					fmt.Println(err)
					return m, nil
				}

				m.status.container = true

				resp, statusCode, err := server.DoRequest("/status")
				if err != nil {
					fmt.Println(err)
					return m, nil
				}

				if value := resp["status"]; statusCode != 200 || value != "ok" {
					fmt.Println("Status not responding")
					return m, nil
				}

				m.status.api = true

				resp, statusCode, err = server.DoSignedRequest("/session/expired")
				if err != nil {
					fmt.Println(err)
					return m, nil
				}

				expired, ok := resp["expired"].(bool)
				if !ok {
					fmt.Println("invalid type for expired")
					return m, nil
				}

				if statusCode != 200 || expired {
					fmt.Println("Session is expired")
					return m, nil
				}

				m.status.session = true

				return m, nil
			}

			if m.selected == 3 {
				resp, statusCode, err := server.DoSignedRequest("/status")
				if err != nil {
					fmt.Println(err)
					return m, nil
				}

				value := resp["status"]
				m.profile.name = value.(string)
				fmt.Println(value, statusCode)
				return m, nil
			}
		}
	}

	return m, nil
}

func Temp(m model, tab int) string {

	color := "0"
	if m.selected == tab {
		color = "45"
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(physicalWidth)

	header := headerStyle.Render("si")

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(12)

	return borderStyle.Render(content)
}

func Profile(m model, tab int) string {

	color := "0"
	if m.selected == tab {
		color = "45"
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	textStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(physicalWidth)

	nameLabel := textStyle.Render(m.profile.name)
	username := textStyle.Render("username")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		nameLabel,
		username,
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(22)

	return borderStyle.Render(content)
}

func ServerStatusView(m model, tab int) string {

	color := "0"
	if m.selected == tab {
		color = "45"
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	textStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(physicalWidth)

	containerLabel := "Container -"
	if m.status.container {
		containerLabel = "Container +"
	}
	containerLabel = textStyle.Render(containerLabel)

	apiLabel := "Api -"
	if m.status.api {
		apiLabel = "Api +"
	}
	apiLabel = textStyle.Render(apiLabel)

	sessionLabel := "Session -"
	if m.status.session {
		sessionLabel = "Session +"
	}
	sessionLabel = textStyle.Render(sessionLabel)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		containerLabel,
		apiLabel,
		sessionLabel,
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(22)

	return borderStyle.Render(content)
}

func (m model) View() tea.View {

	content := lipgloss.JoinVertical(lipgloss.Left,
		ServerStatusView(m, 1),
		Temp(m, 2),
		Profile(m, 3),
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
