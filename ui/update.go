package ui

import tea "charm.land/bubbletea/v2"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// startup routing
	if m.Page == "startup" && m.Status.Container && m.Status.Api {

		if !m.Status.Session {
			m.Page = "login"
			return m, nil
		}

		m.Page = "profile"
		return m, nil
	}

	// startup flow
	if m.Page == "startup" {
		return StartupService(m)
	}

	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		// global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// login page input handling
		if m.Page == "login" {

			switch msg.String() {

            case "ctrl+c":
                return m, tea.Quit

			case "tab":
				m.Login.FocusIndex = (m.Login.FocusIndex + 1) % 2
				m = updateLoginFocus(m)

			case "shift+tab":
				m.Login.FocusIndex--
				if m.Login.FocusIndex < 0 {
					m.Login.FocusIndex = 1
				}
				m = updateLoginFocus(m)
			}

			var cmd tea.Cmd

			m.Login.UsernameInput, cmd = m.Login.UsernameInput.Update(msg)
			m.Login.PasswordInput, _ = m.Login.PasswordInput.Update(msg)

			return m, cmd
		}
	}

	return m, nil
}

func updateLoginFocus(m Model) Model {

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