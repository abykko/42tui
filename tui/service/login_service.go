package service

import (
    "fmt"

    "42cli/api"
)

func Login(user, passwd string) (bool, error) {
    resp, _, err := api.DoRequest(
        fmt.Sprintf("/session/login?username=%s&password=%s", user, passwd),
    )
    if err != nil {
        return false, err
    }

    ok, _ := resp["login"].(bool)
    return ok, nil
}