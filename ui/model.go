package ui

import (
	"io"
	"os"
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
)

type Status struct {
	Container bool
	Api       bool
	Session   bool
}

type Project struct {
	Children       []any   `json:"children"`
	FinalMark      int     `json:"final_mark"`
	IsValidated    bool    `json:"is_validated"`
	LastEventDate  string  `json:"last_event_date"`
	Occurrence     int     `json:"occurrence"`
	ProjectName    string  `json:"project_name"`
	ProjectSlug    string  `json:"project_slug"`
	ProjectsUserID  any    `json:"projects_user_id"`
}

type ProfileData struct {
	ID                int        `json:"id"`
	UserDataID        int        `json:"user_data_id"`
	Login             string     `json:"login"`
	DisplayedLogin    string     `json:"displayed_login"`
	Email             string     `json:"email"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Phone             string     `json:"phone"`
	Image             string     `json:"image"`
	ProfilePicture    string     `json:"profile_picture"`
	IsActive          bool       `json:"is_active"`
	Wallet            int        `json:"wallet"`
	EvaluationPoints  int        `json:"evaluation_points"`

	AlumnizedAt       *string    `json:"alumnized_at"`
	Close             *string    `json:"close"`
	DataErasureDate   string     `json:"data_erasure_date"`
	Location          *string    `json:"location"`

	Groups            []any      `json:"groups"`
	Titles            []any      `json:"titles"`

	Projects		  []Project	 `json:"projects"`
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
	Dump				io.Writer
	Login  				LoginForm
	Status  			Status
	Profile 			ProfileData
	ProfileLastUpdate 	int64
	Page    			string
	Err					error
}

func InitialModel() Model {

	var dump *os.File
    var err error
    
    dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
    if err != nil {
        os.Exit(1)
    }

	return Model{
		Dump:		dump,
		Login:		NewLoginForm(),
		Status:		Status{},
		Profile:	ProfileData{},
		Page:		"startup",
		ProfileLastUpdate: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
