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
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(60)

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
	location := "Málaga(THIS LOCATION IS MOCKED)"
	bio := "Building cool CLI apps with Go 🚀"

	link := func(url, text string) string {
		return "\033]8;;" + url + "\033\\" + text + "\033]8;;\033\\"
	}

	avatar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("69")).
		Render(
			link(m.Profile.ProfilePicture,
				" ◉ ◉ \n"+
					"  ◉  \n"+
					" ◉ ◉ ",
			),
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

		pStatus := "✘"
		color := "#d32020"
		if p.IsValidated {
			color = "#59ff00"
			pStatus = "✓"
		}

		lineStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))
		
		line := fmt.Sprintf("%s (%d) %s", p.ProjectName, p.FinalMark, pStatus)

		projects += lineStyle.Render(line) + "\n"
	}

	parent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#59ff00")).
		Padding(1, 1).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false)
	
	return parent.Render(projects)
}

func ProfileView(m Model) string {

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	root := lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center).
		Padding(2, 0)

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		profile(m),
		projects(m),
	)

	return root.Render(content)
}
