package ui

import (
	"fmt"
	"os"
	"golang.org/x/term"
	"charm.land/lipgloss/v2"
)

func languageEmoji(lang string) string {
	flags := map[string]string{
		"es": "🇪🇸",
		"en": "🇬🇧",
		"fr": "🇫🇷",
		"de": "🇩🇪",
		"it": "🇮🇹",
		"pt": "🇵🇹",
		"ru": "🇷🇺",
		"ja": "🇯🇵",
		"ko": "🇰🇷",
		"zh": "🇨🇳",
	}

	if emoji, ok := flags[lang]; ok {
		return emoji
	}
	return lang
}

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

	firstName := m.Profile.FirstName
	lastName := m.Profile.LastName
	name := fmt.Sprintf("%s %s", firstName, lastName)

	username := fmt.Sprintf("@%s", m.Profile.Login)

	email := m.Profile.Email

	evPoint := m.Profile.EvaluationPoints

	wallet := m.Profile.Wallet

	location := "Málaga(THIS LOCATION IS MOCKED)"

	bio := "Vamos a hacer intra pero mejor, será buena idea ;)"

	points := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Ev.P %d", evPoint),
		fmt.Sprintf("₳ %d", wallet),
		languageEmoji(m.Profile.Language),
	)

	info := lipgloss.JoinVertical(lipgloss.Left,
		nameStyle.Render(name),
		usernameStyle.Render(username),
		labelStyle.Render("Email: ")+valueStyle.Render(email),
		labelStyle.Render("Location: ")+valueStyle.Render(location),
		valueStyle.Render(bio),
	)

	profileContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		points,
		"    ",
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
	
	id := m.Profile.ID
	if id == 0 {
		return root.Render("Loading...")
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		profile(m),
		projects(m),
	)

	return root.Render(content)
}
