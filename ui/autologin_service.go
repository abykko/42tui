package ui

import (
	"log"

	"42cli/conf"
	tea "charm.land/bubbletea/v2"
)

func AutologinService(m Model) (Model, tea.Cmd) {

	log.Println("[Autologin] autologin activado")

	loginUser, err := conf.GetString("user_login")
	if err != nil {
		log.Println("[Autologin] error leyendo user_login:", err)
		return m, nil
	}

	loginPasswd, err := conf.GetString("passwd_login")
	if err != nil {
		log.Println("[Autologin] error leyendo passwd_login:", err)
		return m, nil
	}

	if loginUser == "" || loginPasswd == "" {
		log.Println("[Autologin] error no hay credenciales guardadas:", err)
		m.Page = "login"
		return m, nil
	}

	log.Println("[Autologin] credenciales cargadas para:", loginUser)

	m.Login.UsernameInput.SetValue(loginUser)
	m.Login.PasswordInput.SetValue(loginPasswd)

	log.Println("[Autologin] ejecutando LoginService...")

	newModel, _ := LoginService(m)
	m = newModel.(Model)

	log.Println("[Autologin] login completado")

	return m, func() tea.Msg { return 1 }
}