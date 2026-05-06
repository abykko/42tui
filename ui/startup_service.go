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

	log.Println("[Startup] iniciado")

	defer func() {
		if m.Status.Container == false {
			log.Println("[Startup] container no activo, haciendo stop")
			deployment.Stop()
		}
	}()

	// Build
	log.Println("[Startup] building container...")
	if err := deployment.Build(); err != nil {
		log.Println("[Startup] build error:", err)
		return m, nil
	}

	// Run
	log.Println("[Startup] running container...")
	if _, err := deployment.Run(); err != nil {
		log.Println("[Startup] run error:", err)
		return m, nil
	}

	m.Status.Container = true
	log.Println("[Startup] container running")

	// Wait API
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

	// Autologin
	autoLogin, err := conf.GetBool("autologin")
	if err != nil {
		log.Println("[Startup] error leyendo autologin:", err)
		return m, nil
	}

	if !autoLogin {
		log.Println("[Startup] autologin desactivado")
		m.Page = "login"
		return m, func() tea.Msg { return 1 }
	}

	log.Println("[Startup] autologin activado")

	loginUser, err := conf.GetString("user_login")
	if err != nil {
		log.Println("[Startup] error leyendo user_login:", err)
		return m, nil
	}

	loginPasswd, err := conf.GetString("passwd_login")
	if err != nil {
		log.Println("[Startup] error leyendo passwd_login:", err)
		return m, nil
	}

	if loginUser == "" || loginPasswd == "" {
		log.Println("[Startup] error no hay credenciales guardadas:", err)
		m.Page = "login"
		return m, nil
	}

	log.Println("[Startup] credenciales cargadas para:", loginUser)

	m.Login.UsernameInput.SetValue(loginUser)
	m.Login.PasswordInput.SetValue(loginPasswd)

	log.Println("[Startup] ejecutando LoginService...")

	newModel, _ := LoginService(m)
	m = newModel.(Model)

	log.Println("[Startup] login completado")

	return m, func() tea.Msg { return 1 }
}