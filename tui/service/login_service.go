package service

import (
    "fmt"

    "42tui/api"
    "42tui/conf"
)

func Login(user, passwd string) (bool, error) {
    resp, _, err := api.DoRequest(
        fmt.Sprintf("/session/login?username=%s&password=%s", user, passwd),
    )
    if err != nil {
        return false, err
    }

    ok, _ := resp["login"].(bool)

    if ok {
        err := conf.Set("logged_with", user)
        if err != nil {
            return ok, err
        }
    }

    return ok, nil
}