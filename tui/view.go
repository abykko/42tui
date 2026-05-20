package tui

import (
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"

	"42cli/tui/views"
)

func (m Model) View() tea.View {

	m.ProjectsViewport.SetContent(views.Projects(m.Profile.Projects))

	// espacio entre elementos de una row
	itemStyle := lipgloss.NewStyle().MarginRight(2)

	// espacio entre rows
	rowStyle := lipgloss.NewStyle().MarginBottom(1)

	rows := [][]string{
		{
			// Uncomment this line to showup a simple status bar of the
			// 	itemStyle.Render(ServerStatusView(m)),
			itemStyle.Render(views.Hotbar(m.Clock.Day, m.Clock.Time)),
		},
		{},
	}

	if m.Page == "startup" {
		rows[1] = append(rows[1], itemStyle.Render(views.Startup()))
	}
	if m.Page == "autologin" {
		rows[1] = append(rows[1], itemStyle.Render(views.Autologin()))
	}
	if m.Page  == "profile" {
		rows[1] = append(rows[1], itemStyle.Render(views.Profile(m.Profile, m.ProjectsViewport)))
	}

	var renderedRows []string

	for _, row := range rows {
		renderedRows = append(
			renderedRows,
			rowStyle.Render(
				lipgloss.JoinHorizontal(lipgloss.Left, row...),
			),
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		renderedRows...,
	)

	v := tea.NewView(content)
	v.AltScreen = true

	return v
}