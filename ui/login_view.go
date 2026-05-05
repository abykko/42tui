package ui

import (
	"charm.land/lipgloss/v2"
)

func LoginView(m Model) string {

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	label := lipgloss.NewStyle().Bold(true)

	button := lipgloss.NewStyle().
		Padding(0, 3).
		MarginTop(1).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("7"))

	if m.Login.FocusIndex == 2 {
		button = button.
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("15"))
	}

	user := lipgloss.JoinVertical(lipgloss.Left,
		label.Render("Username:"),
		m.Login.UsernameInput.View(),
	)

	pass := lipgloss.JoinVertical(lipgloss.Left,
		label.Render("Password:"),
		m.Login.PasswordInput.View(),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title.Render("SYSTEM AUTH"),
		user,
		pass,
		button.Render("LOGIN"),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("23")).
		Padding(1, 4).
		Width(30).
		Render(content)

	return box
}
