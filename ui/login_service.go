package ui

import (
    "fmt"
    "log"

    tea "charm.land/bubbletea/v2"
    "42cli/api"
    "42cli/conf"
)

func LoginService(m Model) (tea.Model, tea.Cmd) {
    user := m.Login.UsernameInput.Value()
    passwd := m.Login.PasswordInput.Value()

    if user == "" || passwd == "" {
        log.Println("[Login] Error: Campos vacíos")
        return m, nil
    }

    log.Printf("[Login] Intentando login para usuario: %s\n", user)

    resp, _, err := api.DoRequest(fmt.Sprintf("/session/login?username=%s&password=%s", user, passwd))
    if err != nil {
        log.Printf("[Login] Error de red/API: %v\n", err)
        return m, nil
    }

    loginSuccess, ok := resp["login"].(bool)
    if !ok || !loginSuccess {
        log.Println("[Login] Login rechazado: credenciales incorrectas")
        return m, nil
    }

    log.Println("[Login] Login exitoso, cambiando a perfil")
    m.Status.Session = true
    m.Page = "profile"
    conf.Set("logged_with", user)

    return m, nil
}
