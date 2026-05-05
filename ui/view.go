package ui

import (
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {

	// espacio entre elementos de una row
	itemStyle := lipgloss.NewStyle().
		MarginRight(2)

	// espacio entre rows
	rowStyle := lipgloss.NewStyle().
		MarginBottom(1)

	rows := [][]string{
		{
			itemStyle.Render(ServerStatusView(m)),
		},
        {},
	}

    if m.Status.Session == false && m.Status.Api == true {
        rows[1] = append(rows[1], itemStyle.Render(LoginView(m)))
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