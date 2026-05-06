package ui

import (
    "fmt"
    "log"
    "os"

    tea "charm.land/bubbletea/v2"
    "42cli/api"
    "42cli/conf"
)

func LoginService(m Model) (tea.Model, tea.Cmd) {
    user := m.Login.UsernameInput.Value()
    passwd := m.Login.PasswordInput.Value()

    // Abrir o crear archivo de log
    f, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    logger := log.New(f, "[LOGIN] ", log.LstdFlags)

    if user == "" || passwd == "" {
        logger.Println("Error: Campos vacíos")
        return m, nil
    }

    logger.Printf("Intentando login para usuario: %s\n", user)

    resp, _, err := api.DoRequest(fmt.Sprintf("/session/login?username=%s&password=%s", user, passwd))
    if err != nil {
        logger.Printf("Error de red/API: %v\n", err)
        return m, nil
    }

    loginSuccess, ok := resp["login"].(bool)
    if !ok || !loginSuccess {
        logger.Println("Login rechazado: credenciales incorrectas")
        return m, nil
    }

    logger.Println("Login exitoso, cambiando a perfil")
    m.Status.Session = true
    m.Page = "profile"
    conf.Set("logged_with", user)

    return m, nil
}
