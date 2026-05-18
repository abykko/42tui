package ui

import (
	"fmt"
	"os"
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

	name := fmt.Sprintf("%s %s", m.Profile.FirstName, m.Profile.LastName)
	username := "@" + m.Profile.Login

	location := m.Profile.Location
	if location == "" {
		location = "Not logged in campus."
	}

	points := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Ev.P %d", m.Profile.EvaluationPoints),
		fmt.Sprintf("₳ %d", m.Profile.Wallet),
		fmt.Sprintf("Lan %s", languageEmoji(m.Profile.Language)),
	)

	info := lipgloss.JoinVertical(lipgloss.Left,
		nameStyle.Render(name),
		usernameStyle.Render(username),
		labelStyle.Render("Email: ")+valueStyle.Render(m.Profile.Email),
		labelStyle.Render("Location: ")+valueStyle.Render(location),
		valueStyle.Render("Hello World Bio!"),
	)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		points,
		"    ",
		info,
	)

	return profile.Render(content)
}

func projects(m Model) string {

	const (
		totalWidth = 30
		nameWidth  = 22
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

		// truncate name (safe)
		name := p.ProjectName
		r := []rune(name)

		if len(r) > nameWidth {
			if nameWidth > 3 {
				r = append(r[:nameWidth-3], []rune("...")...)
			} else {
				r = r[:nameWidth]
			}
		}

		name = string(r)

		// IMPORTANT: let lipgloss handle width (no manual padding, no len hacks)
		namePart := nameStyle.Width(nameWidth).Render(name)

		scorePart := scoreStyle.Width(scoreWidth).Render(
			fmt.Sprintf("%d", int(p.FinalMark)),
		)

		statusPart := statusStyle.Render(status)

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			namePart,
			scorePart,
			statusPart,
		)

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

	if m.Profile.ID == 0 {

		ascii := `         
		  .-'''-.                                                                       
.---.   '   _    \         _______                                                     
|   | /   /' '.   \        \  ___ ' '.   .--.   _..._                                   
|   |.   |     \  '         ' |--.\  \  |__| .'     '.   .--./)                        
|   ||   '      |  '        | |    \  ' .--..   .-.   . /.''\\                         
|   |\    \     / /  __     | |     |  '|  ||  '   '  || |  | |                        
|   | '.   ' ..' /.:--.'.   | |     |  ||  ||  |   |  | \' -' /                         
|   |    '-...-' '/ |   \ |  | |     ' .'|  ||  |   |  | /("'')                          
|   |            '"' __ | |  | |___.' /' |  ||  |   |  | \ '---.   ,.--.  ,.--.  ,.--.  
|   |             .'.''| | /_______.'/  |__||  |   |  |  /'""'.\ //    \//    \//    \ 
'---'            / /   | |_\_______|/       |  |   |  | ||     ||\\    /\\    /\\    / 
                 \ \._,' '/                 |  |   |  | \'. __//  ''--'  ''--'  ''--'  
                  '---'  '"                  '--'   '--'  '---'                       
`
		return root.Render(ascii)
	}

	vp := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#59ff00")).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		Padding(0, 2).
		Render(m.ProjectsViewport.View())

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		profile(m),
		vp,
	)

	return root.Render(content)
}
