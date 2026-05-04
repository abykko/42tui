package ui

import (
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func ProfileView(m Model, tab int) string {

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

	name := style.Render(m.Profile.Name)
	user := style.Render("username")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		name,
		user,
	)

	border := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(22)

	return border.Render(content)
}