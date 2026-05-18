package ui

import (
	"charm.land/lipgloss/v2"
)

func ServerStatusView(m Model) string {

	// estilo base de cada card (CLAVE: Width fijo)
	card := lipgloss.NewStyle().
		Background(lipgloss.Color("#4b6a79")).
		Padding(0, 1).
		MarginRight(1).
		Width(18)

	// estilos de texto
	okStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33ff00"))

	offStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	status := func(name string, active bool) string {
		if active {
			return name + " +"
		}
		return name + " -"
	}

	renderStatus := func(name string, active bool) string {
		if active {
			return okStyle.Render(status(name, true))
		}
		return offStyle.Render(status(name, false))
	}

	container := card.Render(renderStatus("Container", m.Status.Container))
	api := card.Render(renderStatus("Api", m.Status.Api))
	session := card.Render(renderStatus("Session", m.Status.Session))

	line := lipgloss.JoinHorizontal(
		lipgloss.Left,
		container,
		api,
		session,
	)

	return line
}