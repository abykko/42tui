package ui

import (
    "charm.land/lipgloss/v2"
    tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {
    content := lipgloss.JoinVertical(lipgloss.Left,
        ServerStatusView(m, 1),
        TempView(m, 2),
        ProfileView(m, 3),
    )

    v := tea.NewView(content)
    v.AltScreen = true
    return v
}
