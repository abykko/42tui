package views

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"42cli/tui/service"
)

/* ----------------------------- PROFILE PICTURE ----------------------------- */

func profilePicture(userID int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", userID)))
	hash := h[:]

	palette := []string{
		"#5F5FAF",
		"#FF5FAF",
		"#7CFF6B",
		"#FFD166",
		"#06D6A0",
		"#EF476F",
		"#A78BFA",
	}

	const width = 6
	const height = 6

	var b strings.Builder

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			idx := int(hash[(i*width+j)%len(hash)]) % len(palette)

			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(palette[idx]))

			b.WriteString(style.Render("██"))
		}
		if i != height-1 {
			b.WriteString("\n")
		}
	}

	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(b.String())
}

/* -------------------------------- PROFILE -------------------------------- */

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
					Width(maxSize-len(fullName)-len(username)).
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

	avatar := profilePicture(p.ID)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		avatar,
		" ",
		info,
	)
}

/* -------------------------------- PROJECTS -------------------------------- */

func Projects(projects []service.Project) string {
	const (
		nameWidth  = 15
		scoreWidth = 5
	)

	var sb strings.Builder

	for i, p := range projects {
		status := "✘"
		statusColor := lipgloss.Color("#ff4848")
		scoreColor := lipgloss.Color("#ff4848")
		nameColor := lipgloss.Color("#ffffff")

		if p.IsValidated {
			status = "✓"
			statusColor = lipgloss.Color("#2ecc71")
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

		name := []rune(p.ProjectName)
		if len(name) > nameWidth {
			if nameWidth > 3 {
				name = append(name[:nameWidth-3], []rune("...")...)
			} else {
				name = name[:nameWidth]
			}
		}

		namePart := nameStyle.Width(nameWidth).Render(string(name))
		scorePart := scoreStyle.Width(scoreWidth).Render(fmt.Sprintf("%d", int(p.FinalMark)))
		statusPart := statusStyle.Render(status)

		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, namePart, scorePart, statusPart))
		sb.WriteByte('\n')
	}

	sb.WriteString("EOF")

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5FAF")).
		Bold(true).
		Render("Submitted Projects ⛶")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		sb.String(),
	)
}

/* ---------------------------------- PACE ---------------------------------- */

func pace(p service.ProfileData) string {
	currentMilestone := p.Pace.Milestone
	currentPace := p.Pace.Pace
	currentDeadline := p.Pace.Deadline

	if len(p.Pace.Milestones) > 6 {
		if p.Pace.Milestones[6].ValidatedAt != "" {
			currentMilestone = 7
		}
	}

	daysUntil := func(deadlineStr string) int {
		now := time.Now().UTC()
		currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		deadlineDate, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			return 0
		}

		return int(deadlineDate.Sub(currentDate).Hours() / 24)
	}

	paces := []int{8, 12, 15, 18, 22, 24}

	nextPace := func(current int) (int, bool) {
		for i, v := range paces {
			if v == current {
				if i+1 < len(paces) {
					return paces[i+1], true
				}
				return 0, false
			}
		}
		return 0, false
	}

	buildHeader := func() string {
		next, ok := nextPace(currentPace)
		days := daysUntil(currentDeadline)

		if !ok {
			return fmt.Sprintf("Pace %d | blackhole", currentPace)
		}

		return fmt.Sprintf(
			"Pace %d | %d days left to Pace %d",
			currentPace,
			days,
			next,
		)
	}

	header := buildHeader()

	passedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ecc71")).
		Padding(0, 1)

	currentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3e7eff")).
		Bold(true).
		Padding(0, 1)

	pendingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7f8c8d")).
		Padding(0, 1)

	var sb strings.Builder

	for i := 0; i <= 6; i++ {
		label := fmt.Sprintf("Milestone %d", i)

		var line string
		switch {
		case i < currentMilestone:
			line = passedStyle.Render("✔ " + label)
		case i == currentMilestone:
			line = currentStyle.Render(label + " ◀")
		default:
			line = pendingStyle.Render("○ " + label)
		}

		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(line)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		sb.String(),
	)
}

/* --------------------------------- PROFILE -------------------------------- */

func Profile(p service.ProfileData, projectsVp viewport.Model) string {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if p.ID == 0 {
		return lipgloss.NewStyle().
			Width(w).
			AlignHorizontal(lipgloss.Center).
			Render("loading")
	}

	profileInformation := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderRight(false).
		BorderForeground(lipgloss.Color("#FF5FAF")).
		Padding(0, 1).
		Render(profile(p, 40))

	projectsList := lipgloss.NewStyle().
		Padding(0, 1).
		Render(projectsVp.View())

	cursusStatus := lipgloss.NewStyle().
		Padding(2, 1).
		Render(pace(p))

	profileSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			profileInformation,
			cursusStatus,
		),
		" ",
		projectsList,
	)

	return lipgloss.NewStyle().
		Width(w).
		AlignHorizontal(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				profileSection,
				"New sections coming soon...",
			),
		)
}