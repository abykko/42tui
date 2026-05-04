package ui

import (
    "strconv"
	"42cli/api"
    "42cli/server-deployment"
    tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit

        case "1", "2", "3":
            m.Selected, _ = strconv.Atoi(msg.String())
            return m, nil

        case "q":
            if m.Selected == 2 {
                deployment.Stop()
                m.Status.Container = false
            }
            return m, nil

        case "enter":
            return handleEnter(m)
        }
    }
    return m, nil
}

func handleEnter(m Model) (Model, tea.Cmd) {
    if m.Selected == 2 {
        if err := deployment.Build(); err != nil { return m, nil }
        if _, err := deployment.Run(); err != nil { return m, nil }
        
        m.Status.Container = true
        m.Status.Api = true
        m.Status.Session = true
    }

    if m.Selected == 3 {
        resp, _, err := api.DoSignedRequest("/status")
        if err == nil {
            m.Profile.Name = resp["status"].(string)
        }
    }
    return m, nil
}
