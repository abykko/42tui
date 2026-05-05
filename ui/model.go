package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
)

type Status struct {
	Container bool
	Api       bool
	Session   bool
}

type ProfileData struct {
	Name string
	User string
}

type LoginForm struct {
	FocusIndex    int
	UsernameInput textinput.Model
	PasswordInput textinput.Model
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

type Model struct {
	Login   LoginForm
	Status  Status
	Profile ProfileData
	Page    string
}

func InitialModel() Model {
	return Model{
		Login:   NewLoginForm(),
		Status:  Status{},
		Profile: ProfileData{},
		Page:    "startup",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
