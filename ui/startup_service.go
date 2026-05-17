package ui

import (
	"log"
	"time"

	"42cli/conf"
	"42cli/api"
	deployment "42cli/server-deployment"
	tea "charm.land/bubbletea/v2"
)

func StartupService(m Model) (Model, tea.Cmd) {

	defer func() {
		if m.Status.Container == false {
			log.Println("[Startup] container no activo, haciendo stop")
			deployment.Stop()
		}
	}()

	// Build container with the server
	log.Println("[Startup] building container...")
	if err := deployment.Build(); err != nil {
		log.Println("[Startup] build error:", err)
		return m, nil
	}

	// Run the server container once is ready
	log.Println("[Startup] running container...")
	if _, err := deployment.Run(); err != nil {
		log.Println("[Startup] run error:", err)
		return m, nil
	}

	m.Status.Container = true
	log.Println("[Startup] container is running")

	// Wait until the API is ready (has a timeout)
	log.Println("[Startup] esperando API /status...")
	err := api.WaitForRequestTo(
		"/status",
		func(resp map[string]interface{}) bool {
			log.Println("[Startup] /status response:", resp)
			return resp["status"] == "ok"
		},
		10*time.Second,
		30*time.Millisecond,
	)
	if err != nil {
		log.Println("[Startup] API no respondió:", err)
		return m, nil
	}

	m.Status.Api = true
	log.Println("[Startup] API OK")

	// Check if autologin is enabled
	autoLogin, err := conf.GetBool("autologin")
	if err != nil {
		log.Println("[Startup] error leyendo autologin:", err)
		return m, nil
	}

	if !autoLogin {
		log.Println("[Autologin] autologin desactivado")

		// Open login page
		m.Page = "login"

		return m, func() tea.Msg { return 1 }
	}

	// Trigger autologin
	m.Page = "autologin"
	return m, func() tea.Msg { return "autologin" }
}