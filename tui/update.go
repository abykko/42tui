package tui

import (
	// "time"
	"log"
	tea "charm.land/bubbletea/v2"

	"42cli/conf"
	"42cli/tui/service"
	deployment "42cli/server-deployment"
)

// Messages

type AutologinMsg struct {
	User	string
	Pass	string
}

type LoginMsg struct {}

type ProfileMsg struct {
	User	string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	m.ProjectsViewport, cmd = m.ProjectsViewport.Update(msg)

	// Keyboard interruption handler
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		log.Println("KEY:", msg.String())
		if msg.String() == "ctrl+c" {
			log.Println("[System] Saliendo del programa (Ctrl+C)...")
			return m, tea.Quit
		}

		var err error
		if msg.String() == "a" {
			m.Profile, err = service.Profile("msouiyeh")
		}
		if msg.String() == "b" {
			m.Profile, err = service.Profile("iamrani-")
		}

		if err != nil {
			return m, tea.Quit
		}
	
	case AutologinMsg:

		m.Page = "autologin"

		log.Println("[Autologin] Obtaining credentials")
		log.Println("[Autologin] Welcome", msg.User)

		ok, err := service.Login(msg.User, msg.Pass)
		if err != nil {
			log.Println("[Autologin] Error:", err)
			return m, tea.Quit
		}
		if !ok {
			log.Println("[Autologin] Autologin failed")
			log.Println("[Autologin] Autologin failed")
			return m, func() tea.Msg {
				return LoginMsg{}
			}
		}
		if ok {
			log.Println("[Autologin] Autologin success")
			return m, func() tea.Msg {
				return ProfileMsg{
					User: msg.User,
				}
			}
		}

	case LoginMsg:

		m.Page = "login"
		log.Println("Loginmsg")
	
	case ProfileMsg:

		m.Page = "profile"

		var err error
		m.Profile, err = service.Profile(msg.User)
		if err != nil {
			log.Println("[ProfileMsg] Error", err)
			return m, tea.Quit
		}

		return m, func() tea.Msg {
			return 0
		}
	}

	if m.Page == "startup" {
		m.Page = ""
		// Deploy server
		ok, err := service.Startup()
		if err != nil {
			deployment.Stop()
			return m, tea.Quit
		}
		if ok {
			// Check if autologin is enabled
			autoLogin, err := conf.GetBool("autologin")
			if err != nil {
				log.Println("[Startup] error reading autologin setting:", err)
				return m, tea.Quit
			}

			if !autoLogin {
				log.Println("[Startup] GO LOGIN PAGE")
				return m, func() tea.Msg {
					return LoginMsg{}
				}
			}

			// Start autologin
			log.Println("[Startup] autologin enabled")

			user, pass, _ := service.Autologin()

			return m, func() tea.Msg {
				return AutologinMsg{
					User: user,
					Pass: pass,
				}
			}
		}
	}

	return m, cmd
}