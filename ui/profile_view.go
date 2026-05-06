package ui

import (
	"fmt"
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func profile(m Model) string {

	profile := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2, 4).
		Width(50)

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	usernameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	name := fmt.Sprintf("%s %s", m.Profile.FirstName, m.Profile.LastName)
	username := fmt.Sprintf("@%s", m.Profile.Login)
	email := m.Profile.Email
	location := "Málaga(MOCKDATA)"
	bio := "Building cool CLI apps with Go 🚀"

	avatar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("69")).
		Render(
			" ◉ ◉ \n" +
				"  ◉  \n" +
				" ◉ ◉ ",
		)

	info := lipgloss.JoinVertical(
		lipgloss.Left,
		nameStyle.Render(name),
		usernameStyle.Render(username),
		"",
		labelStyle.Render("Email: ")+valueStyle.Render(email),
		labelStyle.Render("Location: ")+valueStyle.Render(location),
		"",
		valueStyle.Render(bio),
	)

	profileContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		avatar,
		"   ",
		info,
	)

	return profile.Render(profileContent)
}

func projects(m Model) string {

	var projects string

	for _, p := range m.Profile.Projects {

		pStatus := "~"

		if p.IsValidated {
			pStatus = "✓"
		}

		line := fmt.Sprintf("• %s (%d) %s", p.ProjectName, p.FinalMark, pStatus)

		projects += lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Render(line) + "\n"
	}

	return projects
}

func ProfileView(m Model) string {

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	root := lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center).
		Padding(2, 0)

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		projects(m),
		profile(m),
		projects(m),
	)

	return root.Render(content)
}
