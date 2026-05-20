package tui

import (
	"log"
	"time"

	tea "charm.land/bubbletea/v2"

	"42cli/conf"
	"42cli/tui/service"
	"42cli/tui/views"
	deployment "42cli/server-deployment"
)

type AutologinMsg struct {
	User string
	Pass string
}

type LoginMsg struct{}

type ProfileMsg struct {
	User string
}

type ClockUpdate struct{}

func getDateTime() (string, string) {
	now := time.Now()
	day := now.Format("02 Jan")
	timeStr := now.Format("15:04:05")
	return day, timeStr
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Initialize the clock 
	if m.Clock.Time == "" {
		m.Clock.Day, m.Clock.Time = getDateTime()
		return m, func() tea.Msg { return ClockUpdate{} } 
	}

	var cmds []tea.Cmd

	var vpCmd tea.Cmd
	m.ProjectsViewport, vpCmd = m.ProjectsViewport.Update(msg)
	cmds = append(cmds, vpCmd)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		log.Println("KEY:", msg.String())

		if msg.String() == "ctrl+c" {
			log.Println("[System] exiting...")
			return m, tea.Quit
		}

		// TMP DEBUG
		if msg.String() == "a" {
			return m, func() tea.Msg {
				return ProfileMsg{User: "msouiyeh"}
			}
		}
		if msg.String() == "a" {
			return m, func() tea.Msg {
				return ProfileMsg{User: "iamrani-"}
			}
		}

	case ClockUpdate:
		m.Clock.Day, m.Clock.Time = getDateTime()

		// keep ticking
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(1 * time.Second)
			return ClockUpdate{}
		})

	case AutologinMsg:

		m.Page = "autologin"

		ok, err := service.Login(msg.User, msg.Pass)
		if err != nil {
			log.Println("[Autologin] error:", err)
			return m, tea.Quit
		}

		if !ok {
			return m, func() tea.Msg { return LoginMsg{} }
		}

		return m, func() tea.Msg {
			return ProfileMsg{User: msg.User}
		}

	case LoginMsg:
		m.Page = "login"

	case ProfileMsg:
		m.Page = "profile"

		var err error
		m.Profile, err = service.Profile(msg.User)
		if err != nil {
			log.Println("[Profile] error:", err)
			return m, tea.Quit
		}

		m.ProjectsViewport.SetContent(views.Projects(m.Profile.Projects))

	case tea.WindowSizeMsg:
		// m.ProjectsViewport.Width = msg.Width
		// m.ProjectsViewport.Height = msg.Height
	}

	if m.Page == "startup" {
		m.Page = ""

		ok, err := service.Startup()
		if err != nil {
			deployment.Stop()
			return m, tea.Quit
		}

		if ok {
			autoLogin, err := conf.GetBool("autologin")
			if err != nil {
				log.Println("[Startup] autologin config error:", err)
				return m, tea.Quit
			}

			if !autoLogin {
				return m, func() tea.Msg { return LoginMsg{} }
			}

			user, pass, _ := service.Autologin()

			return m, func() tea.Msg {
				return AutologinMsg{
					User: user,
					Pass: pass,
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}