package ui

import (
	"time"
	"log"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	log.Println(m.Page) // Tmp debug
	
	// Keyboard interruption handler
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// m.ProjectsViewport.Width = msg.Width - 4
		// m.ProjectsViewport.Height = msg.Height - 6
	case tea.KeyPressMsg:
		log.Println("KEY:", msg.String())
		if msg.String() == "ctrl+c" {
			log.Println("[System] Saliendo del programa (Ctrl+C)...")
			return m, tea.Quit
		}
	}
	
	// startup setup
	if m.Page == "startup" { return StartupService(m) }
	
	// If we recieve "autologin" msg we must call the autologin service
	// This only happends at the startup if is enabled in settings.
	if s, ok := msg.(string); ok {
		switch s {
		case "autologin":
			return AutologinService(m)
		}
	}
	
	// In case of being in the login page form we handle key inputs
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.Page != "login" { break }

		key := msg.String()
		switch key {
		case "tab":
			m.Login.FocusIndex = (m.Login.FocusIndex + 1) % 2
			m = updateLoginFocus(m)
		case "shift+tab":
			m.Login.FocusIndex--
			if m.Login.FocusIndex < 0 {
				m.Login.FocusIndex = 1
			}
			m = updateLoginFocus(m)
		case "enter":
			log.Println("[Login] Enter presionado. LoginService...")
			return LoginService(m)
		}

		var cmd tea.Cmd
		m.Login.UsernameInput, cmd = m.Login.UsernameInput.Update(msg)
		m.Login.PasswordInput, _ = m.Login.PasswordInput.Update(msg)

		return m, cmd
	}

	// Load profile
	lastUpdate := m.ProfileLastUpdate
	elapsedTime := time.Now().Unix() - lastUpdate
	
	if elapsedTime >= 30 && m.Page == "profile" {
		return ProfileService(m)
	}
	
	// Update project list viewport
	var cmd tea.Cmd
	cmd = nil

	// ONLY update viewport when visible
	if m.Page == "profile" {
		m.ProjectsViewport, cmd = m.ProjectsViewport.Update(msg)
	}


	return m, cmd
}

// Helper function in login form called from key handling
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