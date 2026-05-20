package views

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"charm.land/lipgloss/v2"

	"42cli/conf"
)

func Hotbar(day string, time string) string {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	loggedUser, err := conf.GetString("logged_with")
	if err != nil || loggedUser == "" {
		return ""
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		Width(w)

	// Línea 1
	h1Left := "Intranet"
	h1Center := day
	h1Right := "42cli"

	h1 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.PlaceHorizontal(w/3, lipgloss.Left, h1Left),
		lipgloss.PlaceHorizontal(w/3, lipgloss.Center, h1Center),
		lipgloss.PlaceHorizontal(w/3, lipgloss.Right, h1Right),
	)

	// Línea 2
	h2Left := fmt.Sprintf("Welcome %s", loggedUser)
	h2Center := time
	h2Right := "v1.0.0"

	h2 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.PlaceHorizontal(w/3, lipgloss.Left, h2Left),
		lipgloss.PlaceHorizontal(w/3, lipgloss.Center, h2Center),
		lipgloss.PlaceHorizontal(w/3, lipgloss.Right, h2Right),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		h1,
		h2,
	)

	return style.Render(content)
}