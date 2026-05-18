package ui

import (
	"io"
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
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
	Location          string    `json:"location"`

	Groups            []any      `json:"groups"`
	Titles            []any      `json:"titles"`

	Projects		  []Project	 `json:"projects"`
	
	Language		  string
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
	Dump				io.Writer			// Tmp debug
	Login  				LoginForm			// Login UI
	Status  			Status				// Server status when starting up

	Profile 			ProfileData			// Hold information to render profile
	ProfileLastUpdate 	int64				// Track last time profile was updated

	Page    			string				// Current page (ui flow)

	Err					error				// Tmp debug

	ProjectsViewport    viewport.Model		// Scrollview for project list section
}

func InitialModel() Model {

	// Projects viewport
	vp := viewport.New(
		viewport.WithWidth(30),
		viewport.WithHeight(16),
	)

	return Model{
		Login:		NewLoginForm(),
		Status:		Status{},
		Profile:	ProfileData{},
		Page:		"startup",
		ProfileLastUpdate: 0,
		ProjectsViewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
