package ui

import (
    tea "charm.land/bubbletea/v2"
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

type Model struct {
    Selected int
    Status   Status
    Profile  ProfileData
}

func InitialModel() Model {
    return Model{
        Selected: 1,
        Status: Status{
            Container: false,
            Api:       false,
            Session:   false,
        },
        Profile: ProfileData{
            Name: "",
            User: "",
        },
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}