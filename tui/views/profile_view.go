package views

import (
	"fmt"
	"os"
	"strings"
	"golang.org/x/term"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"42cli/tui/service"
)

func profile(p service.ProfileData) string {

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

	name := fmt.Sprintf("%s %s", p.FirstName, p.LastName)
	username := "@" + p.Login

	location := p.Location
	if location == "" {
		location = "Not logged in campus."
	}

	points := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Ev.P %d", p.EvaluationPoints),
		fmt.Sprintf("₳ %d", p.Wallet),
		fmt.Sprintf("Lan %s", p.Language),
	)

	info := lipgloss.JoinVertical(lipgloss.Left,
		nameStyle.Render(name),
		usernameStyle.Render(username),
		labelStyle.Render("Email: ")+valueStyle.Render(p.Email),
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

func Projects(projects []service.Project) string {

	const (
		totalWidth = 30
		nameWidth  = 15
		scoreWidth = 5
	)

	var b strings.Builder

	for i, p := range projects {

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
			// Rosa pastel suave (Base/Destacado sutil)
			nameColor = lipgloss.Color("#F7CAD0")
		case 1:
			// Malva / Rosa viejo (Transición elegante)
			nameColor = lipgloss.Color("#DEA6AF")
		case 2:
			// Lavanda grisáceo (Contraste suave)
			nameColor = lipgloss.Color("#B0A8B9")
		case 3:
			// Salvia / Verde oliva muy suave (Punto de quiebre armónico)
			nameColor = lipgloss.Color("#BBCCB4")
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

func Profile(p service.ProfileData, projectsVp viewport.Model) string {

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	root := lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center)

	if p.ID == 0 {
		ascii := `loading`
		return root.Render(ascii)
	}

	vp := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#59ff00")).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		Padding(0, 2).
		Render(projectsVp.View())

	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		profile(p),
		vp,
	)

	return root.Render(content)
}
