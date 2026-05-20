package views

import (
	"fmt"
	"os"
	"strings"
	"crypto/sha256"
	"golang.org/x/term"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"42cli/tui/service"
)

func profilePicture(userID int) string {
	// hash del usuario
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", userID)))
	hash := h[:]

	palette := []string{
		"#5F5FAF", // azul violeta
		"#FF5FAF", // rosa
		"#7CFF6B", // verde
		"#FFD166", // amarillo
		"#06D6A0", // aqua
		"#EF476F", // rojo suave
		"#A78BFA", // violeta claro
	}

	width := 6
	height := 6

	var b strings.Builder

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			// usa el hash para elegir color determinista
			idx := int(hash[(i*width+j)%len(hash)]) % len(palette)

			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(palette[idx]))

			b.WriteString(style.Render("██"))
		}
		if i != height - 1 {
			b.WriteString("\n")
		}
	}

	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(b.String())
}

func profile(p service.ProfileData, maxSize int) string {

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5FAF")).
		Bold(true)

	usernameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#615c65"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6f6f6f"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D0D0D0"))

	fullName := fmt.Sprintf("%s %s", p.FirstName, p.LastName)
	username := "@" + p.Login

	location := p.Location
	if location == "" {
		location = "Not logged."
	}

	header := lipgloss.NewStyle().
		Width(maxSize).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				nameStyle.Render(fullName),
				lipgloss.NewStyle().
					Width(maxSize - len(fullName) - len(username)).
					Render(""),
				usernameStyle.Render(username),
			),
		)

	info := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		labelStyle.Render("Email: ")+valueStyle.Render(p.Email),
		labelStyle.Render("Location: ")+valueStyle.Render(location),
		labelStyle.Render("Bio {"),
		valueStyle.Render(fmt.Sprintf("    text: Hello, I'm %s! 🐌", p.FirstName)),
		labelStyle.Render("};"),
	)

	// We pass the ID to generate a unique noise profile picture
	avatar := profilePicture(p.ID)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		avatar,
		" ",
		info,
	)

	return content
}

func Projects(projects []service.Project) string {

	const (
		totalWidth = 30
		nameWidth  = 15
		scoreWidth = 5
	)

	var projectsBuilder strings.Builder

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
			nameColor = lipgloss.Color("#F7CAD0")
		case 1:
			nameColor = lipgloss.Color("#DEA6AF")
		case 2:
			nameColor = lipgloss.Color("#B0A8B9")
		case 3:
			nameColor = lipgloss.Color("#BBCCB4")
		}

		nameStyle := lipgloss.NewStyle().Foreground(nameColor)
		scoreStyle := lipgloss.NewStyle().Foreground(scoreColor)
		statusStyle := lipgloss.NewStyle().Foreground(statusColor)

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

		projectsBuilder.WriteString(line)
		projectsBuilder.WriteByte('\n')
	}

	projectsBuilder.WriteString("EOF")

	projectsList := projectsBuilder.String()

	projectsHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e75fff"))

	projectsHeader := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectsHeaderStyle.Render("⛶ Projects ⛶"),
	)

	content := lipgloss.JoinVertical(lipgloss.Center,
		projectsHeader,
		projectsList,
	)

	return content
}

func Profile(p service.ProfileData, projectsVp viewport.Model) string {

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if p.ID == 0 {
		return lipgloss.NewStyle().
			Width(w).
			AlignHorizontal(lipgloss.Center).
			Render("loading")
	}

	// Profile box con altura fija
	profileBox := lipgloss.NewStyle().
		AlignVertical(lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).BorderBottom(false).BorderRight(false).
		BorderForeground(lipgloss.Color("#FF5FAF")).
		Padding(0, 1).
		Render(profile(p, 42))

	// Projects box con altura fija
	projectsBox := lipgloss.NewStyle().
		AlignVertical(lipgloss.Top).
		AlignVertical(lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).BorderBottom(false).BorderLeft(false).
		BorderForeground(lipgloss.Color("#FF5FAF")).
		Padding(0, 1).
		Render(projectsVp.View())

	// Layout horizontal con alineación consistente
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		profileBox,
		" ",
		projectsBox,
	)

	// Root centrado
	root := lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center)

	return root.Render(content)
}