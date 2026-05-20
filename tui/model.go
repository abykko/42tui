package tui

import (
	"io"
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	"42cli/tui/service"
)

type LoginForm struct {
	FocusIndex    int
	UsernameInput textinput.Model
	PasswordInput textinput.Model
}

type ClockData struct {
	Day		string
	Time	string
}

type Model struct {
	Dump				io.Writer				// Tmp debug
	Login  				LoginForm				// Login UI
	
	Profile 			service.ProfileData	// Hold information to render profile
	ProfileLastUpdate 	int64					// Track last time profile was updated
	
	Page    			string					// Current page (ui flow)
	
	Err					error					// Tmp debug
	
	ProjectsViewport    viewport.Model			// Scrollview for project list section
	
	Clock				ClockData
}

func NewLoginForm() LoginForm {
	u := textinput.New()
	u.Placeholder = "Username"
	u.Focus()
	u.SetWidth(20)

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.SetWidth(20)

	return LoginForm{
		FocusIndex:    0,
		UsernameInput: u,
		PasswordInput: p,
	}
}

func InitialModel() Model {

	// Projects viewport
	vp := viewport.New(
		viewport.WithWidth(30),
		viewport.WithHeight(16),
	)

	return Model{
		Login:		NewLoginForm(),
		Profile:	service.ProfileData{},
		Page:		"startup",
		ProfileLastUpdate: 0,
		ProjectsViewport: vp,
		Clock: ClockData{Day: "", Time: "",},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
