package ui

import (
	"fmt"
	"os"
	// "log"
	"strings"
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

	location := m.Profile.Location
	if location == "" {
		location = "Not logged in campus."
	}

	bio := "Hello World Bio!"

	points := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Ev.P %d", evPoint),
		fmt.Sprintf("₳ %d", wallet),
		fmt.Sprintf("Lan %s", languageEmoji(m.Profile.Language)),
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

	const (
		totalWidth = 30
		nameWidth  = 23
		scoreWidth = 3
	)

	var b strings.Builder

	for i, p := range m.Profile.Projects {

		status := "✘"
		statusColor := lipgloss.Color("#ff4848")
		scoreColor := lipgloss.Color("#ff4848")
		nameColor := lipgloss.Color("#ffffff")

		if p.IsValidated {
			status = "✓"
			statusColor = lipgloss.Color("#9cff67")
			scoreColor = lipgloss.Color("#9cff67")
		}

		switch i % 4 {
		case 0:
			nameColor = lipgloss.Color("#EDEEC9")
		case 1:
			nameColor = lipgloss.Color("#EDEEC9")
		case 2:
			nameColor = lipgloss.Color("#BFD8BD")
		case 3:
			nameColor = lipgloss.Color("#98C9A3")
		}

		nameStyle := lipgloss.NewStyle().Foreground(nameColor)
		scoreStyle := lipgloss.NewStyle().Foreground(scoreColor)
		statusStyle := lipgloss.NewStyle().Foreground(statusColor)

		// truncate name
		nameRunes := []rune(p.ProjectName)

		if len(nameRunes) > nameWidth {
			if nameWidth > 3 {
				nameRunes = append(nameRunes[:nameWidth-3], []rune("...")...)
			} else {
				nameRunes = nameRunes[:nameWidth]
			}
		}

		name := string(nameRunes)

		namePart := nameStyle.Render(fmt.Sprintf("%-*s", nameWidth, name))
		scorePart := scoreStyle.Render(
			fmt.Sprintf("%*d", scoreWidth, int(p.FinalMark)),
		)
		statusPart := statusStyle.Render(status)

		line := fmt.Sprintf("%s %s %s", namePart, scorePart, statusPart)

		// opcional: asegurar padding final si algo falla
		runes := []rune(line)
		if len(runes) < totalWidth {
			line += strings.Repeat(" ", totalWidth-len(runes))
		}

		line = line + "\x1b[0K"

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
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

	vp := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#59ff00")).
		BorderTop(false).BorderRight(false).BorderBottom(false).
		Padding(0, 2).
		Render(m.ProjectsViewport.View())

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		profile(m),
		vp,
	)

	return root.Render(content)
}
