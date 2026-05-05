package ui

import (
	"fmt"
	"time"
	"42cli/api"
    "42cli/server-deployment"
    tea "charm.land/bubbletea/v2"
)

func StartupService(m Model) (Model, tea.Cmd) {
	// Stop button
	if m.Status.Container == true {
		deployment.Stop()
		m.Status.Container = false
		return m, nil
	}

	defer func() {
		if m.Status.Container == false {
			deployment.Stop()
		}
	}()

	// In case this guard conditions throw an error we stop the server with defer function
	if err := deployment.Build(); err != nil {return m, nil}
	if _, err := deployment.Run(); err != nil {return m, nil}
	
	m.Status.Container = true // Container is running
	
	// We wait for the api
	err := api.WaitForRequestTo(
		"/status",
		func(resp map[string]interface{}) bool {
			return resp["status"] == "ok"
		},
		10*time.Second,         // timeout
		500*time.Millisecond,    // interval
	)

	if err != nil {
		fmt.Println(err)
		return m, nil
	}

	m.Status.Api = true

	resp, _, err := api.DoRequest("/session/expired")
	if err == nil {
		expired := resp["expired"].(bool)
		if expired == true {
			m.Status.Session = false
		} else {
			m.Status.Session = true // Session is valid
		}
	} else {
		return m, nil
	}
    return m, nil
}
