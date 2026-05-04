package ui

import (
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func TempView(m Model, tab int) string {

	color := "0"
	if m.Selected == tab {
		color = "45"
	}

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	header := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(w).
		Render("si")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
	)

	border := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		Width(12)

	return border.Render(content)
}