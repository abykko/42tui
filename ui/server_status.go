package ui

import (
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func ServerStatusView(m Model, tab int) string {

	color := "0"
	if m.Selected == tab {
		color = "45"
	}

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(w)

	container := "Container -"
	if m.Status.Container {
		container = "Container +"
	}

	api := "Api -"
	if m.Status.Api {
		api = "Api +"
	}

	session := "Session -"
	if m.Status.Session {
		session = "Session +"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		style.Render(container),
		style.Render(api),
		style.Render(session),
	)

	border := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(22)

	return border.Render(content)
}