package views

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"charm.land/lipgloss/v2"

	"42tui/conf"
)

// Hotbar renders a single-line, 3-column top status bar matching the terminal width.
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

	// Subdivide terminal width dynamically into three chunks.
	// The right column takes the remaining balance to absorb integer division losses.
	colWidth := w / 3
	rightColWidth := w - (colWidth * 2)

	// Combine data into a single-line format for each column
	leftContent := fmt.Sprintf("Intranet (%s)", loggedUser)
	centerContent := fmt.Sprintf("%s %s", day, time)
	rightContent := "42tui v1.0.0"

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.PlaceHorizontal(colWidth, lipgloss.Left, leftContent),
		lipgloss.PlaceHorizontal(colWidth, lipgloss.Center, centerContent),
		lipgloss.PlaceHorizontal(rightColWidth, lipgloss.Right, rightContent),
	)

	return style.Render(content)
}