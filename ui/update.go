package ui

import (
	"log"
	"github.com/davecgh/go-spew/spew"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	log.Println(m.Page)
	if m.Dump != nil {
		spew.Fdump(m.Dump, msg)
	}

	// Log de cada mensaje recibido para debug profundo
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			log.Println("[System] Saliendo del programa (Ctrl+C)...")
			return m, tea.Quit
		}
	}
	
	// startup flow
	if m.Page == "startup" {
		log.Println("[Flow] Ejecutando StartupService...")
		return StartupService(m)
	}

	if m.Page == "profile" {
		return ProfileService(m)
	}

	// Manejo de lógica por teclado y página
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		// login page input handling
		if m.Page == "login" {
			log.Printf("[Login] Manejando input en login. Foco actual: %d", m.Login.FocusIndex)

			switch msg.String() {
			case "tab", "shift+tab":
				if msg.String() == "tab" {
					m.Login.FocusIndex = (m.Login.FocusIndex + 1) % 2
				} else {
					m.Login.FocusIndex--
					if m.Login.FocusIndex < 0 {
						m.Login.FocusIndex = 1
					}
				}
				log.Printf("[Login] Cambio de foco a: %d", m.Login.FocusIndex)
				m = updateLoginFocus(m)

			case "enter":
				log.Println("[Login] Enter presionado. Intentando LoginService...")
				return LoginService(m)
			}

			var cmd tea.Cmd
			m.Login.UsernameInput, cmd = m.Login.UsernameInput.Update(msg)
			m.Login.PasswordInput, _ = m.Login.PasswordInput.Update(msg)
			return m, cmd
		}

		if m.Page == "profile" {
			if msg.String() == "enter" {
				log.Println("[Profile] Enter presionado en página de perfil")
				return ProfileService(m)
			}
		}
	}

	return m, nil
}

func updateLoginFocus(m Model) Model {
	log.Printf("[UI] Actualizando foco visual: %d", m.Login.FocusIndex)
	switch m.Login.FocusIndex {
	case 0:
		m.Login.UsernameInput.Focus()
		m.Login.PasswordInput.Blur()
	case 1:
		m.Login.UsernameInput.Blur()
		m.Login.PasswordInput.Focus()
	}
	return m
}